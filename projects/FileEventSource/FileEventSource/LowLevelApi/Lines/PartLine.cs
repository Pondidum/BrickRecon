namespace FileEventSource.LowLevelApi.Lines
{
	public class PartLine : Line
	{
		public string Part { get; }

		public PartLine(string part)
		{
			Part = part;
		}
	}
}