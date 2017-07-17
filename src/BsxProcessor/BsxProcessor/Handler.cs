using System.Collections.Generic;
using System.IO;
using System.Threading.Tasks;
using Amazon.Lambda.Core;
using Amazon.Lambda.S3Events;
using Amazon.Lambda.Serialization.Json;
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
				var document = await reader.Read(record.S3.Bucket.Name, record.S3.Object.Key);
				var model = modelBuilder.Build(document);

				var outputPath = "models/" + Path.GetFileNameWithoutExtension(record.S3.Object.Key) + ".json";

				await writer.Write(record.S3.Bucket.Name, outputPath, model);
			}
		}
	}
}
