using System;
using System.Collections.Generic;
using System.Linq;

namespace FileEventSource.LowLevelApi
{
	public class LineParser
	{
		public IEnumerable<Line> Parse(IEnumerable<string> lines)
		{
			var lineNumber = 0;
			foreach (var line in lines)
			{
				var tokens = TokensFromLine(line);

				if (tokens.Length == 0 || string.IsNullOrWhiteSpace(tokens.First()))
					continue;

				var type = (LineTypes)Enum.Parse(typeof(LineTypes), tokens.First());

				Line command = null;
				if (type == LineTypes.CommentOrMeta)
					command = BuildMetaOrCommentLine(lineNumber, tokens);

				lineNumber++;

				if (command == null)
					continue;

				yield return command;
			}
		}

		private string[] TokensFromLine(string line)
		{
			return line.Split(' ');
		}

		private Line BuildMetaOrCommentLine(int lineNumber, string[] tokens)
		{
			if (tokens.Length == 1)
				return null;

			var command = tokens[1].TrimEnd(':');

			if (LineCommands.Official.Contains(command) || LineCommands.Unofficial.Contains(command))
				return new CommandLine(command, string.Join(" ", tokens.Skip(2)));

			if (lineNumber == 0)
				return new TitleLine(command);

			return new CommentLine(string.Join(" ", tokens.Skip(1)));
		}
	}

	public class CommentLine : Line
	{
		public string Comment { get; }

		public CommentLine(string comment)
		{
			Comment = comment;
		}
	}

	public class TitleLine : Line
	{
		public string Title { get; }

		public TitleLine(string title)
		{
			Title = title;
		}
	}

	public class CommandLine : Line
	{
		public string Command { get; }
		public string Arguments { get; }

		public CommandLine(string command, string arguments)
		{
			Command = command;
			Arguments = arguments;
		}
	}

	public class Line
	{
	}
}
