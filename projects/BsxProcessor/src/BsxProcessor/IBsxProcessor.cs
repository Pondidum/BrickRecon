using System.Collections.Generic;
using System.Threading.Tasks;
using System.Xml.Linq;
using BsxProcessor.Infrastructure;

namespace BsxProcessor
{
	public interface IBsxProcessor
	{
		Task Execute(IEnumerable<FileData<XDocument>> records);
	}
}
