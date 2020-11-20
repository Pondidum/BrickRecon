package actions

import (
	"brickrecon/bricklink"
	"brickrecon/brickowl"
	"brickrecon/eventstore"
	"brickrecon/lego"
	"brickrecon/lego/projections/allparts"
	"brickrecon/stud_io"
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
	storeRoot  string
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
		storeRoot:  "./app/static/img/parts",
	}

	return pb, nil
}

func (pb *PartBuilder) FromBrickOwl(ctx context.Context, readPart *brickowl.BrickOwlPart) error {

	return pb.storePart(ctx, readPart.Key, func() *lego.Part {

		p := lego.NewPart(readPart.Key)
		p.AddNames(readPart.Name, readPart.ColourName)
		p.AddBrickOwl(readPart.Boid, readPart.ColourBoid)
		p.AddBrickLink(readPart.BrickLinkID, readPart.BrickLinkColour)

		return p
	})
}

func (pb *PartBuilder) FromWantedList(ctx context.Context, readPart *stud_io.ListPart) error {

	return pb.storePart(ctx, readPart.Key, func() *lego.Part {

		p := lego.NewPart(readPart.Key)
		p.AddNames(readPart.Name, readPart.ColourName)
		p.AddBrickLink(readPart.BrickLinkID, readPart.BrickLinkColour)

		return p
	})

}

func (pb *PartBuilder) storePart(ctx context.Context, key lego.PartKey, createPart func() *lego.Part) error {

	var err error
	ctx, span := beeline.StartSpan(ctx, "store_"+string(key))

	defer func() {
		if err != nil {
			beeline.AddField(ctx, "err", err)
		}
		span.Send()
	}()

	var p *lego.Part

	partExists := pb.knownParts[key]
	beeline.AddField(ctx, "new_part", !partExists)

	if partExists {

		if pb.hasImage[key] {
			beeline.AddField(ctx, "part_has_image", true)
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

	hasImage := p.HasImage()
	beeline.AddField(ctx, "part_has_image", hasImage)

	if hasImage == false {

		path, err := pb.getImage(ctx, p)
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

func (pb *PartBuilder) getImage(ctx context.Context, p *lego.Part) (string, error) {

	storeName := "bricklink"

	imageRelative := path.Join(storeName, fmt.Sprintf("%s-%v.png", p.PartID, p.ColourID))
	fullImagePath := path.Join(pb.storeRoot, imageRelative)

	beeline.AddField(ctx, "image_relative", imageRelative)

	if err := os.MkdirAll(path.Join(pb.storeRoot, storeName), os.ModePerm); err != nil {
		beeline.AddField(ctx, "store_creation_err", err)
		return "", err
	}

	inFilestore := fileExists(fullImagePath)
	beeline.AddField(ctx, "image_in_filestore", inFilestore)

	if inFilestore {
		return imageRelative, nil
	}

	image, err := bricklink.GetImage(ctx, p.BrickLink.PartNumber, p.BrickLink.Colour)
	if err != nil {
		beeline.AddField(ctx, "image_fetch_err", err)
		return "", err
	}

	beeline.AddField(ctx, "image_fetched", true)
	beeline.AddField(ctx, "image_size", len(image))

	if err := ioutil.WriteFile(fullImagePath, image, 0666); err != nil {
		beeline.AddField(ctx, "image_write_err", err)
		return "", err
	}

	beeline.AddField(ctx, "image_written", true)

	return imageRelative, nil
}

func fileExists(filepath string) bool {
	info, err := os.Stat(filepath)
	if os.IsNotExist(err) {
		return false
	}

	return info.IsDir() == false
}
