using System.Collections.Generic;
using System.IO;
using System.Linq;
using System.Threading.Tasks;
using System.Xml.Linq;
using BsxProcessor.Domain;
using BsxProcessor.Infrastructure;
using NSubstitute;
using Shouldly;
using Xunit;

namespace BsxProcessor.Tests
{
	public class BsxProcessorTests
	{
		private const string BucketName = "TestBucket";

		private readonly BsxProcessor _handler;
		private readonly IImageCacheDispatcher _imageCacheDispatcher;
		private readonly InMemoryFileSystem _fileSystem;
		private readonly IBsxModelBuilder _modelBuilder;
		private readonly Config _config;

		public BsxProcessorTests()
		{
			_fileSystem = new InMemoryFileSystem();
			_imageCacheDispatcher = Substitute.For<IImageCacheDispatcher>();
			_modelBuilder = Substitute.For<IBsxModelBuilder>();
			_config = new Config
			{
				OutputBucketPath = BucketName + "://models/"
			};

			_modelBuilder
				.Build(Arg.Any<FileData<XDocument>>())
				.Returns(ci => new BsxModelBuilder().Build(ci.Arg<FileData<XDocument>>()));

			_handler = new BsxProcessor(_fileSystem, _config, _imageCacheDispatcher, _modelBuilder);
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
		public async Task When_there_is_one_non_existing_record()
		{
			var records = new[]
			{
				new FileData<XDocument> { Drive = BucketName, FullPath = "one.bsx" },
			};

			await _handler.Execute(records);

			_modelBuilder.DidNotReceive().Build(Arg.Any<FileData<XDocument>>());
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

			await _imageCacheDispatcher.Received(1).Dispatch();

			var written = _fileSystem.Writes.OfType<FileData<BsxModel>>().Single();

			written.ShouldSatisfyAllConditions(
				() => written.Drive.ShouldBe(BucketName, StringCompareShould.IgnoreCase),
				() => written.FullPath.ShouldBe("models/one.json")
			);
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

			await _imageCacheDispatcher.Received(1).Dispatch();

			var written = _fileSystem.Writes.OfType<FileData<BsxModel>>().Select(f => f.FullPath).ToArray();

			written.ShouldBe(new[]
			{
				"models/one.json",
				"models/two.json"
			});
		}
	}
}
