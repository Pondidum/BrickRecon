using System.Threading.Tasks;
using Amazon.S3;
using Amazon.S3.Model;
using Newtonsoft.Json;

namespace BsxProcessor
{
	public class FileWriter
	{
		public async Task Write(string bucket, string key, object contents)
		{
			var client = new AmazonS3Client();

			var response = await client.PutObjectAsync(new PutObjectRequest
			{
				BucketName = bucket,
				Key = key,
				ContentType = "application/json",
				ContentBody = JsonConvert.SerializeObject(contents)
			});
		}
	}
}
