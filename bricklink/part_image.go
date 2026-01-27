package bricklink

import (
	"brickrecon/lego"
	"context"
	"fmt"
	"io"
	"net/http"
)

func GetImage(ctx context.Context, part lego.InventoryPart) ([]byte, error) {

	url := fmt.Sprintf(`https://img.bricklink.com/ItemImage/PN/%v/%s.png`, part.ColourId, part.Id)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	if res.StatusCode < 200 || res.StatusCode > 299 {
		return nil, fmt.Errorf("Part not found")
	}

	content, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	return content, nil
}
