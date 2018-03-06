using System;
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
				OutputBucketPath = new Uri($"s3://{BucketName}/models/")
			};

			_modelBuilder
				.Build(Arg.Any<string>(), Arg.Any<XDocument>())
				.Returns(ci => new BsxModelBuilder().Build(ci.Arg<string>(), ci.Arg<XDocument>()));

			_handler = new BsxProcessor(_fileSystem, _config, _imageCacheDispatcher, _modelBuilder);
		}

		private static BsxRequest CreateFile(string path, string data) => new BsxRequest
		{
			ModelName = Path.GetFileNameWithoutExtension(path),
			Content = XDocument.Parse(data)
		};

		[Fact]
		public async Task When_there_are_no_records_to_process()
		{
			var records = Enumerable.Empty<BsxRequest>();

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
