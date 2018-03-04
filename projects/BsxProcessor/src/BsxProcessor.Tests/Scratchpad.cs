using System;
using System.Collections.Generic;
using System.Xml.Linq;
using BsxProcessor.Infrastructure;
using Newtonsoft.Json;
using Shouldly;
using Xunit;
using Xunit.Abstractions;

namespace BsxProcessor.Tests
{
	public class Scratchpad
	{
		private readonly ITestOutputHelper _output;

		public Scratchpad(ITestOutputHelper output)
		{
			_output = output;
		}
		
		[Fact]
		public void When_testing_something()
		{
		}
	}
}
