package handlers

import (
	"encoding/json"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/KyleIWS/EmailReceipt/email-server/models"
	"gopkg.in/mgo.v2/bson"
)

const DirReceiptPNG = "./pngs/"

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

	width := 1
	height := width

	// Create a static asset (png file) somewhere that can be GET requested
	img := image.NewNRGBA(image.Rect(0, 0, width, height))
	for w := 0; w < width; w++ {
		randColor := color.NRGBA{
			R: uint8(255),
			G: uint8(255),
			B: uint8(255),
		}
		for h := 0; h < height; h++ {
			img.Set(w, h, randColor)
		}
	}

	f, err := os.Create(DirReceiptPNG + string(rpt.ReceiptID.Hex()) + ".png")
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

func (ctx *ReceiptCtx) DeleteAllHandler(w http.ResponseWriter, r *http.Request) {
	err := ctx.ms.DeleteAll()
	if err != nil {
		http.Error(w, fmt.Sprintf("error deleting all receipts: %v", err), http.StatusInternalServerError)
		return
	}
}

func (ctx *ReceiptCtx) ServeFile(w http.ResponseWriter, r *http.Request) {
	pngPath := r.URL.RequestURI()
	splitPath := strings.Split(pngPath, "/")
	path := splitPath[len(splitPath)-1]
	id := strings.Split(path, ".")[0]
	bsonID := bson.ObjectIdHex(id)
	err := ctx.ms.SetRead(bsonID)
	if err != nil {
		log.Printf("error trying to set given receipt id to read: %v", err)
		http.Error(w, fmt.Sprintf("error trying to fetch image"), http.StatusNotFound)
	}
	http.ServeFile(w, r, DirReceiptPNG+"/"+path)
}
