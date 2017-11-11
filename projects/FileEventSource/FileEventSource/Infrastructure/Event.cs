using System;

namespace FileEventSource.Infrastructure
{
	public abstract class Event
	{
		public Guid AggregateID { get; set; }
		public DateTime TimeStamp { get; set; }
	}
}
