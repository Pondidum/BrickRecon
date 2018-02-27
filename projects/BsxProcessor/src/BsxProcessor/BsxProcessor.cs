using System.Collections.Generic;
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
		private readonly IImageCacheDispatcher _imageCacheDispatch;
		private readonly BsxModelBuilder _modelBuilder;

		public BsxProcessor(IFileSystem fileSystem, IImageCacheDispatcher imageCacheDispatch, BsxModelBuilder modelBuilder)
		{
			_fileSystem = fileSystem;
			_imageCacheDispatch = imageCacheDispatch;
			_modelBuilder = modelBuilder;
		}

		public async Task Execute(IEnumerable<FileData<XDocument>> records)
		{
			var tasks = records.Select(record => record
				.Start(ConvertToModel)
				.Then(QueueParts)
				.Then(WriteJsonFile));

			await Task.WhenAll(tasks);

			await _imageCacheDispatch.Dispatch();
		}

		private Task<FileData<BsxModel>> ConvertToModel(FileData<XDocument> document)
		{
			var model = _modelBuilder.Build(document);

			return Task.FromResult(new FileData<BsxModel>
			{
				Drive = document.Drive,
				FullPath = $"models/{model.Name}.json",
				Content = model
			});
		}

		private Task<FileData<BsxModel>> QueueParts(FileData<BsxModel> file)
		{
			_imageCacheDispatch.Add(file.Content.Parts);
			return Task.FromResult(file);
		}

		private async Task WriteJsonFile(FileData<BsxModel> file)
		{
			await _fileSystem.WriteJson(file);
		}
	}
}
