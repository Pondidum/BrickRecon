using System;
using System.Threading.Tasks;
using Amazon.Lambda;
using Amazon.Lambda.Core;
using Amazon.Lambda.S3Events;
using Amazon.Lambda.SNSEvents;
using Amazon.S3;
using BsxProcessor.Infrastructure;
using JsonSerializer = Amazon.Lambda.Serialization.Json.JsonSerializer;

namespace BsxProcessor
{
	public class Handler
	{
		private readonly S3Handler _s3Handler;
		private readonly SnsHandler _snsHandler;

		public Handler()
		{
			var config = Config.FromEnvironment();
			var lambdaClient = new AmazonLambdaClient();

			var imageCacheDispatch = new ImageCacheDispatcher(config, req => lambdaClient.InvokeAsync(req));
			var modelBuilder = new BsxModelBuilder();

			var fileSystem = new S3FileSystem(new AmazonS3Client());
			var bsxProcessor = new BsxProcessor(fileSystem, imageCacheDispatch, modelBuilder);

			_s3Handler = new S3Handler(fileSystem, bsxProcessor);
			_snsHandler = new SnsHandler(bsxProcessor);
		}

		[LambdaSerializer(typeof(JsonSerializer))]
		public async Task FromS3(S3Event s3Event) => await _s3Handler.Handle(s3Event);

		[LambdaSerializer(typeof(JsonSerializer))]
		public async Task FromSns(SNSEvent snsEvent) => await _snsHandler.Handle(snsEvent);
	}
}
