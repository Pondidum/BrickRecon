using System.Xml.Linq;
using BsxProcessor.Domain;
using BsxProcessor.Infrastructure;

namespace BsxProcessor
{
	public interface IBsxModelBuilder
	{
		BsxModel Build(FileData<XDocument> file);
	}
}
