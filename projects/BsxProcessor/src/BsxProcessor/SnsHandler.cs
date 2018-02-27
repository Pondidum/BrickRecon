using System.Linq;
using System.Threading.Tasks;
using System.Xml.Linq;
using Amazon.Lambda.SNSEvents;
using BsxProcessor.Infrastructure;
using Newtonsoft.Json;

namespace BsxProcessor
{
	public class SnsHandler
	{
		private readonly BsxProcessor _bsxProcessor;

		public SnsHandler(BsxProcessor bsxProcessor)
		{
			_bsxProcessor = bsxProcessor;
		}

		public async Task Handle(SNSEvent snsEvent)
		{
			var files = snsEvent
				.Records
				.Select(record => record.Sns)
				.Where(IsProcessBsxRequest)
				.Select(record => JsonConvert.DeserializeObject<FileData<XDocument>>(record.Message));

			await _bsxProcessor.Execute(files);
		}

		private static bool IsProcessBsxRequest(SNSEvent.SNSMessage record) =>
			record.MessageAttributes.TryGetValue("EventType", out var attribute) &&
			attribute.Value == "PROCESS_BSX_REQUEST";
	}
}
