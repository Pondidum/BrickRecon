using System.IO;
using System.Linq;
using System.Text;
using FileEventSource.LowLevelApi;
using Shouldly;
using Xunit;

namespace FileEventSource.Tests.LowLevelApi
{
	public class LineParserTests
	{
		private const string SingleModelFile =
			@"0 Untitled
0 Name: dual-rail-gun.ldr
0 Author: LDraw
0 Unofficial Model
0 ROTATION CENTER 0 0 0 1 ""Custom""
0 ROTATION CONFIG 0 0
1 0 0 0 0 -1 0 0 0 -1 0 0 0 1 4595.dat
1 0 0 8 10 1 0 0 0 -1 0 0 0 -1 2555.dat
1 0 0 8 -10 0 0 1 0 -1 0 1 0 0 2555.dat
1 0 0 10 -16 1 0 0 0 1 0 0 0 1 2432.dat
1 0 -18 0 10 0 1 0 1 0 0 0 0 -1 3069b.dat
1 0 18 0 10 0 -1 0 1 0 0 0 0 1 3069b.dat
1 0 0 -40 0 1 0 0 0 1 0 0 0 1 4595.dat
1 0 18 -30 -10 0 -1 0 1 0 0 0 0 1 2420.dat
1 0 -18 -30 -10 0 1 0 0 0 1 1 0 0 2420.dat
1 0 -42 -10 0 0 1 0 0 0 -1 -1 0 0 4865a.dat
1 0 42 -10 0 0 -1 0 0 0 -1 1 0 0 4865a.dat
1 0 26 -30 0 0 -1 0 0 0 1 -1 0 0 4590.dat
1 0 -26 -30 0 0 1 0 0 0 1 1 0 0 4590.dat
1 0 0 -48 10 0 0 -1 0 1 0 1 0 0 3024.dat
1 0 0 -48 -10 1 0 0 0 1 0 0 0 1 6019.dat
1 0 0 -56 0 0 0 1 0 1 0 -1 0 0 3069b.dat
1 0 0 -45 -29 -1 0 0 0 0 1 0 1 0 30031.dat
1 0 0 -30 -28 -1 0 0 0 0 1 0 1 0 63965.dat
1 0 0 -10 -28 -1 0 0 0 0 1 0 1 0 63965.dat
0
";

		[Fact]
		public void When_parsing_a_single_model_file()
		{
			var parser = new LineParser();
			var lines = parser.Parse(SingleModelFile.Split('\n'));

			lines.Select(line => line.GetType()).ShouldBe(new[]
			{
				typeof(TitleLine),
				typeof(CommandLine),
				typeof(CommandLine),
				typeof(CommentLine),
				typeof(CommandLine),
				typeof(CommandLine)
			});
		}
	}
}
