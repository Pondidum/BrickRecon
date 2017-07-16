using System;
using System.Collections.Generic;
using System.IO;
using System.Threading.Tasks;
using System.Xml.Linq;
using Amazon.Lambda.Core;
using Amazon.Lambda.S3Events;
using Amazon.Lambda.Serialization.Json;
using Amazon.S3;
using Amazon.S3.Model;
using Amazon.S3.Util;

namespace BsxProcessor
{
	public class Handler
	{
		[LambdaSerializer(typeof(JsonSerializer))]
		public void Handle(S3Event s3Event)
		{
			HandleRecords(s3Event.Records).Wait();
		}

		private async Task HandleRecords(IEnumerable<S3EventNotification.S3EventNotificationRecord> records)
		{
			var reader = new FileReader();
			var writer = new FileWriter();
			var modelBuilder = new BsxModelBuilder();

			foreach (var record in records)
			{
				if (ShouldHandle(record) == false)
					continue;

				var document = await reader.Read(record.S3.Bucket.Name, record.S3.Object.Key);
				var model = modelBuilder.Build(document);

				var outputPath = "website/models/" + Path.GetFileNameWithoutExtension(record.S3.Object.Key) + ".json";

				await writer.Write(record.S3.Bucket.Name, outputPath, model);
			}
		}

		private static bool ShouldHandle(S3EventNotification.S3EventNotificationRecord record)
		{
			var path = record.S3.Object.Key;

			if (path.StartsWith("upload/") == false)
				return false;

			return string.Equals(Path.GetExtension(path), "bsx", StringComparison.OrdinalIgnoreCase);
		}
	}

	public class FileWriter
	{
		public async Task Write(string bucket, string key, object contents)
		{
			var client = new AmazonS3Client();

			using (var ms = new MemoryStream())
			{
				new JsonSerializer().Serialize(contents, ms);
				ms.Position = 0;

				var response = await client.PutObjectAsync(new PutObjectRequest
				{
					BucketName = bucket,
					Key = key,
					ContentType = "application/json",
					InputStream = ms
				});
			}
		}
	}

	public class FileReader
	{
		public async Task<XDocument> Read(string bucket, string key)
		{
			var client = new Amazon.S3.AmazonS3Client();
			var response = await client.GetObjectAsync(new GetObjectRequest
			{
				BucketName = bucket,
				Key = key
			});

			using (var stream = response.ResponseStream)
				return XDocument.Load(stream);
		}
	}

	public class BsxModel
	{
		public IEnumerable<Part> Parts { get; set; }
	}

	public class Part
	{
		public int PartNumber { get; set; }
		public string Name { get; set; }
		public Colors Color { get; set; }
		public int Quantity { get; set; }
		public string Category { get; set; }
	}
}
