using Amazon.Lambda;
using Amazon.Lambda.Core;
using Amazon.Lambda.S3Events;
using Amazon.Lambda.Serialization.Json;
using Amazon.S3;
using BsxProcessor.Infrastructure;

namespace BsxProcessor
{
	public class Handler
	{
		private readonly RecordHandler _recordHandler;

		public Handler()
		{
			var config = Config.FromEnvironment();
			var lambdaClient = new AmazonLambdaClient();

			var fileSystem = new S3FileSystem(new AmazonS3Client());
			var imageCacheDispatch = new ImageCacheDispatcher(config, req => lambdaClient.InvokeAsync(req));
			var modelBuilder = new BsxModelBuilder();
			
			_recordHandler = new RecordHandler(fileSystem, imageCacheDispatch, modelBuilder);
		}

		[LambdaSerializer(typeof(JsonSerializer))]
		public void Handle(S3Event s3Event)
		{
			_recordHandler.Execute(s3Event.Records).Wait();
		}
	}
}
