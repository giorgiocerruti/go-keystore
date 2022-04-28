package core

const (
	_                     = iota
	EventDelete EventType = iota
	EventPut
)

// The interface that defines the transaction logger
type TransactionLogger interface {
	WriteDelete(key string)
	WritePut(key, value string)
	Err() <-chan error

	Read() (<-chan Event, <-chan error)

	Run()
}

type KeyValueStore struct {
	m        map[string]string
	transact TransactionLogger //This is the port, the constructor accepts the adapter of type TransactionLogger
}

func (store *KeyValueStore) Delete(key string) error {
	delete(store.m, key)
	store.transact.WriteDelete(key)

	return nil
}

func (store *KeyValueStore) Put(key, value string) error {
	store.m[key] = value
	store.transact.WritePut(key, value)
	return nil
}

func (store *KeyValueStore) Restore() error {
	var err error
	events, errors := store.transact.Read()
	e, ok := Event{}, true

	for ok && err == nil {
		select {
		case err, ok = <-errors:
		case e, ok = <-events:
			switch e.EventType {
			case EventDelete:
				err = store.Delete(e.Key)
			case EventPut:
				err = store.Put(e.Key, e.Value)
			}
		}
	}

	store.transact.Run()

	return err
}

type EventType int

//Used to send events throught channels
type Event struct {
	Sequence  uint64
	EventType EventType
	Key       string
	Value     string
}
