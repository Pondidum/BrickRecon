package app

import (
	"brickrecon/bricklink"
	"brickrecon/brickowl"
	"brickrecon/eventstore"
	"brickrecon/lego"
	"brickrecon/lego/projections/allparts"
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"path"

	"github.com/honeycombio/beeline-go"
)

type PartBuilder struct {
	store      eventstore.EventStore
	knownParts map[lego.PartKey]bool
	hasImage   map[lego.PartKey]bool
}

func NewPartBuilder(ctx context.Context, store eventstore.EventStore) (*PartBuilder, error) {

	view := allparts.NewAllPartsView()
	if err := store.ReadView(ctx, allparts.ProjectionName, &view); err != nil {
		return nil, err
	}

	pb := &PartBuilder{
		store:      store,
		knownParts: view.KnownParts,
		hasImage:   view.HasImage,
	}

	return pb, nil
}

func (pb *PartBuilder) storeImage(ctx context.Context, sourceName string, number lego.LDrawPart, colour lego.LDrawColour, image []byte) (string, error) {
	location := "./app/static/img/parts"
	filename := fmt.Sprintf("%s-%v.png", number, colour)

	directory := path.Join(location, sourceName)
	file := path.Join(directory, filename)

	if err := os.MkdirAll(directory, os.ModePerm); err != nil {
		return "", err
	}

	if err := ioutil.WriteFile(file, image, 0666); err != nil {
		return "", err
	}

	return path.Join(sourceName, filename), nil
}

func (pb *PartBuilder) FromBrickOwl(ctx context.Context, readPart *brickowl.BrickOwlPart) error {

	return pb.storePart(ctx, readPart.Key, func() *lego.PartAggregate {

		p := lego.NewPart(readPart.Key)
		p.AddNames(readPart.Name, readPart.ColourName)
		p.AddBrickOwl(readPart.Boid, readPart.ColourBoid)
		p.AddBrickLink(readPart.BrickLinkID, readPart.BrickLinkColour)

		return p
	})
}

func (pb *PartBuilder) FromWantedList(ctx context.Context, readPart *lego.Part) error {

	return pb.storePart(ctx, readPart.Key, func() *lego.PartAggregate {

		p := lego.NewPart(readPart.Key)
		p.AddNames(readPart.Name, readPart.Colour.Name)
		p.AddBrickLink(readPart.Aliases.BrickLinkID, readPart.Colour.Aliases.BrickLinkID)

		return p
	})

}

func (pb *PartBuilder) storePart(ctx context.Context, key lego.PartKey, createPart func() *lego.PartAggregate) error {
	var err error
	defer func() {
		if err != nil {
			beeline.AddField(ctx, string(key)+"_error", err)
		}
	}()

	var p *lego.PartAggregate

	if pb.knownParts[key] {

		if pb.hasImage[key] {
			return nil
		}

		p = lego.BlankPart()
		if err = pb.store.LoadAggregate(ctx, eventstore.AggregateID(key), p); err != nil {
			return err
		}
	} else {

		p = createPart()

		if err = pb.store.SaveAggregate(ctx, p); err != nil {
			return err
		}

		pb.knownParts[p.Key] = true
	}

	if p.HasImage() == false {
		image, err := bricklink.GetImage(ctx, p.BrickLink.PartNumber, p.BrickLink.Colour)
		if err != nil {
			return err
		}

		path, err := pb.storeImage(ctx, "bricklink", p.PartID, p.ColourID, image)
		if err != nil {
			return err
		}

		p.AttachImage("bricklink", path)

		if err = pb.store.SaveAggregate(ctx, p); err != nil {
			return err
		}
	}

	return nil
}
