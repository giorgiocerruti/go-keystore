package main

import (
	"fmt"
	"log"
	"os"

	"github.com/giorgiocerruti/go-keystore/core"
	"github.com/giorgiocerruti/go-keystore/frontend"
	v1 "github.com/giorgiocerruti/go-keystore/frontend/rest/v1"
	"github.com/giorgiocerruti/go-keystore/transact"
	"github.com/spf13/cobra"
)

var tlType string
var tlFrontend string
var address string
var port string
var conf transact.TlConfig

var rootCmd = &cobra.Command{
	Use:  "flags",
	Long: "Flags to enable/disble options",
	Run:  startFunc,
}

func init() {
	rootCmd.Flags().StringVarP(&tlType, "tlog", "t", "file", "transactionsl log type")
	rootCmd.Flags().StringVarP(&tlFrontend, "frontend", "f", "rest", "front-end type")
	rootCmd.Flags().StringVarP(&conf.Filename, "file", "", "transaction.log", "file name to store transactions")
	rootCmd.Flags().StringVarP(&conf.DbConf.DbName, "dbName", "", "transactions", "DB name")
	rootCmd.Flags().StringVarP(&conf.DbConf.Host, "dbHost", "", "localhost", "DB host")
	rootCmd.Flags().StringVarP(&conf.DbConf.User, "dbUser", "", "postgres", "db username")
	rootCmd.Flags().StringVarP(&conf.DbConf.Password, "dbPassword", "", "postgres", "DB password")
	rootCmd.Flags().StringVarP(&conf.DbConf.TableName, "dbTbName", "", "transaction", "DB table name")
	rootCmd.Flags().StringVarP(&address, "fAddress", "", "localhost", "Frontend address to listen on")
	rootCmd.Flags().StringVarP(&port, "fPort", "", "", "Frontend port")
}

func startFunc(cmd *cobra.Command, args []string) {
	//Create our TransactionLogger
	fmt.Printf("Backend type: %s\n", tlType)
	fmt.Printf("Frontend type: %s\n", tlFrontend)

	tl, err := transact.NewTransactionLogger(tlType, conf)
	if err != nil {
		log.Fatal(err)
	}

	//Creare the core and tell it wich TL to use
	store := core.NewKeyValueStore(tl)
	err = store.Restore()
	if err != nil {
		log.Fatal(err)
	}

	//Create the frontend
	fe, err := frontend.NewFrontEnd(tlFrontend)
	fe.(*v1.RestFrontend).Config.Address = address
	fe.(*v1.RestFrontend).Config.Port = port

	if err != nil {
		log.Fatal(err)
	}

	log.Fatal(fe.Start(store))
}
func main() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

}
