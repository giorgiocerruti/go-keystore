package db

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"

	"github.com/giorgiocerruti/go-keystore/pkg/logger"
)

type PostgresTransactionLogger struct {
	events chan<- logger.Event
	errors <-chan error
	db     *sql.DB
	dbConf PostgresDBParams
}

//Logs the PUT request
func (l *PostgresTransactionLogger) WritePut(key, value string) {
	l.events <- logger.Event{EventType: logger.EventPut, Key: key, Value: value}
}

//Logs the DELETE request
func (l *PostgresTransactionLogger) WriteDelete(key string) {
	l.events <- logger.Event{EventType: logger.EventDelete, Key: key}
}

//Return a channel of errors
func (l *PostgresTransactionLogger) Err() <-chan error {
	return l.errors
}

//Verify if the table exists
func (l *PostgresTransactionLogger) VerifyTableExists() (bool, error) {
	q := `SELECT EXISTS (
		SELECT FROM 
		    pg_tables
		WHERE 
		    schemaname = 'public' AND 
		    tablename  = '?'
		);`
	result, err := l.db.Exec(q, l.dbConf.tableName)
	if err != nil {
		return false, fmt.Errorf("error checking table exists")
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return false, fmt.Errorf("errors retrieving affected rows: %w", err)
	}

	if rows != 1 {
		return false, nil
	}

	return true, nil
}

//create a table
func (l *PostgresTransactionLogger) CreateTable() error {
	q := `CREATE TABLE ? (
		sequence serial PRIMARY KEY,
		eventType INT NOT NULL,
		key VARCHAR(255) NOT NULL,
		value VARCHAR(255)
		);`

	result, err := l.db.Exec(q, l.dbConf.tableName)
	if err != nil {
		return fmt.Errorf("error checking table exists")
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("errors retrieving affected rows: %w", err)
	}

	if rows != 1 {
		return fmt.Errorf("table has not created")
	}

	return nil
}

func (l *PostgresTransactionLogger) ReadEvents() (<-chan logger.Event, <-chan error) {
	outEvent := make(chan logger.Event)
	outError := make(chan error)
	q := `SELECT sequende, event_type, key. value FROM ? ORDER BY sequence`

	go func() {

		defer close(outEvent)
		defer close(outError)

		rows, err := l.db.Query(q, l.dbConf.tableName)
		if err != nil {
			outError <- fmt.Errorf("sql query error: %w", err)
			return
		}

		defer rows.Close()
		e := logger.Event{}

		for rows.Next() {
			err = rows.Scan(
				&e.Sequence,
				&e.EventType,
				&e.Key,
				&e.Value,
			)

			if err != nil {
				outError <- fmt.Errorf("error reading row: %w", err)
				return
			}

			outEvent <- e
		}

		err = rows.Err()
		if err != nil {
			outError <- fmt.Errorf("log read failure: %w", err)
		}
	}()

	return outEvent, outError
}

//Run the gorutinr for insert items into the DB
func (l *PostgresTransactionLogger) Run() {
	events := make(chan logger.Event, 16)
	l.events = events

	errors := make(chan error, 1)
	l.errors = errors

	go func() {
		q := `INSERT INTO $1 (event_type, key, value) VALUES ($2, $3, $4)`

		for e := range events {
			_, err := l.db.Exec(q, l.dbConf.tableName, e.EventType, e.Key, e.Value)
			if err != nil {
				errors <- err
			}
		}
	}()
}
func NewPostgresTransactionLogger(conf PostgresDBParams) (logger.TransactionLogger, error) {

	//Create connection string
	connString := fmt.Sprintf("host=%s dbName=%s user=%s password=%s",
		conf.host, conf.dbName, conf.user, conf.password)

	//Open connections
	db, err := sql.Open("postgres", connString)
	if err != nil {
		return nil, fmt.Errorf("failed to open db: %w", err)
	}

	//database/sql doesn't open connection
	//nedd to ping to force opne it
	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("failed to open db connection: %w", err)
	}

	log := &PostgresTransactionLogger{db: db}

	//Check if the db exists
	exists, err := log.VerifyTableExists()
	if err != nil {
		return nil, fmt.Errorf("failde to verify table exists: %w", err)
	}

	if !exists {
		//create the table is doesn't exist
		if err = log.CreateTable(); err != nil {
			return nil, fmt.Errorf("failde to create table: %w", err)
		}
	}

	return log, err
}
