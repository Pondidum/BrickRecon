using System.Collections.Generic;
using System.Threading.Tasks;
using Amazon.Lambda;
using Amazon.Lambda.Core;
using Amazon.Lambda.S3Events;
using Amazon.Lambda.Serialization.Json;
using Amazon.S3;
using Amazon.S3.Util;
using Environment = System.Environment;

namespace BsxProcessor
{
	public class Handler
	{
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
			var s3Client = new AmazonS3Client();

			var reader = new FileReader();
			var writer = new FileWriter(s3Client);
			var modelBuilder = new BsxModelBuilder();

			var imageCacheDispatch = new ImageCacheDispatcher(config, req => lambdaClient.InvokeAsync(req));

			foreach (var record in records)
			{
				var document = await reader.Read(record.S3.Bucket.Name, record.S3.Object.Key);
				var model = modelBuilder.Build(record.S3.Object.Key, document);

				imageCacheDispatch.Add(model.Parts);

				await writer.Write(record.S3.Bucket.Name, $"models/{model.Name}.json", model);
			}

			await imageCacheDispatch.Dispatch();
		}
	}
}
