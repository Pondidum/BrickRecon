using System;
using System.Collections.Generic;
using System.Linq;
using System.Threading.Tasks;
using Amazon.Lambda.Model;
using BsxProcessor.Domain;
using Newtonsoft.Json;
using Shouldly;
using Xunit;

namespace BsxProcessor.Tests
{
	public class ImageCacheDispatcherTests
	{
		private static readonly Random random = new Random();

		private readonly List<string> _payloads;
		private readonly List<RequestEvent> _requests;
		private readonly ImageCacheDispatcher _cache;

		public ImageCacheDispatcherTests()
		{
			_payloads = new List<string>();
			_requests = new List<RequestEvent>();
			_cache = new ImageCacheDispatcher(new Config(), request =>
			{
				_payloads.Add(request.Payload);
				_requests.Add(JsonConvert.DeserializeObject<RequestEvent>(request.Payload));
				return Task.FromResult(new InvokeResponse());
			});
		}

		private static Part CreatePart()
		{
			var colors = Enum.GetValues(typeof(Colors)).Cast<Colors>().ToArray();

			return new Part
			{
				PartNumber = random.Next(100, 1000).ToString(),
				Color = colors[random.Next(0, colors.Length)]
			};
		}

		[Fact]
		public async Task When_there_are_no_items_to_dispatch()
		{
			await _cache.Dispatch();

			_requests.ShouldBeEmpty();
		}

		[Fact]
		public async Task When_there_is_one_item_to_dispatch()
		{
			var part = CreatePart();
			_cache.Add(new[] { part });

			await _cache.Dispatch();

			_requests.Single().Parts.ShouldHaveSingleItem().PartNo.ShouldBe(part.PartNumber);
			_requests.Single().Parts.ShouldHaveSingleItem().Color.ShouldBe((int)part.Color);
		}

		[Fact]
		public async Task When_the_same_item_is_added_multiple_times()
		{
			var part = CreatePart();
			_cache.Add(new[] { part, part, part });

			await _cache.Dispatch();

			_requests.Single().Parts.ShouldHaveSingleItem().PartNo.ShouldBe(part.PartNumber);
			_requests.Single().Parts.ShouldHaveSingleItem().Color.ShouldBe((int)part.Color);
		}

		[Fact]
		public async Task When_there_are_multiple_items_to_dispatch()
		{
			_cache.Add(new[] { CreatePart(), CreatePart(), CreatePart() });

			await _cache.Dispatch();

			_requests.Single().Parts.Count().ShouldBe(3);
		}

		[Fact]
		public async Task When_there_are_more_items_than_batch_size()
		{
			_cache.Add(Enumerable.Range(0, ImageCacheDispatcher.BatchSize + 5).Select(i => CreatePart()));

			await _cache.Dispatch();

			_requests.Count.ShouldBe(2);
		}

		[Fact]
		public async Task When_serializing_the_correct_format_is_used()
		{
			var parts = new[] { CreatePart(), CreatePart() };

			_cache.Add(parts);
			await _cache.Dispatch();

			var expected = string.Join(",", parts.Select(p => $"{{\"partno\":\"{p.PartNumber}\",\"color\":{(int)p.Color}}}"));

			_payloads.Single().ShouldBe($"{{\"parts\":[{expected}]}}");
		}

		private class RequestEvent
		{
			public IEnumerable<RequestPart> Parts { get; set; }
		}

		private class RequestPart
		{
			public string PartNo { get; set; }
			public int Color { get; set; }
		}
	}
}
