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

		[Fact]
		public async Task When_there_are_no_records()
		{
			await _handler.Handle(new SNSEvent
			{
				Records = new List<SNSEvent.SNSRecord>()
			});

			await _processor
				.Received()
				.Execute(Arg.Is<IEnumerable<FileData<XDocument>>>(e => e.Any() == false));
		}

		[Fact]
		public async Task When_there_are_only_non_bsx_records_to_process()
		{
			await _handler.Handle(new SNSEvent
			{
				Records = new[]
				{
					new SNSEvent.SNSRecord { Sns = new SNSEvent.SNSMessage
					{
						MessageAttributes = new Dictionary<string, SNSEvent.MessageAttribute>
						{
							{ "EventType", new SNSEvent.MessageAttribute { Value = "WAT" } }
						}
					} }
				}.ToList()
			});

			await _processor
				.Received()
				.Execute(Arg.Is<IEnumerable<FileData<XDocument>>>(e => e.Any() == false));
		}

		[Fact]
		public async Task When_there_is_a_bsx_request_to_process()
		{
			await _handler.Handle(new SNSEvent
			{
				Records = new[]
				{
					new SNSEvent.SNSRecord { Sns = new SNSEvent.SNSMessage
					{
						MessageAttributes = new Dictionary<string, SNSEvent.MessageAttribute>
						{
							{ "EventType", new SNSEvent.MessageAttribute { Value = "PROCESS_BSX_REQUEST" } }
						},
						Message = JsonConvert.SerializeObject(new FileData<XDocument>
						{
							Drive = "wat",
							FullPath = "somefile.bsx",
							Content = XDocument.Parse("<nope />"),
							Exists = true
						})
					} }
				}.ToList()
			});

			await _processor
				.Received()
				.Execute(Arg.Is<IEnumerable<FileData<XDocument>>>(e => e.Single().FullPath == "somefile.bsx"));
		}
	}
}
