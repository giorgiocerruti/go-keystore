package core

//The constructor accepts a TransactionLogger type
//and return a pointer to a KeyStoreValue
//Attaching the drive at the port
func NewKeyValueStore(tl TransactionLogger) *KeyValueStore {
	return &KeyValueStore{
		m:        make(map[string]string),
		transact: tl,
	}
}
