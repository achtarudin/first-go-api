package message

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
)

type MessageSuite struct {
	suite.Suite
}

// TestGoroutine_WG tests multiple goroutines using sync.WaitGroup
// to ensure that all goroutines complete before proceeding.
// It also measures the total time taken to ensure it is within expected limits.
func (suite *MessageSuite) TestGoroutine_WG() {
	var wg sync.WaitGroup
	var mu sync.Mutex
	timeStart := time.Now()
	defer func() {
		timeEnd := time.Now()
		suite.T().Logf("Time taken for TestGoroutine_WG: %v", timeEnd.Sub(timeStart))
		suite.True(timeEnd.Sub(timeStart) <= 4*time.Second, "Time taken should be less than 4 seconds")
	}()

	messages := []Message{
		{name: "Charlie is slow 3 second", slow: true},
		{name: "Charlie is slow 3 second", slow: true},
		{name: "Charlie is slow 3 second", slow: true},
		{name: "Charlie is slow 3 second", slow: true},
		{name: "Charlie is slow 3 second", slow: true},
		{name: "Alice", slow: false},
		{name: "Bob", slow: false},
		{name: "Diana", slow: false},
	}

	results := []string{}
	for _, message := range messages {
		wg.Add(1)

		go func(msg Message) {
			defer wg.Done()
			var result string
			if msg.slow {
				result = msg.SlowGreat()
			} else {
				result = msg.Great()
			}
			mu.Lock()
			results = append(results, result)
			mu.Unlock()
		}(message)
	}

	wg.Wait()
	suite.Equal(len(messages), len(results), "messages and results should have the same length")
}

// TestGoroutine_Channel tests multiple goroutines using channels
// to communicate results and errors back to the main goroutine.
// It also measures the total time taken to ensure it is within expected limits.
func (suite *MessageSuite) TestGoroutine_Channel() {
	timeStart := time.Now()
	defer func() {
		timeEnd := time.Now()
		suite.T().Logf("Time taken for TestGoroutine_Channel: %v", timeEnd.Sub(timeStart))
		suite.True(timeEnd.Sub(timeStart) <= 4*time.Second, "Time taken should be less than 4 seconds")
	}()

	ctx := context.Background()

	messages := []MessageWithError{
		{Message: Message{name: "Message Foo is slow 3 second and no error", slow: true}, err: nil},
		{Message: Message{name: "Message Danu is slow 3 second and error", slow: true}, err: ErrorSlowAndError},
		{Message: Message{name: "Message Alice", slow: false}, err: Error},
		{Message: Message{name: "Message Bob", slow: false}, err: nil},
		{Message: Message{name: "Message Diana", slow: false}, err: nil},
	}

	doneChans := make([]chan string, len(messages))
	errorChans := make([]chan error, len(messages))

	var mu sync.Mutex

	for index, message := range messages {
		doneChans[index] = make(chan string, 1)
		errorChans[index] = make(chan error, 1)

		go func(msg MessageWithError, doneChan chan string, errorChan chan error) {

			if msg.slow {
				msg.SlowGreatCtx(ctx, doneChan, errorChan)
			} else {
				msg.GreatCtx(ctx, doneChan, errorChan)
			}
		}(message, doneChans[index], errorChans[index])

	}

	totalErrors := 0
	totalSuccess := 0
	for i := range messages {
		select {
		case err := <-errorChans[i]:
			if err != nil {
				mu.Lock()
				totalErrors++
				mu.Unlock()

			}
		case msg := <-doneChans[i]:
			mu.Lock()
			totalSuccess++
			mu.Unlock()
			_ = msg
		}
	}

	suite.Equal(len(messages), totalErrors+totalSuccess, "all messages should be processed")
}

