package brickowl

import (
	"brickrecon/lego"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

type Owlette interface {
	getInventory(boid string) ([]inventoryItem, error)
	lookupSetBoid(setNumber lego.SetId) (string, error)
	lookupParts(boids []lego.BrickOwlPart) (map[lego.BrickOwlPart]lookupItem, error)
	lookup(boid string) (*lookupItem, error)
	listColours() (map[flexInt]colourItem, error)
}

type lowLevelApi struct {
	apiKey string
}

func newLowLevelApi(key string) Owlette {
	return &lowLevelApi{apiKey: key}
}

func (api *lowLevelApi) getInventory(boid string) ([]inventoryItem, error) {

	args := map[string]string{
		"boid": boid,
	}

	var dto inventoryResponse

	if err := api.makeRequest("https://api.brickowl.com/v1/catalog/inventory", args, &dto); err != nil {
		return nil, err
	}

	return dto.Inventory, nil
}

func (api *lowLevelApi) lookupSetBoid(setNumber lego.SetId) (string, error) {

	args := map[string]string{
		"type":    "Set",
		"id_type": "set_number",
		"id":      string(setNumber),
	}

	var dto idlookupResponse

	if err := api.makeRequest("https://api.brickowl.com/v1/catalog/id_lookup", args, &dto); err != nil {
		return "", err
	}

	if len(dto.Boids) == 0 {
		return "", errors.New("No Boids found")
	}

	return dto.Boids[0], nil
}

func (api *lowLevelApi) lookupParts(boids []lego.BrickOwlPart) (map[lego.BrickOwlPart]lookupItem, error) {
	if len(boids) > 100 {
		return nil, errors.New("Max 100 ids")
	}

	args := map[string]string{
		"boids": boidCsv(boids),
	}

	var dto bulkLookupResponse

	if err := api.makeRequest("https://api.brickowl.com/v1/catalog/bulk_lookup", args, &dto); err != nil {
		return nil, err
	}

	return dto.Items, nil
}

func (api *lowLevelApi) lookup(boid string) (*lookupItem, error) {
	args := map[string]string{
		"boid": boid,
	}

	var dto lookupItem

	if err := api.makeRequest("https://api.brickowl.com/v1/catalog/lookup", args, &dto); err != nil {
		return nil, err
	}

	return &dto, nil
}

func (api *lowLevelApi) listColours() (map[flexInt]colourItem, error) {
	args := map[string]string{}

	var dto map[flexInt]colourItem

	if err := api.makeRequest("https://api.brickowl.com/v1/catalog/color_list", args, &dto); err != nil {
		return nil, err
	}

	return dto, nil
}

func (bo *lowLevelApi) makeRequest(url string, args map[string]string, dto interface{}) error {
	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		return err
	}

	q := req.URL.Query()
	q.Add("key", bo.apiKey)
	for name, value := range args {
		q.Add(name, value)
	}
	req.URL.RawQuery = q.Encode()

	client := http.Client{}
	res, err := client.Do(req)

	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode < 200 || res.StatusCode > 299 {
		return fmt.Errorf("Unexpected statusCode: %v", res.StatusCode)
	}

	content, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}

	return json.Unmarshal(content, &dto)
}

func boidCsv(boids []lego.BrickOwlPart) string {
	if len(boids) == 0 {
		return ""
	}
	var (
		sep = []byte(",")
		// preallocate for len(sep) + assume at least 1 character
		out = make([]byte, 0, (1+len(sep))*len(boids))
	)
	for _, s := range boids {
		out = append(out, s...)
		out = append(out, sep...)
	}
	return string(out[:len(out)-len(sep)])
}

type bulkLookupResponse struct {
	Items map[lego.BrickOwlPart]lookupItem
}

type lookupItem struct {
	Boid     string
	Name     string
	Type     string
	ColourID flexInt `json:"color_id"`
	IDs      idMap
}
type idlookupResponse struct {
	Boids []string
}

type inventoryResponse struct {
	Inventory []inventoryItem
}

type inventoryItem struct {
	Boid     lego.BrickOwlPart
	Quantity flexInt
}

type colourItem struct {
	ID   string
	Name lego.ColourName
	Hex  lego.HexColour

	LDrawIDs     []flexInt `json:"ldraw_ids"`
	BrickLinkIDs []flexInt `json:"bl_ids"`
}
