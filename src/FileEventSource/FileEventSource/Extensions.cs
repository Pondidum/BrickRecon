using System;

namespace FileEventSource
{
	public static class Extensions
	{
		public static bool EqualsIgnore(this string first, string second)
		{
			return string.Equals(first, second, StringComparison.OrdinalIgnoreCase);
		}
	}
}
