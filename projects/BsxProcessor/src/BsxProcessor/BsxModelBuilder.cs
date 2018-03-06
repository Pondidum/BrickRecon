using System;
using System.IO;
using System.Linq;
using System.Xml.Linq;
using BsxProcessor.Domain;

namespace BsxProcessor
{
	public class BsxModelBuilder : IBsxModelBuilder
	{
		public BsxModel Build(string modelName, XDocument content)
		{
			return new BsxModel
			{
				Name = Path.GetFileNameWithoutExtension(modelName),
				Parts = content.Descendants("Item").Select(PartFromItem)
			};
		}

		private Part PartFromItem(XElement element)
		{
			return new Part
			{
				PartNumber = element.Element("ItemID").Value,
				Name = element.Element("ItemName").Value,
				Category = element.Element("CategoryName").Value,
				Color = (Colors)Enum.Parse(typeof(Colors), element.Element("ColorID").Value),
				Quantity = int.Parse(element.Element("Qty").Value),
			};
		}
	}
}
