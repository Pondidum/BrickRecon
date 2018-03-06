using System.Linq;
using System.Threading.Tasks;
using System.Xml.Linq;
using Amazon.Lambda.SNSEvents;

namespace BsxProcessor
{
	public class SnsHandler
	{
		private readonly IBsxProcessor _bsxProcessor;

		public SnsHandler(IBsxProcessor bsxProcessor)
		{
			_bsxProcessor = bsxProcessor;
		}

		public async Task Handle(SNSEvent snsEvent)
		{
			var files = snsEvent
				.Records
				.Select(record => record.Sns)
				.Where(IsProcessBsxRequest)
				.Select(record => new BsxRequest
				{
					ModelName = ReadModelName(record),
					Content  = XDocument.Parse(record.Message)
				});

			await _bsxProcessor.Execute(files);
		}

		private static bool IsProcessBsxRequest(SNSEvent.SNSMessage record) =>
			record.MessageAttributes.TryGetValue("EventType", out var attribute) &&
			attribute.Value == "PROCESS_BSX_REQUEST";

		private static string ReadModelName(SNSEvent.SNSMessage record) =>
			record.MessageAttributes.TryGetValue("ModelName", out var attribute)
				? attribute.Value
				: string.Empty;
	}
}
