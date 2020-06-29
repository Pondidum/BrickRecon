package brickowl

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

func (bo *BrickOwlApi) GetSetName(setNumber string) (string, error) {

	setBoid, err := bo.getSetBoid(setNumber)
	if err != nil {
		return "", err
	}

	info, err := bo.lookup(setBoid)
	if err != nil {
		return "", err
	}

	return info.Name, nil
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

			if itemData.Type == "Part" {
				part := createPart(colours, item, itemData)
				parts = append(parts, part)
			}
		}
	}

	return parts, nil
}

func (bo *BrickOwlApi) makeRequest(url string, args map[string]string, dto interface{}) error {
	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		return err
	}

	q := req.URL.Query()
	q.Add("key", bo.key)
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

func (bo *BrickOwlApi) getSetBoid(setNumber string) (string, error) {

	args := map[string]string{
		"type":    "Set",
		"id_type": "set_number",
		"id":      setNumber,
	}

	var dto idlookupResponse

	if err := bo.makeRequest("https://api.brickowl.com/v1/catalog/id_lookup", args, &dto); err != nil {
		return "", err
	}

	if len(dto.Boids) == 0 {
		return "", errors.New("No Boids found")
	}

	return dto.Boids[0], nil
}

func (bo *BrickOwlApi) getInventory(boid string) ([]inventoryItem, error) {

	args := map[string]string{
		"boid": boid,
	}

	var dto inventoryResponse

	if err := bo.makeRequest("https://api.brickowl.com/v1/catalog/inventory", args, &dto); err != nil {
		return nil, err
	}

	return dto.Inventory, nil
}

func (bo *BrickOwlApi) lookup(boid string) (*lookupItem, error) {
	args := map[string]string{
		"boid": boid,
	}

	var dto lookupItem

	if err := bo.makeRequest("https://api.brickowl.com/v1/catalog/lookup", args, &dto); err != nil {
		return nil, err
	}

	return &dto, nil
}

func (bo *BrickOwlApi) lookupParts(boids []string) (map[string]lookupItem, error) {
	if len(boids) > 100 {
		return nil, errors.New("Max 100 ids")
	}

	args := map[string]string{
		"boids": strings.Join(boids, ","),
	}

	var dto bulkLookupResponse

	if err := bo.makeRequest("https://api.brickowl.com/v1/catalog/bulk_lookup", args, &dto); err != nil {
		return nil, err
	}

	return dto.Items, nil
}

func (bo *BrickOwlApi) loadColours() (map[flexInt]colourItem, error) {
	args := map[string]string{}

	var dto map[flexInt]colourItem

	if err := bo.makeRequest("https://api.brickowl.com/v1/catalog/color_list", args, &dto); err != nil {
		return nil, err
	}

	return dto, nil
}

func createPart(colours map[flexInt]colourItem, item inventoryItem, additional lookupItem) lego.Part {
	ldrawID, found := additional.IDs["ldraw"]

	if !found {
		ldrawID, found = additional.IDs["design_id"]
	}

	colour := partColour(colours, additional.ColourID)

	name := sanitiseName(additional.Name, ldrawID, colour)

	return lego.Part{
		ID:       lego.PartID(ldrawID),
		Name:     lego.PartName(name),
		Quantity: int(item.Quantity),
		Colour:   colour,
		Aliases: lego.PartAliases{
			LDrawID: ldrawID,
			Boid:    item.Boid,
		},
	}
}

func sanitiseName(name string, id string, colour lego.Colour) string {

	name = strings.TrimPrefix(name, "LEGO ")
	name = strings.TrimPrefix(name, colour.Name)
	name = name[0:strings.LastIndex(name, "(")]
	name = strings.TrimSpace(name)

	return name
}

func partColour(colours map[flexInt]colourItem, colourID flexInt) lego.Colour {
	colourInfo := colours[colourID]

	colourAliases := lego.ColourAliases{
		BrickLinkID: lego.BrickLinkColour(colourInfo.BrickLinkIDs[0]),
		LDrawID:     lego.LDrawColour(colourInfo.LDrawIDs[0]),
		Boid:        lego.BrickOwlColour(colourID),
	}

	return lego.Colour{
		ID:      colourAliases.BrickLinkID,
		Aliases: colourAliases,
		Name:    colourInfo.Name,
	}
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

type bulkLookupResponse struct {
	Items map[string]lookupItem
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
	Boid     string
	Quantity flexInt
}

type colourItem struct {
	ID   string
	Name string

	LDrawIDs     []flexInt `json:"ldraw_ids"`
	BrickLinkIDs []flexInt `json:"bl_ids"`
}

type flexInt int

func (fi *flexInt) UnmarshalJSON(b []byte) error {
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
	*fi = flexInt(i)
	return nil
}

type idMap map[string]string

func (idmap *idMap) UnmarshalJSON(b []byte) error {

	var ids []struct {
		ID   string
		Type string
	}

	if err := json.Unmarshal(b, &ids); err != nil {
		return err
	}

	m := idMap{}

	for _, pair := range ids {
		m[pair.Type] = pair.ID
	}

	*idmap = m

	return nil
}
