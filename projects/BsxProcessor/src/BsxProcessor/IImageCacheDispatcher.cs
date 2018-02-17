using System.Collections.Generic;
using System.Threading.Tasks;
using BsxProcessor.Domain;

namespace BsxProcessor
{
	public interface IImageCacheDispatcher
	{
		void Add(IEnumerable<Part> parts);
		Task Dispatch();
	}
}
