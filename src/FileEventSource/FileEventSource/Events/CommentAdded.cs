using FileEventSource.Infrastructure;

namespace FileEventSource.Events
{
	public class CommentAdded : Event
	{
		public string Comment { get; }

		public CommentAdded(string comment)
		{
			Comment = comment;
		}
	}

	public class CommandAdded : Event
	{
		public string Command { get; }
		public string Arguments { get; }

		public CommandAdded(string command, string arguments)
		{
			Command = command;
			Arguments = arguments;
		}
	}
}