using System.Net;
using System.Threading.Tasks;
using System.Xml.Linq;
using Amazon.S3;
using Amazon.S3.Model;
using Newtonsoft.Json;
using Newtonsoft.Json.Serialization;

namespace BsxProcessor.Infrastructure
{
	public class S3FileSystem : IFileSystem
	{
		private static readonly JsonSerializerSettings Settings = new JsonSerializerSettings
		{
			ContractResolver = new CamelCasePropertyNamesContractResolver()
		};
		
		private readonly IAmazonS3 _client;

		public S3FileSystem(IAmazonS3 client)
		{
			_client = client;
		}

		public async Task<FileData<XDocument>> ReadXml(string drive, string path)
		{
			var response = await _client.GetObjectAsync(new GetObjectRequest
			{
				BucketName = drive,
				Key = path
			});

			if (response.HttpStatusCode != HttpStatusCode.OK)
				return new FileData<XDocument>
				{
					Drive = drive,
					FullPath = path,
					Exists = false,
				};

			using (var stream = response.ResponseStream)
				return new FileData<XDocument>
				{
					Drive = drive,
					FullPath = path,
					Content = XDocument.Load(stream),
					Exists = true
				};
		}
		
		public async Task WriteJson<TContent>(FileData<TContent> file)
		{
			await _client.PutObjectAsync(new PutObjectRequest
			{
				BucketName = file.Drive,
				Key = file.FullPath,
				ContentType = "application/json",
				ContentBody = JsonConvert.SerializeObject(file.Content, Settings)
			});
		}
	}
}
