package distributor

type MessageMeta struct {
}

func (m *MessageMeta) meta() *MessageMeta { return m }

type Message interface {
	meta() *MessageMeta
}
