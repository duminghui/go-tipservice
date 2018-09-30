// Package db provides ...
package db

import (
	"time"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
)

const (
	colTxProcessInfo = "tx_process_info"
)

// TxProcessInfo
// Status:
//  -1: wait process
//  1: process compile
type TxProcessInfo struct {
	Symbol      string `bson:"symbol,omitempty"`
	TxID        string `bson:"txid,omitempty"`
	Status      int    `bson:"status,omitempty"`
	ProcessTime int64  `bson:"process_time,omitempty"`
}

func (db *DB) SaveTxProcess(symbol, txid string) error {
	session := mgoSession.Clone()
	defer session.Close()
	// symbol := db.symbol
	selector := TxProcessInfo{
		Symbol: symbol,
		TxID:   txid,
	}
	data := TxProcessInfo{
		Symbol:      symbol,
		TxID:        txid,
		Status:      -1,
		ProcessTime: -1,
	}
	_, err := session.DB(db.database).C(colTxProcessInfo).Upsert(selector, data)
	if err != nil {
		log.Errorf("[%s]Save Process Error: %s [%s]", symbol, err, txid)
		return err
	}
	// log.Infof("[%s]Save Process Success:[%s] %#v", symbol, txid, changeInfo)
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

func (db *DB) TxProcessUpdate(sessionUse *mgo.Session, txid string, status int, extendSecond int64) error {
	session, closer := session(sessionUse)
	defer closer()
	symbol := db.symbol
	selector := &TxProcessInfo{
		Symbol: symbol,
		TxID:   txid,
	}
	processTime := int64(0)
	if extendSecond != 0 {
		processTime = time.Now().Add(time.Duration(extendSecond) * time.Second).Unix()
	}
	data := &TxProcessInfo{
		Status:      status,
		ProcessTime: processTime,
	}
	err := session.DB(db.database).C(colTxProcessInfo).Update(selector,
		bson.M{
			"$set": data,
		})
	if err != nil {
		log.Errorf("[%s]Update Process Tx Error:%s[%s]", symbol, err, txid)
		return err
	}
	return nil
}
