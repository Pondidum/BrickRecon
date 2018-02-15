using System.Threading.Tasks;
using Amazon.S3;
using Amazon.S3.Model;
using BsxProcessor.Infrastructure;
using NSubstitute;
using Xunit;

namespace BsxProcessor.Tests.Infrastructure
{
	public class S3FileSystemTests
	{
		private readonly IAmazonS3 _client;
		private readonly IFileSystem _fileSystem;

		public S3FileSystemTests()
		{
			_client = Substitute.For<IAmazonS3>();
			_fileSystem = new S3FileSystem(_client);
		}

		[Fact]
		public async Task When_serializing_camal_case_is_used()
		{
			var contents = new { SomeProperty = "nothing" };
			var expected = "{\"someProperty\":\"nothing\"}";

			await _fileSystem.WriteJson(new FileData<object>
			{
				Drive = "",
				FullPath = "",
				Content = contents
			});

			await _client
				.Received()
				.PutObjectAsync(Arg.Is<PutObjectRequest>(req => req.ContentBody == expected));
		}
	}
}
