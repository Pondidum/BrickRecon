using System;
using System.Collections.Generic;
using System.IO;
using System.Linq;
using System.Threading.Tasks;
using System.Xml.Linq;
using BsxProcessor.Domain;
using BsxProcessor.Infrastructure;

namespace BsxProcessor
{
	public class BsxProcessor : IBsxProcessor
	{
		private readonly IFileSystem _fileSystem;
		private readonly Config _config;
		private readonly IImageCacheDispatcher _imageCacheDispatch;
		private readonly IBsxModelBuilder _modelBuilder;

		public BsxProcessor(IFileSystem fileSystem, Config config, IImageCacheDispatcher imageCacheDispatch, IBsxModelBuilder modelBuilder)
		{
			_fileSystem = fileSystem;
			_config = config;
			_imageCacheDispatch = imageCacheDispatch;
			_modelBuilder = modelBuilder;
		}

		public async Task Execute(IEnumerable<BsxRequest> records)
		{
			var tasks = records
				.Select(record => record
					.Start(ConvertToModel)
					.Then(QueueParts)
					.Then(WriteJsonFile));

			await Task.WhenAll(tasks);

			await _imageCacheDispatch.Dispatch();
		}

		private Task<BsxModel> ConvertToModel(BsxRequest request)
		{
			return Task.FromResult(_modelBuilder.Build(request.ModelName, request.Content));
		}

		private Task<BsxModel> QueueParts(BsxModel model)
		{
			_imageCacheDispatch.Add(model.Parts);
			return Task.FromResult(model);
		}

		private async Task WriteJsonFile(BsxModel model)
		{
			var path = new Uri( _config.OutputBucketPath, model.Name + ".json");

			await _fileSystem.WriteJson(new FileData<BsxModel>
			{
				Content = model,
				Drive = path.Host,
				FullPath = path.LocalPath.TrimStart('/')
			});
		}
	}
}
