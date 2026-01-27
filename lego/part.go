package lego

type PartName string
type PartId string

type Part struct {
	Id   PartId
	Name PartName

	Sources []Source
}

type Source struct {
	SourceName string
	PartId     string
}

func NewPart(partId PartId, name PartName) *Part {
	return &Part{Id: partId, Name: name}
}

type InventoryPart struct {
	Part

	ColourId ColorId
	Quantity int
}
