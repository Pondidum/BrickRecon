using System.Collections.Generic;
using FileEventSource.Events;

namespace FileEventSource
{
	public class LegoModel
	{
		public IEnumerable<object> Events => _events;

		private readonly List<object> _events;

		public LegoModel()
		{
			_events = new List<object>();
		}

		private void Apply(object @event)
		{
			_events.Add(@event);
		}

		public void AddComment(string comment) => Apply(new CommentAdded(comment));
		public void AddRotation() => Apply(new RotationAdded());
	}
}