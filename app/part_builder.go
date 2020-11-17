package app

import (
	"brickrecon/bricklink"
	"brickrecon/eventstore"
	"brickrecon/lego"
	"brickrecon/lego/projections/allparts"
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"path"
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

func (pb *PartBuilder) StorePart(ctx context.Context, readPart *lego.Part) error {

	var p *lego.PartAggregate

	if pb.knownParts[readPart.Key] {

		if pb.hasImage[readPart.Key] {
			return nil
		}

		p = lego.BlankPart()
		if err := pb.store.LoadAggregate(ctx, eventstore.AggregateID(readPart.Key), p); err != nil {
			return err
		}
	} else {
		p = lego.NewPartFromLDraw(
			readPart.Key,
			readPart.Aliases.LDrawID,
			readPart.Name,
			readPart.Colour.Aliases.LDrawID,
			readPart.Colour.Name,
			readPart.Colour.Category)

		if err := pb.store.SaveAggregate(ctx, p); err != nil {
			return err
		}

		pb.knownParts[p.Key] = true
	}

	if p.HasImage() == false {
		image, err := bricklink.GetImage(ctx, p.Number, readPart.Colour.Aliases.BrickLinkID)
		if err != nil {
			return err
		}

		path, err := pb.storeImage(ctx, "bricklink", p.Number, p.Colour, image)
		if err != nil {
			return err
		}

		p.AttachImage("bricklink", path)

		if err := pb.store.SaveAggregate(ctx, p); err != nil {
			return err
		}
	}

	return nil
}
