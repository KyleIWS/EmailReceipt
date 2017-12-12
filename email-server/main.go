package main

import (
	"log"
	"net/http"
	"os"

	"github.com/KyleIWS/EmailReceipt/email-server/handlers"
	"github.com/KyleIWS/EmailReceipt/email-server/models"
	mgo "gopkg.in/mgo.v2"
)

const SVCADDR = "localhost:4444"
const DBADDR = "localhost:27017"

func main() {

	addr := os.Getenv("ADDR")
	if len(addr) == 0 {
		addr = SVCADDR
	}

	dbaddr := os.Getenv("DBADDR")
	if len(dbaddr) == 0 {
		dbaddr = DBADDR
	}

	session, err := mgo.Dial(dbaddr)
	if err != nil {
		log.Fatalf("error dialing mongo db: %v", err)
	}
	mongoStore := models.NewMongoStore(session, "ReceiptDB", "ReceiptCollectionTesting")

	ctx := handlers.NewReceiptCtx(mongoStore)

	mux := http.NewServeMux()
	// Handler is called on my end to create a new receipt
	mux.HandleFunc("/create", ctx.CreateReceiptHandler)
	// here we will serve up all the receipts and let the user go from there
	mux.HandleFunc("/all", ctx.GetAllReceiptsHandler)
	// handle deleting all of them for testing purposes
	mux.HandleFunc("/delete-all", ctx.DeleteAllHandler)
	// Need handler will be in charge of detecting when an image asset is requested
	mux.HandleFunc("/static/", ctx.ServeFile)

	// Add CORS stuff
	corsHandler := handlers.NewCORS(mux)

	log.Printf("server started listeing on %s...", addr)
	log.Fatal(http.ListenAndServe(addr, corsHandler))
}
