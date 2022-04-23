package logger

import (
	"bufio"
	"fmt"
	"os"
)

type FileTransactionLogger struct {
	events       chan<- Event //Write-only channell for sending events
	errors       <-chan error //Read-only channel for errors
	lastSequence uint64       //The last used event sequence
	file         *os.File     //file location
}

//Logs the PUT request
func (l *FileTransactionLogger) WritePut(key, value string) {
	l.events <- Event{EventType: EventPut, Key: key, Value: value}
}

//Logs the DELETE request
func (l *FileTransactionLogger) WriteDelete(key string) {
	l.events <- Event{EventType: EventDelete, Key: key}
}

//Return a channel of errors
func (l *FileTransactionLogger) Err() <-chan error {
	return l.errors
}

//creates a gorutine that read events from a channel and
// write them onto a file
func (l *FileTransactionLogger) Run() {
	//event channel
	events := make(chan Event, 16)
	l.events = events

	//error channel
	errors := make(chan error, 1)
	l.errors = errors

	go func() {
		for e := range events {
			//increments last sequence
			l.lastSequence++

			//write the event to the log
			_, err := fmt.Fprintf(
				l.file,
				"%d\t%d\t%s\t%s\n",
				l.lastSequence, e.EventType, e.Key, e.Value,
			)

			if err != nil {
				errors <- err
				return
			}
		}
	}()

}

func (l *FileTransactionLogger) ReadEvents() (<-chan Event, <-chan error) {
	scanner := bufio.NewScanner(l.file)
	outEvent := make(chan Event)
	outError := make(chan error)

	go func() {
		var e Event

		defer close(outEvent)
		defer close(outError)

		for scanner.Scan() {
			line := scanner.Text()

			if _, err := fmt.Sscanf(line, "%d\t%d\t%s\t%s",
				&e.Sequence, &e.EventType, &e.Key, &e.Value); err != nil {
				outError <- fmt.Errorf("input parse error %w", err)
			}

			//check the sequernce integrity
			if l.lastSequence >= e.Sequence {
				outError <- fmt.Errorf("transaction numbers out of sequence")
				return
			}

			l.lastSequence = e.Sequence
			outEvent <- e

		}

	}()

	return outEvent, outError
}

//Constructor
func NewFileTransactionLogger(filename string) (TransactionLogger, error) {
	file, err := os.OpenFile(filename, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0755)
	if err != nil {
		return nil, fmt.Errorf("cannot otpen transactionasl log file %w", err)
	}

	return &FileTransactionLogger{file: file}, nil
}
