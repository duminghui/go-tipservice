// Package db provides ...
package db

import (
	"time"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
)

type txProcessStatus int

const (
	TxProcessStatusWait = txProcessStatus(-1)
	// TxProcessStatusNoChange = txProcessStatus(0)
	TxProcessStatusDone = txProcessStatus(1)
)

func (db *DB) cTxProcess(session *mgo.Session) *mgo.Collection {
	return session.DB(db.database).C("tx_process_info")
}

// TxProcessInfo
// Status:
//  -1: wait process
//  1: process compile
type TxProcessInfo struct {
	Symbol      string          `bson:"symbol,omitempty"`
	TxID        string          `bson:"txid,omitempty"`
	Status      txProcessStatus `bson:"status,omitempty"`
	ProcessTime int64           `bson:"process_time,omitempty"`
}

func (db *DB) TxProcessIsDone(sessionUse *mgo.Session, txID string) (bool, error) {
	session, closer := session(sessionUse)
	defer closer()
	c := db.cTxProcess(session)
	// symbol := db.symbol
	selector := TxProcessInfo{
		Symbol: db.symbol,
		TxID:   txID,
	}
	result := new(TxProcessInfo)
	err := c.Find(selector).One(result)
	if err != nil && err != mgo.ErrNotFound {
		return false, err
	}
	if err != nil && err == mgo.ErrNotFound {
		return false, nil
	}
	if result.Status == TxProcessStatusDone {
		return true, nil
	}
	return false, nil
}

func (db *DB) txProcess(sessionOrg *mgo.Session, symbol, txID string) (*TxProcessInfo, error) {
	session, closer := session(sessionOrg)
	defer closer()
	selector := TxProcessInfo{
		Symbol: symbol,
		TxID:   txID,
	}
	c := db.cTxProcess(session)
	info := new(TxProcessInfo)
	err := c.Find(selector).One(info)
	if err != nil && err == mgo.ErrNotFound {
		return nil, nil
	}
	if err != nil && err != mgo.ErrNotFound {
		log.Errorf("[%s]txProcess Error:%s[%s]", symbol, err, txID)
		return nil, err
	}
	return info, nil
}

func (db *DB) TxProcessAddNew(symbol, txID string, extendSecond float64) error {
	session, closer := session(nil)
	defer closer()
	info, err := db.txProcess(session, symbol, txID)
	if err != nil {
		return err
	}
	if info != nil {
		return nil
	}
	processTime := time.Now().Add(time.Duration(extendSecond) * time.Second).Unix()
	data := &TxProcessInfo{
		Symbol:      symbol,
		TxID:        txID,
		Status:      TxProcessStatusWait,
		ProcessTime: processTime,
	}
	c := db.cTxProcess(session)
	err = c.Insert(data)
	if err != nil {
		log.Errorf("[%s]TxProcessAddNew Error:%s[%s]", symbol, err, txID)
		return err
	}
	return nil
}

func (db *DB) TxProcessExtendTime(sessionUse *mgo.Session, symbol, txID string, extendSecond int64) error {
	session, closer := session(sessionUse)
	defer closer()
	selector := TxProcessInfo{
		Symbol: symbol,
		TxID:   txID,
	}

	processTime := time.Now().Add(time.Duration(extendSecond) * time.Second).Unix()
	data := bson.M{
		"$set": &TxProcessInfo{
			ProcessTime: processTime,
		},
	}
	err := db.cTxProcess(session).Update(selector, data)
	if err != nil {
		log.Errorf("[%s]Update TxProcess Time Error: %s [%s]", symbol, err, txID)
		return err
	}
	return nil
}

func (db *DB) TxProcessStatusDone(sessionUse *mgo.Session, symbol, txID string) error {
	session, closer := session(sessionUse)
	defer closer()
	selector := TxProcessInfo{
		Symbol: symbol,
		TxID:   txID,
	}

	data := bson.M{
		"$set": &TxProcessInfo{
			Status: TxProcessStatusDone,
		},
	}
	err := db.cTxProcess(session).Update(selector, data)
	if err != nil {
		log.Errorf("[%s]Update TxProcess Time Error: %s [%s]", symbol, err, txID)
		return err
	}
	return nil
}

// func (db *DB) UpsertTxProcess(sessionUse *mgo.Session, symbol, txID string, status txProcessStatus, extendSecond int64) error {
// 	session, closer := session(sessionUse)
// 	defer closer()
// 	// symbol := db.symbol
// 	selector := TxProcessInfo{
// 		Symbol: symbol,
// 		TxID:   txID,
// 	}
// 	processTime := int64(0)
// 	if extendSecond != 0 {
// 		processTime = time.Now().Add(time.Duration(extendSecond) * time.Second).Unix()
// 	}
// 	data := bson.M{
// 		"$set": &TxProcessInfo{
// 			Symbol:      symbol,
// 			TxID:        txID,
// 			Status:      status,
// 			ProcessTime: processTime,
// 		},
// 	}
// 	_, err := session.DB(db.database).C(colTxProcessInfo).Upsert(selector, data)
// 	if err != nil {
// 		log.Errorf("[%s]Upsert Process Error: %s [%s]", symbol, err, txID)
// 		return err
// 	}
// 	return nil
// }

func (db *DB) TxProcessInfos() ([]*TxProcessInfo, error) {
	session := mgoSession.Clone()
	defer session.Close()
	nowTime := time.Now().Unix()
	selector := bson.M{
		"status": -1,
		// "status": bson.M{
		// "$in": []int{-1, 1},
		// },
		"process_time": bson.M{
			"$lt": nowTime,
		},
	}
	query := db.cTxProcess(session).Find(selector).Limit(10)
	txs := make([]*TxProcessInfo, 0)
	err := query.All(&txs)
	if err != nil {
		log.Errorf("[All]Read Process Tx Count Error: %s", err)
		return nil, err
	}
	// log.Info(txs)

	return txs, nil
}
