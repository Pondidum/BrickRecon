using System;
using System.Collections.Generic;
using System.Linq;
using FileEventSource.LowLevelApi.Lines;

namespace FileEventSource.LowLevelApi
{
	public class LineParser
	{
		private readonly Dictionary<LineTypes, Func<int, string[], Line>> _builders;

		public LineParser()
		{
			_builders = new Dictionary<LineTypes, Func<int, string[], Line>>
			{
				{ LineTypes.CommentOrMeta, BuildMetaOrCommentLine },
				{ LineTypes.SubFileReference, BuildSubReferenceLine }
			};
		}

		public IEnumerable<Line> Parse(IEnumerable<string> lines)
		{
			var lineNumber = 0;
			foreach (var line in lines)
			{
				var tokens = line.Split(' ');

				if (tokens.Length == 0 || string.IsNullOrWhiteSpace(tokens.First()))
					continue;

				var type = (LineTypes)Enum.Parse(typeof(LineTypes), tokens.First());

				Func<int, string[], Line> builder;

				if (!_builders.TryGetValue(type, out builder))
					continue;

				var command = builder(lineNumber++, tokens);

				if (command != null)
					yield return command;
			}
		}

		private Line BuildMetaOrCommentLine(int lineNumber, string[] tokens)
		{
			if (tokens.Length == 1)
				return null;

			var command = tokens[1].TrimEnd(':');

			if (command.EqualsIgnore("Name"))
				return new NameLine(string.Join(" ", tokens.Skip(2)));

			if (command.EqualsIgnore("Author"))
				return new AuthorLine(string.Join(" ", tokens.Skip(2)));

			if (LineCommands.Official.Contains(command) || LineCommands.Unofficial.Contains(command))
				return new CommandLine(command, string.Join(" ", tokens.Skip(2)));

			if (lineNumber == 0)
				return new TitleLine(command);

			return new CommentLine(string.Join(" ", tokens.Skip(1)));
		}

		private Line BuildSubReferenceLine(int lineNumber, string[] tokens)
		{
			if (tokens.Length == 1)
				return null;

			var part = tokens.Last();

			return new PartLine(part); // this will clearly need a lot more
		}
	}
}
