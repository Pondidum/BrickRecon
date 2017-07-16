using System.Linq;
using System.Xml.Linq;
using Shouldly;
using Xunit;

namespace BsxProcessor.Tests
{
	public class BsxModelBuilderTests
	{
		[Fact]
		public void When_mapping_a_file()
		{
			var document = XDocument.Parse(Xml);
			var model = new BsxModelBuilder().Build(document);

			var part = model.Parts.First();
			model.Parts.ShouldSatisfyAllConditions(
				() => model.Parts.Count().ShouldBe(2),
				() => part.PartNumber.ShouldBe(3039),
				() => part.Name.ShouldBe("Slope 45 2 x 2"),
				() => part.Category.ShouldBe("Slope"),
				() => part.Color.ShouldBe(Colors.Black),
				() => part.Quantity.ShouldBe(5)
			);
		}

		private const string Xml =
@"<?xml version=""1.0"" encoding=""UTF-8\""?>
<!DOCTYPE BrickStockXML>
<BrickStockXML>
	<Inventory>
		<Item>
			<ItemID>3039</ItemID>
			<ItemTypeID>P</ItemTypeID>
			<ColorID>11</ColorID>
			<ItemName>Slope 45 2 x 2</ItemName>
			<ItemTypeName>Part</ItemTypeName>
			<ColorName>Black</ColorName>
			<CategoryID>31</CategoryID>
			<CategoryName>Slope</CategoryName>
			<Status>I</Status>
			<Qty>5</Qty>
			<Price>0.000</Price>
			<Condition>N</Condition>
		</Item>
		<Item>
			<ItemID>32039</ItemID>
			<ItemTypeID>P</ItemTypeID>
			<ColorID>11</ColorID>
			<ItemName>Technic, Axle Connector with Axle Hole</ItemName>
			<ItemTypeName>Part</ItemTypeName>
			<ColorName>Black</ColorName>
			<CategoryID>133</CategoryID>
			<CategoryName>Technic, Connector</CategoryName>
			<Status>I</Status>
			<Qty>0</Qty>
			<Price>0.000</Price>
			<Condition>N</Condition>
		</Item>
	</Inventory>
</BrickStockXML>
";
	}
}
