using System;
using System.Collections.Generic;
using System.IO;
using System.Linq;

namespace FileEventSource
{
	public class MpdFileParser
	{
		public LegoModel Parse(Stream stream)
		{
			using (var reader = new StreamReader(stream))
			{
				var model = new LegoModel();

				int lineNumber = 0;
				string line = null;
				while ((line = reader.ReadLine()) != null)
				{
					var tokens = TokensFromLine(line);

					if (tokens.Length == 0)
						continue;

					var type = (LineTypes)Enum.Parse(typeof(LineTypes), tokens.First());

					if (type == LineTypes.CommentOrMeta)
						HandleCommentOrMeta(model, lineNumber, tokens);


					lineNumber++;
				}

				return model;
			}
		}

		private string[] TokensFromLine(string line)
		{
			return line.Split(' ');
		}

		private void HandleCommentOrMeta(LegoModel model, int lineNumber, string[] tokens)
		{
			if (tokens.Length == 1)
				return;

			var command = tokens[1].TrimEnd(':');

			Console.WriteLine(command);

			if (LineCommands.Official.Contains(command) || LineCommands.Unofficial.Contains(command))
			{
				//command
				model.AddCommand(command, string.Join(" ", tokens.Skip(2)));
			}
			else if (lineNumber == 0)
			{
				model.AddCommand("title", command);
			}
			else
			{
				//comment
				model.AddComment(string.Join(" ", tokens.Skip(1)));
			}
		}
	}
}
