using System;

namespace BsxProcessor
{
	public class Config
	{
		public string ImageCacheLambda { get; set; }
		public Uri OutputBucketPath { get; set; }

		public static Config FromEnvironment() => new Config
		{
			ImageCacheLambda = Environment.GetEnvironmentVariable("IMAGECACHE_LAMBDA"),
			OutputBucketPath = new Uri(Environment.GetEnvironmentVariable("IMAGECACHE_OUTPUTBUCKETPATH"))
		};
	}
}
