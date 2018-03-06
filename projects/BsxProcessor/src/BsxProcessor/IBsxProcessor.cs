using System.Collections.Generic;
using System.Threading.Tasks;

namespace BsxProcessor
{
	public interface IBsxProcessor
	{
		Task Execute(IEnumerable<BsxRequest> records);
	}
}
