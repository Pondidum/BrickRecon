using System;
using System.Collections.Generic;

namespace FileEventSource
{
	// http://www.ldraw.org/article/218#linetypes
	public enum LineTypes
	{
		CommentOrMeta,
		SubFileReference,
		Line,
		Triangle,
		Quadrilateral,
		Optional
	}

	public class LineCommands
	{
		public static readonly HashSet<string> Official= new HashSet<string>(StringComparer.OrdinalIgnoreCase)
		{
			"Author",
			"BFC",
			"!CATEGORY",
			"CLEAR",
			"!CMDLINE",
			"!COLOUR",
			"!HELP",
			"!HISTORY",
			"!KEYWORDS",
			"!LDRAW_ORG",
			"LDRAW_ORG",
			"!LICENSE",
			"Name",
			"PAUSE",
			"PRINT",
			"SAVE",
			"STEP",
			"WRITE"
		};
		
		public static readonly HashSet<string> Unofficial= new HashSet<string>(StringComparer.OrdinalIgnoreCase)
		{
			"ROTSTEP",
			"BACKGROUND",
			"BUFEXCHG",
			"GHOST",
			"GROUP",
			"MLCAD",
			"ROTATION"
		};
	}
}
