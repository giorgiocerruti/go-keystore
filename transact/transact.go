package transact

import (
	"fmt"

	"github.com/giorgiocerruti/go-keystore/core"
)

type TlConfig struct {
	Filename string
	DbConf   PostgresDBParams
}

func NewTransactionLogger(core string, conf TlConfig) (core.TransactionLogger, error) {
	switch core {
	case "file":
		return NewFileTransactionLogger(conf.Filename)
	case "postgres":
		return NewPostgresTransactionLogger(
			conf.DbConf,
		)
	case "":
		return nil, fmt.Errorf("transaction core must be speified")
	default:
		return nil, fmt.Errorf("no such trasnaction core %s", core)
	}
}
