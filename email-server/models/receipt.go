package models

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

type Receipt struct {
	ReceiptID bson.ObjectId `json:"id" bson:"_id"`
	Created   time.Time     `json:"created"`
	Reads     []time.Time   `json:"reads"`
}

func NewReceipt() *Receipt {
	return &Receipt{
		ReceiptID: bson.NewObjectId(),
		Created:   time.Now(),
		Reads:     make([]time.Time, 0),
	}
}
