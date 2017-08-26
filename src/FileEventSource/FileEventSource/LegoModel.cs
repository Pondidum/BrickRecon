using System;
using System.Collections.Generic;
using FileEventSource.Events;
using FileEventSource.Infrastructure;

namespace FileEventSource
{
	public class LegoModel : AggregateRoot
	{
		public string Title { get; private set; }
		public string Name { get; private set; }
		public string Author { get; private set; }
		public IEnumerable<string> Comments => _comments;
		public IEnumerable<object> Parts => _parts;

		private readonly List<string> _comments;
		private readonly List<object> _parts;
		private readonly List<CommandAdded> _unhandledCommands;

		public LegoModel()
		{
			_comments = new List<string>();
			_parts = new List<object>();
			_unhandledCommands = new List<CommandAdded>();

			Register<CommentAdded>(Apply);
			Register<CommandAdded>(Apply);
		}

		public void AddComment(string comment)
		{
			Console.WriteLine($"AddComment: {comment}");

			if (string.IsNullOrWhiteSpace(comment))
				return;

			ApplyEvent(new CommentAdded(comment));
		}

		public void AddCommand(string command, string arguments)
		{
			Console.WriteLine($"AddCommand: {command}, {arguments}");
			if (string.IsNullOrWhiteSpace(command))
				return;

			ApplyEvent(new CommandAdded(command, arguments));
		}

		private void Apply(CommentAdded @event)
		{
			_comments.Add(@event.Comment);
		}

		private void Apply(CommandAdded @event)
		{
			var actions = new Dictionary<string, Action>(StringComparer.OrdinalIgnoreCase)
			{
				{ "name", () => Name = @event.Arguments },
				{ "author", () => Author = @event.Arguments },
				{ "title", () => Title = @event.Arguments }
			};

			if (actions.ContainsKey(@event.Command))
				actions[@event.Command]();
			else
				_unhandledCommands.Add(@event);
		}
	}
}
