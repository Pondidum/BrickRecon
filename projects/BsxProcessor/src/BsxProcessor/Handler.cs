using System.Collections.Generic;
using System.Threading.Tasks;
using Amazon.Lambda;
using Amazon.Lambda.Core;
using Amazon.Lambda.S3Events;
using Amazon.Lambda.Serialization.Json;
using Amazon.S3;
using Amazon.S3.Util;
using BsxProcessor.Domain;
using BsxProcessor.Infrastructure;
using Environment = System.Environment;

namespace BsxProcessor
{
	public class Handler
	{
		private readonly IFileSystem _fileSystem;
		private readonly ImageCacheDispatcher _imageCacheDispatch;

		public Handler()
		{
			var config = Config.FromEnvironment();
			var lambdaClient = new AmazonLambdaClient();

			_fileSystem = new S3FileSystem(new AmazonS3Client());
			_imageCacheDispatch = new ImageCacheDispatcher(config, req => lambdaClient.InvokeAsync(req));
		}

		[LambdaSerializer(typeof(JsonSerializer))]
		public void Handle(S3Event s3Event)
		{
			HandleRecords(s3Event.Records).Wait();
		}

		private async Task HandleRecords(IEnumerable<S3EventNotification.S3EventNotificationRecord> records)
		{
			var modelBuilder = new BsxModelBuilder();

			foreach (var record in records)
			{
				var document = await _fileSystem.ReadXml(record.S3.Bucket.Name, record.S3.Object.Key);
				var model = modelBuilder.Build(document);

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
