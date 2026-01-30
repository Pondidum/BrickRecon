package lego

type SetName string
type SetNumber string

type Set struct {
	Number SetNumber
	Name   SetName

	Parts []*InventoryPart
}
