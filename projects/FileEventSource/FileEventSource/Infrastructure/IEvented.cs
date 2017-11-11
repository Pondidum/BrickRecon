using System.Collections.Generic;

namespace FileEventSource.Infrastructure
{
	public interface IEvented
	{
		IEnumerable<object> GetPendingEvents();
		void ClearPendingEvents();
		void LoadFromEvents(IEnumerable<object> events);
	}
}
