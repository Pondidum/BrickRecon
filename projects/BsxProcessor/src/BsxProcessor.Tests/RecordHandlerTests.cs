using System.Collections.Generic;
using System.Linq;
using System.Threading.Tasks;
using System.Xml.Linq;
using Amazon.S3.Util;
using BsxProcessor.Domain;
using BsxProcessor.Infrastructure;
using NSubstitute;
using Xunit;

namespace BsxProcessor.Tests
{
	public class RecordHandlerTests
	{
		private const string BucketName = "TestBucket";

		private readonly RecordHandler _handler;
		private readonly IImageCacheDispatcher _imageCacheDispatcher;
		private readonly IFileSystem _fileSystem;

		public RecordHandlerTests()
		{
			_fileSystem = Substitute.For<IFileSystem>();
			_imageCacheDispatcher = Substitute.For<IImageCacheDispatcher>();
			var modelBuilder = new BsxModelBuilder();

			_handler = new RecordHandler(_fileSystem, _imageCacheDispatcher, modelBuilder);
		}

		private static IEnumerable<S3EventNotification.S3EventNotificationRecord> CreateRecords(params string[] keys) => keys.Select(key => new S3EventNotification.S3EventNotificationRecord
		{
			S3 = new S3EventNotification.S3Entity
			{
				Bucket = new S3EventNotification.S3BucketEntity { Name = BucketName },
				Object = new S3EventNotification.S3ObjectEntity { Key = key }
			}
		});

		private void CreateFile(string key, string data)
		{
			_fileSystem.ReadXml(BucketName, key).Returns(new FileData<XDocument>
			{
				Drive = BucketName,
				FullPath = key,
				Exists = true,
				Content = XDocument.Parse(data)
			});
		}

		[Fact]
		public async Task When_there_are_no_records_to_process()
		{
			var records = CreateRecords();

			await _handler.Execute(records);

			await _imageCacheDispatcher.Received(1).Dispatch();
		}

		[Fact]
		public async Task When_handling_one_record()
		{
			var records = CreateRecords("one.bsx");
			CreateFile("one.bsx", TestData.BsxWithTwoParts);

			await _handler.Execute(records);

			_imageCacheDispatcher.Received(1).Add(Arg.Any<IEnumerable<Part>>());

			await _fileSystem.Received().WriteJson(Arg.Is<FileData<BsxModel>>(arg => arg.Drive == BucketName && arg.FullPath == "models/one.json"));
			await _imageCacheDispatcher.Received(1).Dispatch();
		}

		[Fact]
		public async Task When_handling_multiple_records()
		{
			var records = CreateRecords("one.bsx", "two.bsx");
			CreateFile("one.bsx", TestData.BsxWithTwoParts);
			CreateFile("two.bsx", TestData.BsxWithFourParts);

			await _handler.Execute(records);

			_imageCacheDispatcher.Received(2).Add(Arg.Any<IEnumerable<Part>>());

			await _fileSystem.Received().WriteJson(Arg.Is<FileData<BsxModel>>(arg => arg.Drive == BucketName && arg.FullPath == "models/one.json"));
			await _fileSystem.Received().WriteJson(Arg.Is<FileData<BsxModel>>(arg => arg.Drive == BucketName && arg.FullPath == "models/two.json"));
			await _imageCacheDispatcher.Received(1).Dispatch();
		}
	}
}
