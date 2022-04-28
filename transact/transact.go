package transact

import (
	"fmt"
	"os"

	"github.com/giorgiocerruti/go-keystore/core"
)

func NewTransactionLogger(logger string) (core.TransactionLogger, error) {
	switch logger {
	case "file":
		return NewFileTransactionLogger(os.Getenv("TLOG_FILENAME"))
	case "postgres":
		return NewPostgresTransactionLogger(
			PostgresDBParams{
				dbName:   os.Getenv("TLOG_DB_DATABASE"),
				host:     os.Getenv("TGLOG_DB_HOST"),
				user:     os.Getenv("TGLOG_DB_USERNAME"),
				password: os.Getenv("TGLOG_DB_PASSWORD"),
			},
		)
	case "":
		return nil, fmt.Errorf("transaction logger must be speified")
	default:
		return nil, fmt.Errorf("no such trasnaction logger %s", logger)
	}
}
