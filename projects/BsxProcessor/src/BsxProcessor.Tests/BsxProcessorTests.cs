using System.Collections.Generic;
using System.Linq;
using System.Threading.Tasks;
using System.Xml.Linq;
using BsxProcessor.Domain;
using BsxProcessor.Infrastructure;
using NSubstitute;
using Xunit;

namespace BsxProcessor.Tests
{
	public class BsxProcessorTests
	{
		private const string BucketName = "TestBucket";

		private readonly BsxProcessor _handler;
		private readonly IImageCacheDispatcher _imageCacheDispatcher;
		private readonly IFileSystem _fileSystem;

		public BsxProcessorTests()
		{
			_fileSystem = Substitute.For<IFileSystem>();
			_imageCacheDispatcher = Substitute.For<IImageCacheDispatcher>();
			var modelBuilder = new BsxModelBuilder();

			_handler = new BsxProcessor(_fileSystem, _imageCacheDispatcher, modelBuilder);
		}

		private static FileData<XDocument> CreateFile(string path, string data) => new FileData<XDocument>
		{
			Drive = BucketName,
			FullPath = path,
			Exists = true,
			Content = XDocument.Parse(data)
		};

		[Fact]
		public async Task When_there_are_no_records_to_process()
		{
			var records = Enumerable.Empty<FileData<XDocument>>();

			await _handler.Execute(records);

			await _imageCacheDispatcher.Received(1).Dispatch();
		}

		[Fact]
		public async Task When_handling_one_record()
		{
			var records = new[]
			{
				CreateFile("one.bsx", TestData.BsxWithTwoParts)
			};

			await _handler.Execute(records);

			_imageCacheDispatcher.Received(1).Add(Arg.Any<IEnumerable<Part>>());

			await _fileSystem.Received().WriteJson(Arg.Is<FileData<BsxModel>>(arg => arg.Drive == BucketName && arg.FullPath == "models/one.json"));
			await _imageCacheDispatcher.Received(1).Dispatch();
		}

		[Fact]
		public async Task When_handling_multiple_records()
		{
			var records = new[]
			{
				CreateFile("one.bsx", TestData.BsxWithTwoParts),
				CreateFile("two.bsx", TestData.BsxWithFourParts)
			};

			await _handler.Execute(records);

			_imageCacheDispatcher.Received(2).Add(Arg.Any<IEnumerable<Part>>());

			await _fileSystem.Received().WriteJson(Arg.Is<FileData<BsxModel>>(arg => arg.Drive == BucketName && arg.FullPath == "models/one.json"));
			await _fileSystem.Received().WriteJson(Arg.Is<FileData<BsxModel>>(arg => arg.Drive == BucketName && arg.FullPath == "models/two.json"));
			await _imageCacheDispatcher.Received(1).Dispatch();
		}
	}
}
