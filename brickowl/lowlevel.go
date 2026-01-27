package brickowl

import (
	"brickrecon/lego"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

type Boid string

type Owlette interface {
	getInventory(boid Boid) ([]inventoryItem, error)
	lookupSetBoid(setNumber lego.SetId) (Boid, error)
	lookupParts(boids []Boid) (map[Boid]lookupItem, error)
	lookup(boid Boid) (*lookupItem, error)
}

type lowLevelApi struct {
	apiKey string
}

func newLowLevelApi(key string) Owlette {
	return &lowLevelApi{apiKey: key}
}

func (api *lowLevelApi) getInventory(boid Boid) ([]inventoryItem, error) {

	args := map[string]string{
		"boid": string(boid),
	}

	var dto inventoryResponse

	if err := api.makeRequest("https://api.brickowl.com/v1/catalog/inventory", args, &dto); err != nil {
		return nil, err
	}

	return dto.Inventory, nil
}

func (api *lowLevelApi) lookupSetBoid(setNumber lego.SetId) (Boid, error) {

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

func (api *lowLevelApi) lookupParts(boids []Boid) (map[Boid]lookupItem, error) {
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

func (api *lowLevelApi) lookup(boid Boid) (*lookupItem, error) {
	args := map[string]string{
		"boid": string(boid),
	}

	var dto lookupItem

	if err := api.makeRequest("https://api.brickowl.com/v1/catalog/lookup", args, &dto); err != nil {
		return nil, err
	}

	return &dto, nil
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

	content, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}

	return json.Unmarshal(content, &dto)
}

func boidCsv(boids []Boid) string {
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
	Items map[Boid]lookupItem
}

type lookupItem struct {
	Boid      Boid
	Name      string
	Type      string
	ColorID   flexInt `json:"color_id"`
	ColorName string  `json:"color_name"`
	IDs       idMap
}
type idlookupResponse struct {
	Boids []Boid
}

type inventoryResponse struct {
	Inventory []inventoryItem
}

type inventoryItem struct {
	Boid     Boid
	Quantity flexInt
}
