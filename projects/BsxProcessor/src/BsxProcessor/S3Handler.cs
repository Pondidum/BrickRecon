using System.Collections.Generic;
using System.IO;
using System.Threading.Tasks;
using System.Xml.Linq;
using Amazon.Lambda.S3Events;
using BsxProcessor.Infrastructure;

namespace BsxProcessor
{
	public class S3Handler
	{
		private readonly IFileSystem _fileSystem;
		private readonly IBsxProcessor _bsxProcessor;

		public S3Handler(IFileSystem fileSystem, IBsxProcessor bsxProcessor)
		{
			_fileSystem = fileSystem;
			_bsxProcessor = bsxProcessor;
		}

		public async Task Handle(S3Event s3Event)
		{
			var files = new List<BsxRequest>(s3Event.Records.Count);

			foreach (var record in s3Event.Records)
				files.Add(new BsxRequest
				{
					ModelName = Path.GetFileNameWithoutExtension(record.S3.Object.Key),
					Content = (await _fileSystem.ReadXml(record.S3.Bucket.Name, record.S3.Object.Key)).Content
				});

			await _bsxProcessor.Execute(files);
		}
	}
}
