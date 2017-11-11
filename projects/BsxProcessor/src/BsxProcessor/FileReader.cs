using System.Threading.Tasks;
using System.Xml.Linq;
using Amazon.S3;
using Amazon.S3.Model;

namespace BsxProcessor
{
	public class FileReader
	{
		public async Task<XDocument> Read(string bucket, string key)
		{
			var client = new AmazonS3Client();
			var response = await client.GetObjectAsync(new GetObjectRequest
			{
				BucketName = bucket,
				Key = key
			});

			using (var stream = response.ResponseStream)
				return XDocument.Load(stream);
		}
	}
}
