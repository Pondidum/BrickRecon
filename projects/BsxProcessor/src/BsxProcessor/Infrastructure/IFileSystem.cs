using System.Threading.Tasks;
using System.Xml.Linq;

namespace BsxProcessor.Infrastructure
{
	public interface IFileSystem
	{
		Task<FileData<XDocument>> ReadXml(string drive, string path);
		Task WriteJson<TContent>(FileData<TContent> file);
	}
}