func (suite *MessageSuite) TestGoroutine_CtxCancel() {

	timeStart := time.Now()
	defer func() {
		timeEnd := time.Now()
		suite.T().Logf("Time taken for TestGoroutine_CtxCancel: %v", timeEnd.Sub(timeStart))
		suite.True(timeEnd.Sub(timeStart) <= 4*time.Second, "Time taken should be less than 4 seconds")
	}()
	ctx, cancel := context.WithCancel(context.Background())

	var wg sync.WaitGroup

	defer cancel()

	messages := []MessageWithError{
		{Message: Message{name: "Message Foo is slow 3 second and no error", slow: true}, err: nil},
		{Message: Message{name: "Message Bar is slow 3 second and no error", slow: true}, err: nil},

		{Message: Message{name: "Message Danu is slow 3 second and error", slow: true}, err: ErrorSlowAndError},
		{Message: Message{name: "Message Alice", slow: false}, err: Error},
		{Message: Message{name: "Message Bob", slow: false}, err: nil},
		{Message: Message{name: "Message Diana", slow: false}, err: nil},
	}

	var totalErrors int
	var mu sync.Mutex

	for _, message := range messages {
		wg.Add(1)

		go func(msg MessageWithError) {
			defer wg.Done()
			select {
			case <-ctx.Done():
				return
			default:
				if msg.err != nil {
					mu.Lock()
					totalErrors++
					mu.Unlock()
					cancel()
					return
				}

				if msg.slow {
					msg.Message.SlowGreat()
				} else {
					msg.Message.Great()
				}
			}

		}(message)

	}

	wg.Wait()
	suite.NotZero(totalErrors)
	suite.Greater(len(messages), totalErrors)

}

func (suite *MessageSuite) TestGoroutine_CtxCancelWithChannel() {

	timeStart := time.Now()
	defer func() {
		timeEnd := time.Now()
		suite.T().Logf("Time taken for TestGoroutine_CtxCancelWithChannel: %v", timeEnd.Sub(timeStart))
		suite.True(timeEnd.Sub(timeStart) <= 4*time.Second, "Time taken should be less than 4 seconds")
	}()

	var wg sync.WaitGroup
	var totalErrors int
	var totalNotErrors int

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	messages := []MessageWithError{
		{Message: Message{name: "Message Baz", slow: false}, err: Error},
		{Message: Message{name: "Message Bar is slow 3 second and error", slow: true}, err: ErrorSlowAndError},
		{Message: Message{name: "Message Foo is slow 3 second and error", slow: true}, err: ErrorSlowAndError},
		{Message: Message{name: "Message Danu is slow 3 second and no error", slow: true}, err: nil},
		{Message: Message{name: "Message Danu is slow 3 second and no error", slow: true}, err: nil},
		{Message: Message{name: "Message Danu is slow 3 second and no error", slow: true}, err: nil},
		{Message: Message{name: "Message Danu is slow 3 second and no error", slow: true}, err: nil},
		{Message: Message{name: "Message Alice", slow: false}, err: nil},
		{Message: Message{name: "Message Bob", slow: false}, err: nil},
		{Message: Message{name: "Message Diana", slow: false}, err: nil},
	}

	type result struct {
		message string
		err     error
	}
	resultChannel := make(chan result, len(messages))

	for _, message := range messages {
		wg.Add(1)

		go func(msg MessageWithError, done chan string, err chan error) {
			defer wg.Done()

			if msg.slow {
				msg.SlowGreatCtx(ctx, done, err)
			} else {
				msg.GreatCtx(ctx, done, err)
			}

			select {
			case err := <-err:
				if err != nil {
					resultChannel <- result{err: err}
					cancel()
				}
			case msg := <-done:
				if msg != "" {
					resultChannel <- result{message: msg}
				}
			}
		}(message, make(chan string, 1), make(chan error, 1))
	}

	go func() {
		wg.Wait()
		close(resultChannel)
	}()

	for result := range resultChannel {
		if result.err != nil {
			totalErrors++
		}

		if result.message != "" {
			totalNotErrors++
		}
	}

	suite.NotZero(totalErrors, "there is at least 1 error")

}

func TestMessage(t *testing.T) {
	suite.Run(t, &MessageSuite{})
}
