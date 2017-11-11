using System.Threading.Tasks;
using Amazon.S3;
using Amazon.S3.Model;
using NSubstitute;
using Xunit;

namespace BsxProcessor.Tests
{
	public class FileWriterTests
	{
		private readonly IAmazonS3 _client;
		private readonly FileWriter _writer;

		public FileWriterTests()
		{
			_client = Substitute.For<IAmazonS3>();
			_writer = new FileWriter(_client);
		}

		[Fact]
		public async Task When_serializing_camal_case_is_used()
		{
			var contents = new { SomeProperty = "nothing" };
			var expeccted = "{\"someProperty\":\"nothing\"}";

			await _writer.Write("", "", contents);

			_client.Received()
				.PutObjectAsync(Arg.Is<PutObjectRequest>(req => req.ContentBody == expeccted));
		}
	}
}
