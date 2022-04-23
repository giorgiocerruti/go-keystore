package logger

//Set an ID for each event
const (
	_                     = iota
	EventDelete EventType = iota
	EventPut
)

type TransactionLogger interface {
	WriteDelete(key string)
	WritePut(key, value string)
	Err() <-chan error

	ReadEvents() (<-chan Event, <-chan error)

	Run()
}

type EventType int

//Used to send events throught channels
type Event struct {
	Sequence  uint64
	EventType EventType
	Key       string
	Value     string
}
