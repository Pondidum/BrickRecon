namespace FileEventSource.LowLevelApi.Lines
{
	public class CommandLine : Line
	{
		public string Command { get; }
		public string Arguments { get; }

		public CommandLine(string command, string arguments)
		{
			Command = command;
			Arguments = arguments;
		}
	}
}