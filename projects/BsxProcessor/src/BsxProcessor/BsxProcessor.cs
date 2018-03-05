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

		public async Task Execute(IEnumerable<FileData<XDocument>> records)
		{
			var tasks = records
				.Where(record => record.Exists)
				.Select(record => record
					.Start(ConvertToModel)
					.Then(QueueParts)
					.Then(WriteJsonFile));

			await Task.WhenAll(tasks);

			await _imageCacheDispatch.Dispatch();
		}

		private Task<BsxModel> ConvertToModel(FileData<XDocument> document)
		{
			return Task.FromResult(_modelBuilder.Build(document));
		}

		private Task<BsxModel> QueueParts(BsxModel model)
		{
			_imageCacheDispatch.Add(model.Parts);
			return Task.FromResult(model);
		}

		private async Task WriteJsonFile(BsxModel model)
		{
			var root = new Uri(_config.OutputBucketPath).Scheme;

			await _fileSystem.WriteJson(new FileData<BsxModel>
			{
				Content = model,
				Drive = root,
				FullPath = Path.Combine(_config.OutputBucketPath.Substring(root.Length + 3), model.Name + ".json")
			});
		}
	}
}
