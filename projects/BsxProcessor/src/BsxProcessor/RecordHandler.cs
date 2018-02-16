using System.Collections.Generic;
using System.Linq;
using System.Threading.Tasks;
using System.Xml.Linq;
using Amazon.S3.Util;
using BsxProcessor.Domain;
using BsxProcessor.Infrastructure;

namespace BsxProcessor
{
	public class RecordHandler
	{
		private readonly IFileSystem _fileSystem;
		private readonly ImageCacheDispatcher _imageCacheDispatch;
		private readonly BsxModelBuilder _modelBuilder;

		public RecordHandler(IFileSystem fileSystem, ImageCacheDispatcher imageCacheDispatch, BsxModelBuilder modelBuilder)
		{
			_fileSystem = fileSystem;
			_imageCacheDispatch = imageCacheDispatch;
			_modelBuilder = modelBuilder;
		}

		public async Task Execute(IEnumerable<S3EventNotification.S3EventNotificationRecord> records)
		{
			var tasks = records
				.Select(ReadBsxFile)
				.Select(ConvertToModel)
				.Select(QueueParts)
				.Select(WriteJsonFile)
				.ToArray();

			await Task.WhenAll(tasks);

			await _imageCacheDispatch.Dispatch();
		}

		private async Task WriteJsonFile(FileData<BsxModel> file)
		{
			await _fileSystem.WriteJson(file);
		}

		private FileData<BsxModel> QueueParts(FileData<BsxModel> file)
		{
			_imageCacheDispatch.Add(file.Content.Parts);

			return file;
		}

		private FileData<BsxModel> ConvertToModel(Task<FileData<XDocument>> documentTask)
		{
			var document = documentTask.Result;
			var model = _modelBuilder.Build(document);

			return new FileData<BsxModel>
			{
				Drive = document.Drive,
				FullPath = $"models/{model.Name}.json",
				Content = model
			};
		}

		private async Task<FileData<XDocument>> ReadBsxFile(S3EventNotification.S3EventNotificationRecord record)
		{
			return await _fileSystem.ReadXml(record.S3.Bucket.Name, record.S3.Object.Key);
		}
	}
}
