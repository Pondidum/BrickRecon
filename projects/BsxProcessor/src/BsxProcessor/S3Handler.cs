using System.Collections.Generic;
using System.Threading.Tasks;
using System.Xml.Linq;
using Amazon.Lambda.S3Events;
using BsxProcessor.Infrastructure;

namespace BsxProcessor
{
	public class S3Handler
	{
		private readonly IFileSystem _fileSystem;
		private readonly BsxProcessor _bsxProcessor;

		public S3Handler(IFileSystem fileSystem, BsxProcessor bsxProcessor)
		{
			_fileSystem = fileSystem;
			_bsxProcessor = bsxProcessor;
		}

		public async Task Handle(S3Event s3Event)
		{
			var files = new List<FileData<XDocument>>(s3Event.Records.Count);

			foreach (var record in s3Event.Records)
				files.Add(await _fileSystem.ReadXml(record.S3.Bucket.Name, record.S3.Object.Key));

			await _bsxProcessor.Execute(files);
		}
	}
}
