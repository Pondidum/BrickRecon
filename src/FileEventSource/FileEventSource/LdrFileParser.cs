using System.IO;

namespace FileEventSource
{
	public class LdrFileParser
	{
		public LegoModel Parse(Stream stream)
		{
			using (var reader = new StreamReader(stream))
			{
				var model = new LegoModel();

				string line = null;
				while ((line = reader.ReadLine()) != null)
				{
					var tokens = TokensFromLine(line);
				}

				return model;
			}
		}

		private string[] TokensFromLine(string line)
		{
			return line.Split(' ');
		}
	}
}
