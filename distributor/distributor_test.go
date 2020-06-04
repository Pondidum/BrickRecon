package distributor

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRegistration(t *testing.T) {

	d := NewDistributor()

	d.RegisterFor(&TestMessage{}, func(m Message) {
		message := m.(*TestMessage)
		assert.Equal(t, "some value", message.Value)
	})

	assert.Len(t, d.topics, 1)
	assert.Contains(t, d.topics, "TestMessage")
	assert.Len(t, d.topics["TestMessage"], 1)

	d.topics["TestMessage"][0](&TestMessage{Value: "some value"})
}

func TestDispatch(t *testing.T) {
	d := NewDistributor()

	messages := make(chan *TestMessage)

	d.RegisterFor(&TestMessage{}, func(m Message) {
		messages <- m.(*TestMessage)
	})

	d.Dispatch(&TestMessage{Value: "Test"})

	handled := <-messages
	assert.Equal(t, "Test", handled.Value, "message handler was not called")
}

func TestMultipleListeners(t *testing.T) {
	d := NewDistributor()

	a := make(chan string)
	b := make(chan string)

	d.RegisterFor(&TestMessage{}, func(m Message) {
		a <- "one"
	})

	d.RegisterFor(&TestMessage{}, func(m Message) {
		b <- "two"
	})

	d.Dispatch(&TestMessage{Value: "Test"})

	assert.Equal(t, []string{"one", "two"}, []string{<-a, <-b})

}

type TestMessage struct {
	MessageMeta

	Value string
}
