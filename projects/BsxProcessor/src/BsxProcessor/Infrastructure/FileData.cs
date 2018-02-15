namespace BsxProcessor.Infrastructure
{
	public class FileData<TContent>
	{
		public string FullPath { get; set; }
		public string Drive { get; set; }
		public TContent Content { get; set; }
		public bool Exists { get; set; }
	}
}
