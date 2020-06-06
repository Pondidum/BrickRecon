package background

import (
	"mvc/distributor"
	"mvc/lego"
)

type PartsAddedMessage struct {
	distributor.Message

	Parts []lego.Part
}
