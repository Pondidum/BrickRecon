using System;

namespace BsxProcessor
{
	public class Config
	{
		public string ImageCacheLambda { get; set; }

		public static Config FromEnvironment() => new Config
		{
			ImageCacheLambda = Environment.GetEnvironmentVariable("IMAGECACHE_LAMBDA")
		};
	}
}
