using System.Threading.Tasks;
using Amazon.S3;
using Amazon.S3.Model;
using Newtonsoft.Json;
using Newtonsoft.Json.Serialization;

namespace BsxProcessor
{
	public class FileWriter
	{
		private static readonly JsonSerializerSettings Settings = new JsonSerializerSettings
		{
			ContractResolver = new CamelCasePropertyNamesContractResolver()
		};

		private readonly IAmazonS3 _client;

		public FileWriter(IAmazonS3 client)
		{
			_client = client;
		}

		public async Task Write(string bucket, string key, object contents)
		{
			await _client.PutObjectAsync(new PutObjectRequest
			{
				BucketName = bucket,
				Key = key,
				ContentType = "application/json",
				ContentBody = JsonConvert.SerializeObject(contents, Settings)
			});
		}
	}
}
