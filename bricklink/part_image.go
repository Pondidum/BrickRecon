package bricklink

import (
	"brickrecon/lego"
	"context"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/honeycombio/beeline-go"
)

func GetImage(ctx context.Context, part lego.LDrawPart, colour lego.BrickLinkColour) ([]byte, error) {

	url := fmt.Sprintf(`https://img.bricklink.com/ItemImage/PN/%v/%s.png`, colour, part)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		beeline.AddField(ctx, "error", err)
		return nil, err
	}

	defer res.Body.Close()

	beeline.AddField(ctx, "status_code", res.StatusCode)

	if res.StatusCode < 200 || res.StatusCode > 299 {
		return nil, fmt.Errorf("Part not found")
	}

	content, err := ioutil.ReadAll(res.Body)
	if err != nil {
		beeline.AddField(ctx, "error", err)
		return nil, err
	}

	beeline.AddField(ctx, "content_length", len(content))
	return content, nil
}
