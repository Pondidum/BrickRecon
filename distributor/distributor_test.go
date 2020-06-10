package distributor

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRegistration(t *testing.T) {

	d := NewDistributor()

	d.RegisterFor(&TestMessage{}, func(ctx context.Context, m Message) {
		message := m.(*TestMessage)
		assert.Equal(t, "some value", message.Value)
	})

	assert.Len(t, d.topics, 1)
	assert.Contains(t, d.topics, "TestMessage")
	assert.Len(t, d.topics["TestMessage"], 1)

	d.topics["TestMessage"][0].action(context.Background(), &TestMessage{Value: "some value"})
	assert.NotEmpty(t, d.topics["TestMessage"][0].name)
}

func TestDispatch(t *testing.T) {
	d := NewDistributor()

	messages := make(chan *TestMessage)

	d.RegisterFor(&TestMessage{}, func(ctx context.Context, m Message) {
		messages <- m.(*TestMessage)
	})

	d.Dispatch(context.Background(), &TestMessage{Value: "Test"})

	handled := <-messages
	assert.Equal(t, "Test", handled.Value, "message handler was not called")
}

func TestMultipleListeners(t *testing.T) {
	d := NewDistributor()

	a := make(chan string)
	b := make(chan string)

	d.RegisterFor(&TestMessage{}, func(ctx context.Context, m Message) {
		a <- "one"
	})

	d.RegisterFor(&TestMessage{}, func(ctx context.Context, m Message) {
		b <- "two"
	})

	d.Dispatch(context.Background(), &TestMessage{Value: "Test"})

	assert.Equal(t, []string{"one", "two"}, []string{<-a, <-b})

}

func TestNoValidListeners(t *testing.T) {
	d := NewDistributor()

	called := false

	d.RegisterFor(&TestMessage{}, func(ctx context.Context, m Message) {
		called = true
	})

	wait := d.Dispatch(context.Background(), &OtherTestMessage{Value: "Test"})
	wait()

	assert.False(t, called, "the handler should not have been invoked")

}

type TestMessage struct {
	MessageMeta

	Value string
}

type OtherTestMessage struct {
	MessageMeta

	Value string
}

func TestSnakeCase(t *testing.T) {
	assert.Equal(t, "test_message", slither("TestMessage"))
	assert.Equal(t, "test_message", slither("testMessage"))
	assert.Equal(t, "is_html", slither("IsHTML"))
	assert.Equal(t, "is_html_value", slither("IsHTMLValue"))
}
