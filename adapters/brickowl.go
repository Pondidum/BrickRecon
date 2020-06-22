package adapters

import (
	"brickrecon/lego"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

type BrickOwlApi struct {
	key string
}

func NewBrickOwlApi(key string) *BrickOwlApi {
	return &BrickOwlApi{key}
}

func (bo *BrickOwlApi) getSetBoid(setNumber string) (string, error) {

	req, err := http.NewRequest("GET", "https://api.brickowl.com/v1/catalog/id_lookup", nil)

	if err != nil {
		return "", err
	}

	q := req.URL.Query()
	q.Add("key", bo.key)
	q.Add("type", "Set")
	q.Add("id_type", "set_number")
	q.Add("id", setNumber)
	req.URL.RawQuery = q.Encode()

	client := http.Client{}
	res, err := client.Do(req)

	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	if res.StatusCode < 200 || res.StatusCode > 299 {
		return "", fmt.Errorf("Unexpected statusCode: %v", res.StatusCode)
	}

	content, err := ioutil.ReadAll(res.Body)

	var dto idlookupResponse
	if err := json.Unmarshal(content, &dto); err != nil {
		return "", err
	}

	if len(dto.Boids) == 0 {
		return "", errors.New("No Boids found")
	}

	return dto.Boids[0], err
}

func (bo *BrickOwlApi) GetParts(setNumber string) ([]lego.Part, error) {

	setBoid, err := bo.getSetBoid(setNumber)
	if err != nil {
		return nil, err
	}

	inventory, err := bo.getInventory(setBoid)
	if err != nil {
		return nil, err
	}

	colours, err := bo.loadColours()
	if err != nil {
		return nil, err
	}

	chunks := split(inventory, 100)

	parts := []lego.Part{}

	for _, items := range chunks {

		partBoids := make([]string, len(items))

		for i, item := range items {
			partBoids[i] = item.Boid
		}

		// takes max 100 items
		partData, err := bo.lookupParts(partBoids)
		if err != nil {
			return nil, err
		}

		for _, item := range items {
			itemData := partData[item.Boid]
			part := createPart(colours, item, itemData)

			parts = append(parts, part)
		}
	}

	return parts, nil
}

func createPart(colours map[FlexInt]colourItem, item inventoryItem, additional lookupItem) lego.Part {
	ldrawID := getID(additional.IDs, "ldraw")

	return lego.Part{
		ID:       lego.NewPartID(ldrawID),
		Name:     additional.Name,
		Quantity: ignore(strconv.Atoi(item.Quantity)),
		Colour: lego.Colour{
			ID:   ignore(strconv.Atoi(colours[additional.ColourID].LDrawIDs[0])),
			Name: colours[additional.ColourID].Name,
		},
		Aliases: lego.PartAliases{
			LDrawID: ldrawID,
			Boid:    item.Boid,
		},
	}
}

func ignore(value int, err error) int {
	return value
}

func getID(ids []lookupID, t string) string {
	for _, value := range ids {
		if value.Type == t {
			return value.ID
		}
	}

	return ""
}

func split(buf []inventoryItem, lim int) [][]inventoryItem {
	var chunk []inventoryItem
	chunks := make([][]inventoryItem, 0, len(buf)/lim+1)
	for len(buf) >= lim {
		chunk, buf = buf[:lim], buf[lim:]
		chunks = append(chunks, chunk)
	}
	if len(buf) > 0 {
		chunks = append(chunks, buf[:])
	}
	return chunks
}

func (bo *BrickOwlApi) getInventory(boid string) ([]inventoryItem, error) {

	req, err := http.NewRequest("GET", "https://api.brickowl.com/v1/catalog/inventory", nil)

	if err != nil {
		return nil, err
	}

	q := req.URL.Query()
	q.Add("key", bo.key)
	q.Add("boid", boid)
	req.URL.RawQuery = q.Encode()

	client := http.Client{}
	res, err := client.Do(req)

	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode < 200 || res.StatusCode > 299 {
		return nil, fmt.Errorf("Unexpected statusCode: %v", res.StatusCode)
	}

	content, err := ioutil.ReadAll(res.Body)

	var dto inventoryResponse
	if err := json.Unmarshal(content, &dto); err != nil {
		return nil, err
	}

	return dto.Inventory, nil
}

func (bo *BrickOwlApi) lookupParts(boids []string) (map[string]lookupItem, error) {
	if len(boids) > 100 {
		return nil, errors.New("Max 100 ids")
	}

	req, err := http.NewRequest("GET", "https://api.brickowl.com/v1/catalog/bulk_lookup", nil)

	if err != nil {
		return nil, err
	}

	q := req.URL.Query()
	q.Add("key", bo.key)
	q.Add("boids", strings.Join(boids, ","))
	req.URL.RawQuery = q.Encode()

	ioutil.WriteFile("url", []byte(req.URL.String()), 0666)

	client := http.Client{}
	res, err := client.Do(req)

	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode < 200 || res.StatusCode > 299 {
		return nil, fmt.Errorf("Unexpected statusCode: %v", res.StatusCode)
	}

	content, err := ioutil.ReadAll(res.Body)

	var dto bulkLookupResponse
	if err := json.Unmarshal(content, &dto); err != nil {
		return nil, err
	}

	return dto.Items, nil
}

func (bo *BrickOwlApi) loadColours() (map[FlexInt]colourItem, error) {

	req, err := http.NewRequest("GET", "https://api.brickowl.com/v1/catalog/color_list", nil)
	if err != nil {
		return nil, err
	}

	q := req.URL.Query()
	q.Add("key", bo.key)
	req.URL.RawQuery = q.Encode()

	client := http.Client{}
	res, err := client.Do(req)

	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode < 200 || res.StatusCode > 299 {
		return nil, fmt.Errorf("Unexpected statusCode: %v", res.StatusCode)
	}

	content, err := ioutil.ReadAll(res.Body)

	var dto map[FlexInt]colourItem
	if err := json.Unmarshal(content, &dto); err != nil {
		return nil, err
	}

	return dto, nil
}

type bulkLookupResponse struct {
	Items map[string]lookupItem
}

type lookupItem struct {
	Boid     string
	Name     string
	ColourID FlexInt `json:"color_id"`
	IDs      []lookupID
}

type lookupID struct {
	ID   string
	Type string
}

type idlookupResponse struct {
	Boids []string
}

type inventoryResponse struct {
	Inventory []inventoryItem
}

type inventoryItem struct {
	Boid     string
	Quantity string
}

type colourItem struct {
	ID   string
	Name string

	LDrawIDs     []string `json:"ldraw_ids"`
	BrickLinkIDs []string `json:"bl_ids"`
}

type FlexInt int

func (fi *FlexInt) UnmarshalJSON(b []byte) error {
	if b[0] != '"' {
		return json.Unmarshal(b, (*int)(fi))
	}
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	i, err := strconv.Atoi(s)
	if err != nil {
		return err
	}
	*fi = FlexInt(i)
	return nil
}
