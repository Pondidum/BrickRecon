using System;
using System.Collections.Generic;
using System.Linq;
using System.Threading.Tasks;
using System.Xml.Linq;
using Amazon.Lambda;
using Amazon.Lambda.Core;
using Amazon.Lambda.S3Events;
using Amazon.Lambda.Serialization.Json;
using Amazon.Lambda.SNSEvents;
using Amazon.S3;
using BsxProcessor.Infrastructure;

namespace BsxProcessor
{
	public class Handler
	{
		private readonly RecordHandler _recordHandler;
		private readonly S3FileSystem _fileSystem;

		public Handler()
		{
			var config = Config.FromEnvironment();
			var lambdaClient = new AmazonLambdaClient();

			var imageCacheDispatch = new ImageCacheDispatcher(config, req => lambdaClient.InvokeAsync(req));
			var modelBuilder = new BsxModelBuilder();

			_fileSystem = new S3FileSystem(new AmazonS3Client());
			_recordHandler = new RecordHandler(_fileSystem, imageCacheDispatch, modelBuilder);
		}

		[LambdaSerializer(typeof(JsonSerializer))]
		public async Task FromS3(S3Event s3Event)
		{
			var files = new List<FileData<XDocument>>(s3Event.Records.Count);

			foreach (var record in s3Event.Records)
				files.Add(await _fileSystem.ReadXml(record.S3.Bucket.Name, record.S3.Object.Key));

			await _recordHandler.Execute(files);
		}

		[LambdaSerializer(typeof(JsonSerializer))]
		public void FromSns(SNSEvent snsEvent)
		{
			throw new NotImplementedException();
		}
	}
}
