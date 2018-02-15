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

		public Handler()
		{
			_fileSystem = new S3FileSystem(new AmazonS3Client());
		}

		[LambdaSerializer(typeof(JsonSerializer))]
		public void Handle(S3Event s3Event)
		{
			var config = new Config
			{
				ImageCacheLambda = Environment.GetEnvironmentVariable("IMAGECACHE_LAMBDA")
			};

			HandleRecords(config, s3Event.Records).Wait();
		}

		private async Task HandleRecords(Config config, IEnumerable<S3EventNotification.S3EventNotificationRecord> records)
		{
			var lambdaClient = new AmazonLambdaClient();
			var modelBuilder = new BsxModelBuilder();
			var imageCacheDispatch = new ImageCacheDispatcher(config, req => lambdaClient.InvokeAsync(req));

			foreach (var record in records)
			{
				var document = await _fileSystem.ReadXml(record.S3.Bucket.Name, record.S3.Object.Key);
				var model = modelBuilder.Build(document);

				imageCacheDispatch.Add(model.Parts);

				await _fileSystem.WriteJson(new FileData<BsxModel>
				{
					Drive = record.S3.Bucket.Name,
					FullPath = $"models/{model.Name}.json",
					Content = model
				});
			}

			await imageCacheDispatch.Dispatch();
		}
	}
}
