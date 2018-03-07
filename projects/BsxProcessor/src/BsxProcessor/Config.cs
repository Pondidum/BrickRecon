using System;

namespace BsxProcessor
{
	public class Config
	{
		public string ImageCacheLambda { get; set; }
		public Uri OutputBucketPath { get; set; }

		public void Validate()
		{
			if (string.IsNullOrWhiteSpace(ImageCacheLambda))
				throw new ArgumentException("The ImageCacheLambda must be set (read from env.IMAGECACHE_LAMBDA)", nameof(ImageCacheLambda));

			if (string.Equals(OutputBucketPath.Scheme, "s3", StringComparison.OrdinalIgnoreCase) == false)
				throw new ArgumentException("The OutputBucketPath must start with 's3://'", nameof(OutputBucketPath));
		}

		public static Config FromEnvironment()
		{
			var config = new Config
			{
				ImageCacheLambda = Environment.GetEnvironmentVariable("IMAGECACHE_LAMBDA"),
				OutputBucketPath = new Uri(Environment.GetEnvironmentVariable("IMAGECACHE_OUTPUTBUCKETPATH"))
			};

			config.Validate();

			return config;
		}
	}
}
