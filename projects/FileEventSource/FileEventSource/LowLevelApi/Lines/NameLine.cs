namespace FileEventSource.LowLevelApi.Lines
{
	public class NameLine : Line
	{
		public string Name { get; }

		public NameLine(string name)
		{
			Name = name;
		}
	}
}