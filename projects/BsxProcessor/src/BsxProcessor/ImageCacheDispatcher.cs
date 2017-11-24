using System;
using System.Collections.Generic;
using System.Linq;
using System.Threading.Tasks;
using Amazon.Lambda;
using Amazon.Lambda.Model;
using BsxProcessor.Domain;
using Newtonsoft.Json;

namespace BsxProcessor
{
	public class ImageCacheDispatcher
	{
		public const int BatchSize = 10;

		private readonly Func<InvokeRequest, Task<InvokeResponse>> _dispatch;
		private readonly HashSet<KeyValuePair<string, Colors>> _parts;

		public ImageCacheDispatcher(Func<InvokeRequest, Task<InvokeResponse>> dispatch)
		{
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
				.ToArray();

			foreach (var batch in batches)
			{
				await _dispatch(new InvokeRequest
				{
					InvocationType = InvocationType.Event,
					FunctionName = "brickrecon_imagecache",
					Payload = JsonConvert.SerializeObject(new { parts = batch.AsEnumerable() })
				});
			}
		}
	}
}