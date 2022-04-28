package frontend

import "github.com/giorgiocerruti/go-keystore/core"

type Frontend interface {
	Start(s *core.KeyValueStore) error
}
