using System.Collections.Generic;

namespace BsxProcessor.Domain
{
	public class BsxModel
	{
		public IEnumerable<Part> Parts { get; set; }
		public string Name { get; set; }
	}
}
