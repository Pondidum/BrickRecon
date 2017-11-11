namespace FileEventSource.LowLevelApi.Lines
{
	public class CommentLine : Line
	{
		public string Comment { get; }

		public CommentLine(string comment)
		{
			Comment = comment;
		}
	}
}