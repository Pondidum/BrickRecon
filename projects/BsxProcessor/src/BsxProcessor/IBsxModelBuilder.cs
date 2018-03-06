using System.Xml.Linq;
using BsxProcessor.Domain;

namespace BsxProcessor
{
	public interface IBsxModelBuilder
	{
		BsxModel Build(string modelName, XDocument content);
	}
}
