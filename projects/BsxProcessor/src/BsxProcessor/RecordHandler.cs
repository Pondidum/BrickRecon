using System.Collections.Generic;
using System.Threading.Tasks;
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
			foreach (var record in records)
			{
				var document = await _fileSystem.ReadXml(record.S3.Bucket.Name, record.S3.Object.Key);
				var model = _modelBuilder.Build(document);

				_imageCacheDispatch.Add(model.Parts);

				await _fileSystem.WriteJson(new FileData<BsxModel>
				{
					Drive = document.Drive,
					FullPath = $"models/{model.Name}.json",
					Content = model
				});
			}

			await _imageCacheDispatch.Dispatch();
		}
	}
}
