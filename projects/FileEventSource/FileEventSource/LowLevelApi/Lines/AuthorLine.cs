namespace FileEventSource.LowLevelApi.Lines
{
	public class AuthorLine : Line
	{
		public string Author { get; }

		public AuthorLine(string author)
		{
			Author = author;
		}
	}
}