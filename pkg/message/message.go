package message

import (
	"context"
	"errors"
	"fmt"
	"time"
)

var (
	ErrorSlowAndError = errors.New("You are slow and error")
	Error             = errors.New("You get error")
)

type Message struct {
	name string
	slow bool
}

func (m *Message) Great() string {
	message := fmt.Sprintf("Hello, %s !", m.name)
	return message
}

func (m *Message) SlowGreat() string {
	time.Sleep(3 * time.Second)
	return m.Great()
}

type MessageWithError struct {
	Message
	err error
}

func (m *MessageWithError) GreatCtx(ctx context.Context, doneChan chan string, errorChan chan error) {
	if m.err != nil {
		errorChan <- fmt.Errorf("%s: %v", m.name, m.err)
		return
	}
	select {
	case <-ctx.Done():
		errorChan <- fmt.Errorf("%s: %v", m.name, ctx.Err())
		return
	default:
		doneChan <- m.Message.Great()
	}
}

func (m *MessageWithError) SlowGreatCtx(ctx context.Context, doneChan chan string, errorChan chan error) {
	if m.err != nil {
		errorChan <- fmt.Errorf("%s: %v", m.name, m.err)
		return
	}
	select {
	case <-ctx.Done():
		errorChan <- fmt.Errorf("%s: %v", m.name, ctx.Err())
		return
	case <-time.After(3 * time.Second):
		select {
		case <-ctx.Done():
			errorChan <- fmt.Errorf("%s: %v", m.name, ctx.Err())
			return
		default:
			doneChan <- m.Message.Great()
		}
	}

}
