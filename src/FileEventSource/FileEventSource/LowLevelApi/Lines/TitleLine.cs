namespace FileEventSource.LowLevelApi.Lines
{
	public class TitleLine : Line
	{
		public string Title { get; }

		public TitleLine(string title)
		{
			Title = title;
		}
	}
}