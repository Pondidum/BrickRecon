package lego

type PartName string
type PartId string

type Part struct {
	Id   PartId
	Name PartName

	Category string
}

func NewPart(partId PartId, name PartName) *Part {
	return &Part{Id: partId, Name: name}
}

type InventoryPart struct {
	Part

	ColorId  ColorId
	Quantity int
}
