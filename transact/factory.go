package transact

import (
	"fmt"
	"os"

	"github.com/s-vvardenfell/CloudKeyValStorage/core"
)

func NewTransactionLogger(logger string) (core.TransactionLogger, error) {
	switch logger {
	case "file":
		return NewFileTransactionLogger(os.Getenv("TLOG_FILENAME"))
	case "postgres":
		return NewPostgresTransactionLogger(
			PostgresDbParams{
				dbName:   os.Getenv("TLOG_DB_HOST"),
				host:     os.Getenv("TLOG_DB_DATABASE"),
				user:     os.Getenv("TLOG_DB_USERNAME"),
				password: os.Getenv("TLOG_DB_PASSWORD"),
			},
		)
	case "":
		return nil, fmt.Errorf("transaction logger type not defined")
	default:
		return nil, fmt.Errorf("no such transaction logger %s", logger)
	}
}
