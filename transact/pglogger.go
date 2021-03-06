package transact

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"

	"github.com/giorgiocerruti/go-keystore/core"
)

type PostgresDBParams struct {
	DbName    string
	Host      string
	User      string
	Password  string
	TableName string
}

type PostgresTransactionLogger struct {
	events chan<- core.Event
	errors <-chan error
	db     *sql.DB
	dbConf PostgresDBParams
}

//Logs the PUT request
func (l *PostgresTransactionLogger) WritePut(key, value string) {
	l.events <- core.Event{EventType: core.EventPut, Key: key, Value: value}
}

//Logs the DELETE request
func (l *PostgresTransactionLogger) WriteDelete(key string) {
	l.events <- core.Event{EventType: core.EventDelete, Key: key}
}

//Return a channel of errors
func (l *PostgresTransactionLogger) Err() <-chan error {
	return l.errors
}

//create a table
func (l *PostgresTransactionLogger) CreateTable() error {
	q := fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s (
		sequence serial PRIMARY KEY,
		eventType INT NOT NULL,
		key VARCHAR(255) NOT NULL,
		value VARCHAR(255)
		);`, l.dbConf.TableName)

	result, err := l.db.Exec(q)
	if err != nil {
		return fmt.Errorf("error creating the table: %w", err)
	}

	_, err = result.RowsAffected()
	if err != nil {
		return fmt.Errorf("errors retrieving affected rows: %w", err)
	}

	return nil
}

func (l *PostgresTransactionLogger) Read() (<-chan core.Event, <-chan error) {
	outEvent := make(chan core.Event)
	outError := make(chan error)
	q := fmt.Sprintf("SELECT sequence, eventType, key, value FROM %s ORDER BY sequence", l.dbConf.TableName)

	go func() {

		defer close(outEvent)
		defer close(outError)

		rows, err := l.db.Query(q)
		if err != nil {
			outError <- fmt.Errorf("sql query error: %w", err)
			return
		}

		defer rows.Close()
		e := core.Event{}
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
	events := make(chan core.Event, 16)
	l.events = events

	errors := make(chan error, 1)
	l.errors = errors

	go func() {
		var q string

		for e := range events {
			q = fmt.Sprintf(`INSERT INTO %s (eventType, key, value) VALUES (%d, '%s', '%s')`, l.dbConf.TableName, e.EventType, e.Key, e.Value)
			fmt.Println(q)
			rows, err := l.db.Exec(q)
			if err != nil {
				errors <- err
				return
			}

			result, err := rows.RowsAffected()
			if err != nil {
				errors <- fmt.Errorf("error insert: %w", err)
				return
			}

			if result != 1 {
				errors <- fmt.Errorf("no rows affected")
			}
		}
	}()
}

func NewPostgresTransactionLogger(conf PostgresDBParams) (core.TransactionLogger, error) {

	//Create connection string
	connString := fmt.Sprintf("host=%s dbname=%s user=%s password=%s sslmode=disable",
		conf.Host, conf.DbName, conf.User, conf.Password)

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

	log := &PostgresTransactionLogger{db: db, dbConf: conf}

	if err = log.CreateTable(); err != nil {
		return nil, fmt.Errorf("failde to create table: %w", err)
	}

	return log, err
}
