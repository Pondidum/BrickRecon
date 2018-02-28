using System.Collections.Generic;
using System.Linq;
using System.Threading.Tasks;
using System.Xml.Linq;
using Amazon.Lambda.SNSEvents;
using BsxProcessor.Infrastructure;
using Newtonsoft.Json;
using NSubstitute;
using Xunit;

namespace BsxProcessor.Tests
{
	public class SnsHandlerTests
	{
		private readonly IBsxProcessor _processor;
		private readonly SnsHandler _handler;

		public SnsHandlerTests()
		{
			_processor = Substitute.For<IBsxProcessor>();
			_handler = new SnsHandler(_processor);
		}

		private static SNSEvent CreateNotification(params SNSEvent.SNSRecord[] records) => new SNSEvent
		{
			Records = records.ToList()
		};

		private static SNSEvent.SNSRecord CreateMessage(string type, object content) => new SNSEvent.SNSRecord
		{
			Sns = new SNSEvent.SNSMessage
			{
				MessageAttributes = new Dictionary<string, SNSEvent.MessageAttribute>
				{
					{ "EventType", new SNSEvent.MessageAttribute { Value = type } }
				},
				Message = JsonConvert.SerializeObject(content)
			}
		};

		[Fact]
		public async Task When_there_are_no_records()
		{
			var notification = CreateNotification();

			await _handler.Handle(notification);

			await _processor
				.Received()
				.Execute(Arg.Is<IEnumerable<FileData<XDocument>>>(e => e.Any() == false));
		}

		[Fact]
		public async Task When_there_are_only_non_bsx_records_to_process()
		{
			var notification = CreateNotification(
				CreateMessage("WAT", null)
			);

			await _handler.Handle(notification);

			await _processor
				.Received()
				.Execute(Arg.Is<IEnumerable<FileData<XDocument>>>(e => e.Any() == false));
		}

		[Fact]
		public async Task When_there_is_a_bsx_request_to_process()
		{
			var notification = CreateNotification(
				CreateMessage("PROCESS_BSX_REQUEST", new FileData<XDocument>
				{
					Drive = "wat",
					FullPath = "somefile.bsx",
					Content = XDocument.Parse("<nope />"),
					Exists = true
				})
			);

			await _handler.Handle(notification);

			await _processor
				.Received()
				.Execute(Arg.Is<IEnumerable<FileData<XDocument>>>(e => e.Single().FullPath == "somefile.bsx"));
		}
	}
}
