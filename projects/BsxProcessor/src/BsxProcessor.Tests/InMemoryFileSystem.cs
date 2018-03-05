using System.Collections.Generic;
using System.Threading.Tasks;
using System.Xml.Linq;
using BsxProcessor.Infrastructure;

namespace BsxProcessor.Tests
{
	public class InMemoryFileSystem : IFileSystem
	{
		public IEnumerable<object> Writes => _writes;

		private readonly List<object> _writes;

		public InMemoryFileSystem()
		{
			_writes = new List<object>();
		}

		public Task<FileData<XDocument>> ReadXml(string drive, string path)
		{
			throw new System.NotImplementedException();
		}

		public Task WriteJson<TContent>(FileData<TContent> file)
		{
			_writes.Add(file);
			return Task.CompletedTask;
		}
	}
}
