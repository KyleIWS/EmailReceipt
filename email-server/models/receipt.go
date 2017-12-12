package models

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

type Receipt struct {
	ReceiptID bson.ObjectId
	Created   time.Time
	Read      time.Time
}

func NewReceipt() *Receipt {
	return &Receipt{
		ReceiptID: bson.NewObjectId(),
		Created:   time.Now(),
	}
}
