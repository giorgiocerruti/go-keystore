package v1

import "errors"

var store = make(map[string]string)

var ErrorNoSuchKey = errors.New("no such key")

func Put(key, value string) error {
	store[key] = value

	return nil
}

func Get(key string) (string, error) {
	v, ok := store[key]

	if !ok {
		return "", ErrorNoSuchKey
	}

	return v, nil
}

func Delete(key string) error {
	delete(store, key)

	return nil
}
