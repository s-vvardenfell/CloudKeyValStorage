package main

import (
	"log"

	"github.com/s-vvardenfell/CloudKeyValStorage/core"
	"github.com/s-vvardenfell/CloudKeyValStorage/frontend"
	"github.com/s-vvardenfell/CloudKeyValStorage/transact"
)

func main() {
	// Create our TransactionLogger. This is an adapter that will plug
	// into the core application's TransactionLogger plug.
	tl, _ := transact.NewTransactionLogger("file")

	// Create Core and tell it which TransactionLogger to use.
	// This is an example of a "driven agent"
	store := core.NewKeyValueStore().WithTransactionLogger(tl)
	store.Restore()

	// Create the frontend.
	// This is an example of a "driving agent".
	fe, _ := frontend.NewFrontEnd("rest")

	log.Fatal(fe.Start(store))
}
