using System;
using System.Collections.Generic;
using System.Linq;
using System.Text.RegularExpressions;
using System.Threading.Tasks;
using Amazon.Lambda;
using Amazon.Lambda.Model;
using BsxProcessor.Domain;
using Newtonsoft.Json;

namespace BsxProcessor
{
	public class ImageCacheDispatcher : IImageCacheDispatcher
	{
		public const int BatchSize = 10;

		private readonly Config _config;
		private readonly Func<InvokeRequest, Task<InvokeResponse>> _dispatch;
		private readonly HashSet<KeyValuePair<string, Colors>> _parts;

		public ImageCacheDispatcher(Config config, Func<InvokeRequest, Task<InvokeResponse>> dispatch)
		{
			_config = config;
			_dispatch = dispatch;
			_parts = new HashSet<KeyValuePair<string, Colors>>();
		}

		public void Add(IEnumerable<Part> parts)
		{
			foreach (var part in parts)
				_parts.Add(new KeyValuePair<string, Colors>(part.PartNumber, part.Color));
		}

		public async Task Dispatch()
		{
			var batches = _parts
				.Select((x, i) => new { Index = i, partno = x.Key, color = x.Value })
				.GroupBy(x => x.Index / BatchSize)
				.Select(group => group.Select(part => new { part.partno, part.color }));

			foreach (var batch in batches)
			{
				await _dispatch(new InvokeRequest
				{
					InvocationType = InvocationType.Event,
					FunctionName = _config.ImageCacheLambda,
					Payload = JsonConvert.SerializeObject(new { parts = batch })
				});
			}
		}
	}
}
