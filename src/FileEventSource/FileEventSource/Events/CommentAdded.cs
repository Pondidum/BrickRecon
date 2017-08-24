namespace FileEventSource.Events
{
	public class CommentAdded
	{
		public string Comment { get; }

		public CommentAdded(string comment)
		{
			Comment = comment;
		}
	}
}