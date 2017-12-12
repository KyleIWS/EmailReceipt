package main

import (
	"encoding/json"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"log"
	"net/http"
	"os"

	"github.com/KyleIWS/EmailReceipt/email-server/models"
	mgo "gopkg.in/mgo.v2"
)

const SVCADDR = "localhost:4444"
const DBADDR = "localhost:27017"
const DirReceiptPNG = "./pngs/"

func NewReceiptCtx(ms *models.MongoStore) *ReceiptCtx {
	return &ReceiptCtx{
		ms: ms,
	}
}

type ReceiptCtx struct {
	ms *models.MongoStore
}

// TODO: Add ability to parse incoming JSON to discern extra qualities about a new
// receipt such as a subject or recipient.
func (ctx *ReceiptCtx) CreateReceiptHandler(w http.ResponseWriter, r *http.Request) {
	// Create a new receipt ObjectId
	rpt := models.NewReceipt()
	// Commit it to the database
	if err := ctx.ms.Insert(rpt); err != nil {
		http.Error(w, fmt.Sprintf("error adding new receipt to database: %v", err), http.StatusInternalServerError)
		return
	}
	// Create a static asset (png file) somewhere that can be GET requested
	img := image.NewNRGBA(image.Rect(0, 0, 1, 1))
	img.Set(1, 1, color.NRGBA{
		R: uint8(255),
		G: uint8(255),
		B: uint8(255),
	})
	f, err := os.Create(DirReceiptPNG + string(rpt.ReceiptID))
	if err != nil {
		http.Error(w, fmt.Sprintf("error opening new file to write to: %v", err), http.StatusInternalServerError)
		return
	}
	if err := png.Encode(f, img); err != nil {
		f.Close()
		http.Error(w, fmt.Sprintf("error encoding new img to png: %v", err), http.StatusInternalServerError)
		return
	}
	f.Close()
	// could return id but why bother
	//json.NewEncoder(w).Encode()
}

func (ctx *ReceiptCtx) GetAllReceiptsHandler(w http.ResponseWriter, r *http.Request) {
	allReceipts, err := ctx.ms.GetAllReceipts()
	if err != nil {
		http.Error(w, fmt.Sprintf("error retrieving all receipts: %v", err), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(allReceipts)
}

func main() {

	addr := os.Getenv("SVCADDER")
	if len("SVCADDER") == 0 {
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

	ctx := NewReceiptCtx(mongoStore)

	mux := http.NewServeMux()
	// Handler is called on my end to create a new receipt
	mux.HandleFunc("/create", ctx.CreateReceiptHandler)
	// here we will serve up all the receipts and let the user go from there
	mux.HandleFunc("/all", ctx.GetAllReceiptsHandler)
	// Need handler will be in charge of detecting when an image asset is requested
	// mux.Handle("/static/", )

	log.Printf("server started listeing on %s...", addr)
	log.Fatal(http.ListenAndServe(addr, mux))
}
