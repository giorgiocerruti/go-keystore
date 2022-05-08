# go-keystore
A simple key-value store for cloud native go book 


### Usage

```
Usage:
  flags [flags]

Flags:
      --dbHost string       DB host (default "localhost")
      --dbName string       DB name (default "transactions")
      --dbPassword string   DB password (default "postgres")
      --dbTbName string     DB table name (default "transaction")
      --dbUser string       db username (default "postgres")
      --fAddress string     Frontend address to listen on (default "localhost")
      --fPort string        Frontend port
      --file string         file name to store transactions (default "transaction.log")
  -f, --frontend string     front-end type (default "rest")
  -h, --help                help for flags
  -t, --tlog string         transactionsl log type (default "file")
  ```