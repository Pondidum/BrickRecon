package brickowl

import (
	"brickrecon/lego"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
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

func (api *FakeApi) getInventory(boid string) ([]inventoryItem, error) {

	var dto inventoryResponse
	err := readFile(&dto, "./responses/catalog-inventory-%s.json", boid)

	return dto.Inventory, err
}

func (api *FakeApi) lookupSetBoid(setNumber lego.SetId) (string, error) {

	var dto idlookupResponse
	if err := readFile(&dto, "./responses/catalog-idlookup-set-%s.json", setNumber); err != nil {
		return "", err
	}

	return dto.Boids[0], nil
}

func (api *FakeApi) lookupParts(boids []lego.BrickOwlPart) (map[lego.BrickOwlPart]lookupItem, error) {
	return nil, errors.New("lookupParts not implemented")
}

func (api *FakeApi) lookup(boid string) (*lookupItem, error) {
	return nil, errors.New("lookup not implemented")
}

func (api *FakeApi) listColours() (map[flexInt]colourItem, error) {
	var dto map[flexInt]colourItem
	err := readFile(&dto, "./responses/catalog-colourlist.json")

	return dto, err
}

func readFile(dto interface{}, path string, args ...interface{}) error {

	content, err := ioutil.ReadFile(fmt.Sprintf(path, args...))
	if err != nil {
		return err
	}

	if err := json.Unmarshal(content, &dto); err != nil {
		return err
	}

	return nil
}
