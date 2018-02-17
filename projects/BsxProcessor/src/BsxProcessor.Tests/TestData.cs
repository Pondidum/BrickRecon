namespace BsxProcessor.Tests
{
	public static class TestData
	{
		public const string BsxWithTwoParts =
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

		public const string BsxWithFourParts =
			@"<?xml version=""1.0"" encoding=""UTF-8\""?>
<!DOCTYPE BrickStockXML>
<BrickStockXML>
	<Inventory>
		<Item>
			<ItemID>6558</ItemID>
			<ItemTypeID>P</ItemTypeID>
			<ColorID>11</ColorID>
			<ItemName>Technic, Pin 3L with Friction Ridges Lengthwise</ItemName>
			<ItemTypeName>Part</ItemTypeName>
			<ColorName>Black</ColorName>
			<CategoryID>139</CategoryID>
			<CategoryName>Technic, Pin</CategoryName>
			<Status>I</Status>
			<Qty>0</Qty>
			<Price>0.000</Price>
			<Condition>N</Condition>
		</Item>
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
			<ItemID>11477</ItemID>
			<ItemTypeID>P</ItemTypeID>
			<ColorID>85</ColorID>
			<ItemName>Slope, Curved 2 x 1 No Studs</ItemName>
			<ItemTypeName>Part</ItemTypeName>
			<ColorName>Dark Bluish Gray</ColorName>
			<CategoryID>438</CategoryID>
			<CategoryName>Slope, Curved</CategoryName>
			<Status>I</Status>
			<Qty>0</Qty>
			<Price>0.000</Price>
			<Condition>N</Condition>
		</Item>
		<Item>
			<ItemID>2412b</ItemID>
			<ItemTypeID>P</ItemTypeID>
			<ColorID>85</ColorID>
			<ItemName>Tile, Modified 1 x 2 Grille with Bottom Groove / Lip</ItemName>
			<ItemTypeName>Part</ItemTypeName>
			<ColorName>Dark Bluish Gray</ColorName>
			<CategoryID>38</CategoryID>
			<CategoryName>Tile, Modified</CategoryName>
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
