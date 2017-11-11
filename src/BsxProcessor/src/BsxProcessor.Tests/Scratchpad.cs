using System;
using System.Collections.Generic;
using Shouldly;
using Xunit;

namespace BsxProcessor.Tests
{
	public class Scratchpad
	{
		[Fact]
		public void When_testing_something()
		{
			var one = new KeyValuePair<string, int>("2033", 11);
			var two = new KeyValuePair<string, int>("2033", 11);
			//var three = new KeyValuePair<string, int>("211", 5);
			
			var hash = new HashSet<KeyValuePair<string, int>>();
			hash.Add(one);
			
			hash.ShouldContain(two);
		}
	}
}
