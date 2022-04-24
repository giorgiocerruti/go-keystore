package v1

import (
	"errors"
	"fmt"
	"sync"

	"github.com/giorgiocerruti/go-keystore/pkg/logger"
)

var store = struct {
	sync.RWMutex
	m map[string]string
}{m: make(map[string]string)}

var ErrorNoSuchKey = errors.New("no such key")
var log logger.TransactionLogger

func Put(key, value string) error {
	store.Lock()
	store.m[key] = value
	store.Unlock()
	return nil
}

func Get(key string) (string, error) {
	store.RLock()
	v, ok := store.m[key]
	store.RUnlock()

	if !ok {
		return "", ErrorNoSuchKey
	}

	return v, nil
}

func Delete(key string) error {
	store.Lock()
	delete(store.m, key)
	store.Unlock()

	return nil
}

func InitializeTransdactionLogger() (logger.TransactionLogger, error) {
	var err error

	log, err = logger.NewPostgresTransactionLogger(logger.FILENAME)
	if err != nil {
		return nil, fmt.Errorf("failed to create event logger: %w", err)
	}

	events, errors := log.ReadEvents()
	e, ok := logger.Event{}, true

	for ok && err == nil {
		select {
		case err, ok = <-errors:
		case e, ok = <-events:
			switch e.EventType {
			case logger.EventDelete:
				err = Delete(e.Key)
			case logger.EventPut:
				err = Put(e.Key, e.Value)
			}
		}
	}

	log.Run()

	return log, err

}
