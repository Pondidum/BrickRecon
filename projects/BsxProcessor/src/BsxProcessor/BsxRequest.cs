using System.Xml.Linq;

namespace BsxProcessor
{
	public class BsxRequest
	{
		public XDocument Content { get; set; }
		public string ModelName { get; set; }
	}
}
