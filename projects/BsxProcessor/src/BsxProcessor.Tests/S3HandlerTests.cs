using System.Collections.Generic;
using System.Linq;
using System.Threading.Tasks;
using System.Xml.Linq;
using Amazon.Lambda.S3Events;
using Amazon.S3.Util;
using BsxProcessor.Infrastructure;
using NSubstitute;
using Xunit;

namespace BsxProcessor.Tests
{
	public class S3HandlerTests
	{
		private readonly IBsxProcessor _processor;
		private readonly IFileSystem _fileSystem;
		private readonly S3Handler _handler;

		public S3HandlerTests()
		{
			_fileSystem = Substitute.For<IFileSystem>();
			_processor = Substitute.For<IBsxProcessor>();
			_handler = new S3Handler(_fileSystem, _processor);
		}

		private static S3Event CreateRecords(params string[] keys) => new S3Event
		{
			Records = keys
				.Select(key => new S3EventNotification.S3Entity
				{
					Bucket = new S3EventNotification.S3BucketEntity { Name = "wat" },
					Object = new S3EventNotification.S3ObjectEntity { Key = key }
				})
				.Select(entity => new S3EventNotification.S3EventNotificationRecord { S3 = entity })
				.ToList()
		};

		private void CreateFiles(params string[] files)
		{
			foreach (var file in files)
				_fileSystem.ReadXml("wat", file).Returns(new FileData<XDocument>
				{
					Drive = "wat",
					FullPath = file,
					Content = XDocument.Parse("<nope />"),
					Exists = true
				});
		}

		[Fact]
		public async Task When_there_are_no_records()
		{
			var records = CreateRecords();

			await _handler.Handle(records);

			await _processor
				.Received()
				.Execute(Arg.Is<IEnumerable<FileData<XDocument>>>(e => e.Any() == false));
		}

		[Fact]
		public async Task When_there_is_one_record()
		{
			var records = CreateRecords("first");
			CreateFiles("first");

			await _handler.Handle(records);

			await _processor
				.Received()
				.Execute(Arg.Is<IEnumerable<FileData<XDocument>>>(e => e.Single().FullPath == "first"));
		}
		
		[Fact]
		public async Task When_there_are_multiple_records()
		{
			var records = CreateRecords("first", "second", "third");
			CreateFiles("first", "second", "third");

			await _handler.Handle(records);

			await _processor
				.Received()
				.Execute(Arg.Is<IEnumerable<FileData<XDocument>>>(e => e.Count() == 3));
		}
	}
}
