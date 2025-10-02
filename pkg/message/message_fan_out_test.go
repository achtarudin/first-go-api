package message

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
)

// gow test -count=1 -v -timeout 30s -run ^TestMessageFanOut$ cutbray/first_api/pkg/message

type MessageFanOutAndInSuite struct {
	suite.Suite
}

func (suite *MessageFanOutAndInSuite) TestFanOut() {
	messages := []*Message{
		{name: "Alice (cepat)", slow: false},
		{name: "Bob (lambat)", slow: true},
		{name: "Charlie (cepat)", slow: false},
		{name: "David (lambat)", slow: true},
		{name: "Eve (cepat)", slow: false},
		{name: "Frank (lambat)", slow: true},
	}

	var wg sync.WaitGroup
	numWorkers := 3
	tasks := make(chan *Message, len(messages))

	for i := 1; i <= numWorkers; i++ {
		wg.Add(1)

		go func(index int, wait *sync.WaitGroup, taskChan <-chan *Message) {
			defer wait.Done()
			fmt.Println("Worker", index, "started")

			for task := range taskChan {
				if task.slow {
					task.SlowGreat()
					suite.T().Log("Processing task SlowGreat:", task.name)
				} else {
					task.Great()
					suite.T().Log("Processing task Great:", task.name)
				}
			}
		}(i, &wg, tasks)
	}

	for _, msg := range messages {
		tasks <- msg
	}

	close(tasks)
	wg.Wait()

	suite.True(true)
}

func (suite *MessageFanOutAndInSuite) TestFanIn() {

	timeStart := time.Now()
	defer func() {
		timeEnd := time.Now()
		suite.T().Logf("Time taken for TestFanIn: %v", timeEnd.Sub(timeStart))
		suite.True(timeEnd.Sub(timeStart) <= 4*time.Second, "Time taken should be less than 4 seconds")
	}()

	var wg sync.WaitGroup

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	messages := []MessageWithError{
		{Message: Message{name: "Alice (cepat)", slow: false}, err: nil},
		{Message: Message{name: "Bob (lambat)", slow: true}, err: nil},
		{Message: Message{name: "Charlie (cepat)", slow: false}, err: nil},
		// {Message: Message{name: "David (lambat) (error)", slow: true}, err: fmt.Errorf("ada error")},
		{Message: Message{name: "Eve (cepat)", slow: false}, err: nil},
		{Message: Message{name: "Frank (lambat)", slow: true}, err: nil},
		// {Message: Message{name: "Grace (cepat) (error)", slow: false}, err: fmt.Errorf("ada error")},
		{Message: Message{name: "Hank (lambat)", slow: true}, err: nil},
	}

	numWorkers := 3
	totalTask := len(messages)

	type result struct {
		message string
		err     error
	}

	tasksChannel := make(chan MessageWithError, totalTask)

	resultChannel := make(chan result, totalTask)

	for index := 1; index <= numWorkers; index++ {

		wg.Add(1)

		go func(workerIndex int, taskChan <-chan MessageWithError, resultsChan chan<- result) {
			defer wg.Done()

			fmt.Println("Worker", workerIndex, "started")
			for task := range taskChan {

				doneChan := make(chan string, 1)
				errorChan := make(chan error, 1)

				if task.slow {
					task.SlowGreatCtx(ctx, doneChan, errorChan)
				} else {
					task.GreatCtx(ctx, doneChan, errorChan)
				}

				select {
				case err := <-errorChan:
					if err != nil {
						resultsChan <- result{err: err}
						cancel()
					}
				case msg := <-doneChan:
					if msg != "" {
						fmt.Println(msg)
						resultsChan <- result{message: msg}
					}
				case <-ctx.Done():
					// suite.T().Logf("Worker %d berhenti karena context dibatalkan.", workerIndex)
					return
				}

			}

		}(index, tasksChannel, resultChannel)
	}

	for _, msg := range messages {
		tasksChannel <- msg
	}

	close(tasksChannel)

	go func() {
		wg.Wait()
		close(resultChannel)
	}()

	for range resultChannel {
		// fmt.Println(result.message)
	}

	suite.True(true)
}

func TestMessageFanOut(t *testing.T) {
	suite.Run(t, &MessageFanOutAndInSuite{})
}
