package models

import (
	"fmt"
	"time"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type MongoStore struct {
	session *mgo.Session
	dbname  string
	colname string
}

func NewMongoStore(sess *mgo.Session, dbName string, collectionName string) *MongoStore {
	if sess == nil {
		panic("nil pointer passed for session")
	}
	return &MongoStore{
		session: sess,
		dbname:  dbName,
		colname: collectionName,
	}
}

func (ms *MongoStore) GetByID(id bson.ObjectId) (*Receipt, error) {
	result := &Receipt{}
	col := ms.session.DB(ms.dbname).C(ms.colname)
	if err := col.Find(bson.M{"_id": id}).One(&result); err != nil {
		return nil, fmt.Errorf("error finding given receipt id: %v", err)
	}
	return result, nil
}

func (ms *MongoStore) Insert(rpt *Receipt) error {
	col := ms.session.DB(ms.dbname).C(ms.colname)
	if err := col.Insert(rpt); err != nil {
		return fmt.Errorf("error adding receipt to database: %v", err)
	}
	return nil
}

func (ms *MongoStore) SetRead(id bson.ObjectId) error {
	entry, err := ms.GetByID(id)
	if err != nil {
		return fmt.Errorf("error trying to set receipt to read: %v", err)
	}
	entry.Reads = append(entry.Reads, time.Now())
	change := mgo.Change{
		Update:    bson.M{"$set": entry},
		ReturnNew: true,
	}
	col := ms.session.DB(ms.dbname).C(ms.colname)
	if _, err := col.FindId(id).Apply(change, entry); err != nil {
		return fmt.Errorf("error updating receipt: %v", err)
	}
	return nil
}

func (ms *MongoStore) GetAllReceipts() ([]*Receipt, error) {
	results := []*Receipt{}
	col := ms.session.DB(ms.dbname).C(ms.colname)
	err := col.Find(bson.M{}).All(&results)
	return results, err
}

func (ms *MongoStore) Delete(id bson.ObjectId) error {
	col := ms.session.DB(ms.dbname).C(ms.colname)
	if err := col.RemoveId(id); err != nil {
		return fmt.Errorf("error deleting receipt: %v", err)
	}
	return nil
}

func (ms *MongoStore) DeleteAll() error {
	col := ms.session.DB(ms.dbname).C(ms.colname)
	if _, err := col.RemoveAll(nil); err != nil {
		return fmt.Errorf("error deleting receipt: %v", err)
	}
	return nil
}
