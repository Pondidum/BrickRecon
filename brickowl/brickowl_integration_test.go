package brickowl

import (
	"brickrecon/lego"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetInventory(t *testing.T) {
	t.Parallel()

	owl := NewBrickOwlApi(os.Getenv("BRICKOWL_API_KEY"))
	// owl := BrickOwlApi{api: &FakeApi{}}

	parts, err := owl.GetParts("75193-1")

	assert.NoError(t, err)
	assert.Len(t, parts, 47)
}

type FakeApi struct{}

var _ Owlette = &FakeApi{}

func (api *FakeApi) getInventory(boid Boid) ([]inventoryItem, error) {

	var dto inventoryResponse
	err := readFile(&dto, "./responses/catalog-inventory-%s.json", boid)

	return dto.Inventory, err
}

func (api *FakeApi) lookupSetBoid(setNumber lego.SetNumber) (Boid, error) {

	var dto idlookupResponse
	if err := readFile(&dto, "./responses/catalog-idlookup-set-%s.json", setNumber); err != nil {
		return "", err
	}

	return dto.Boids[0], nil
}

func (api *FakeApi) lookupParts(boids []Boid) (map[Boid]lookupItem, error) {
	return nil, errors.New("lookupParts not implemented")
}

func (api *FakeApi) lookup(boid Boid) (*lookupItem, error) {
	return nil, errors.New("lookup not implemented")
}

func readFile(dto interface{}, path string, args ...interface{}) error {

	content, err := os.ReadFile(fmt.Sprintf(path, args...))
	if err != nil {
		return err
	}

	if err := json.Unmarshal(content, &dto); err != nil {
		return err
	}

	return nil
}
