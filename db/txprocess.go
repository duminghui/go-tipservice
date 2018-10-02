// Package db provides ...
package db

import (
	"time"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
)

type txProcessStatus int

const (
	colTxProcessInfo        = "tx_process_info"
	TxProcessStatusWait     = txProcessStatus(-1)
	TxProcessStatusNoChange = txProcessStatus(0)
	TxProcessStatusDone     = txProcessStatus(1)
)

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

func (db *DB) IsTxProcessDone(sessionUse *mgo.Session, txID string) (bool, error) {
	session, closer := session(sessionUse)
	defer closer()
	c := session.DB(db.database).C(colTxProcessInfo)
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

func (db *DB) UpsertTxProcess(sessionUse *mgo.Session, symbol, txID string, status txProcessStatus, extendSecond int64) error {
	session, closer := session(sessionUse)
	defer closer()
	// symbol := db.symbol
	selector := TxProcessInfo{
		Symbol: symbol,
		TxID:   txID,
	}
	processTime := int64(0)
	if extendSecond != 0 {
		processTime = time.Now().Add(time.Duration(extendSecond) * time.Second).Unix()
	}
	data := TxProcessInfo{
		Symbol:      symbol,
		TxID:        txID,
		Status:      status,
		ProcessTime: processTime,
	}
	_, err := session.DB(db.database).C(colTxProcessInfo).Upsert(selector, data)
	if err != nil {
		log.Errorf("[%s]Upsert Process Error: %s [%s]", symbol, err, txID)
		return err
	}
	return nil
}

func (db *DB) TxProcessInfos() ([]TxProcessInfo, error) {
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
	query := session.DB(db.database).C(colTxProcessInfo).Find(selector).Limit(10)
	txs := make([]TxProcessInfo, 0)
	err := query.All(&txs)
	if err != nil {
		log.Errorf("[All]Read Process Tx Count Error: %s", err)
		return nil, err
	}
	// log.Info(txs)

	return txs, nil
}

// func (db *DB) TxProcessUpdate(sessionUse *mgo.Session, txid string, status txProcessStatus, extendSecond int64) error {
// 	session, closer := session(sessionUse)
// 	defer closer()
// 	symbol := db.symbol
// 	selector := &TxProcessInfo{
// 		Symbol: symbol,
// 		TxID:   txid,
// 	}
// 	processTime := int64(0)
// 	if extendSecond != 0 {
// 		processTime = time.Now().Add(time.Duration(extendSecond) * time.Second).Unix()
// 	}
// 	data := &TxProcessInfo{
// 		Status:      status,
// 		ProcessTime: processTime,
// 	}
// 	err := session.DB(db.database).C(colTxProcessInfo).Update(selector,
// 		bson.M{
// 			"$set": data,
// 		})
// 	if err != nil {
// 		log.Errorf("[%s]Update Process Tx Error:%s[%s]", symbol, err, txid)
// 		return err
// 	}
// 	return nil
// }
