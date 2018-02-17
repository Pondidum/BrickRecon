using System.Linq;
using System.Xml.Linq;
using BsxProcessor.Domain;
using BsxProcessor.Infrastructure;
using Shouldly;
using Xunit;

namespace BsxProcessor.Tests
{
	public class BsxModelBuilderTests
	{
		[Fact]
		public void When_mapping_a_file()
		{
			var model = new BsxModelBuilder().Build(new FileData<XDocument>
			{
				Drive = "s3",
				FullPath = "some/path/to/a/model.bsx",
				Content = XDocument.Parse(TestData.BsxWithTwoParts)
			});

			var part = model.Parts.First();
			model.ShouldSatisfyAllConditions(
				() => model.Name.ShouldBe("model"),
				() => model.Parts.Count().ShouldBe(2),
				() => part.PartNumber.ShouldBe("3039"),
				() => part.Name.ShouldBe("Slope 45 2 x 2"),
				() => part.Category.ShouldBe("Slope"),
				() => part.Color.ShouldBe(Colors.Black),
				() => part.Quantity.ShouldBe(5)
			);
		}
	}
}
