// Package db provides ...
package db

import (
	"github.com/duminghui/go-tipservice/amount"
	"github.com/duminghui/go-util/utime"
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
)

func (db *DB) cDeposit(session *mgo.Session) *mgo.Collection {
	return session.DB(db.database).C("deposit")
}

// type Deposit struct {
// UserID      string  `bson:"user_id"`
// UserName    string  `bson:"user_name"`
// Amount      float64 `bson:"amount"`
// TxID        string  `bson:"txid"`
// Address     string  `bson:"addresses"`
// IsConfirmed bool    `bson:"isConfirmed"`
// }

type Deposit struct {
	UserID      string  `bson:"user_id,omitempty"`
	UserName    string  `bson:"user_name,omitempty"`
	Amount      float64 `bson:"amount,omitempty"`
	TxID        string  `bson:"txid,omitempty"`
	Address     string  `bson:"addresses,omitempty"`
	Time        string  `bson:"time,omitempty"`
	IsConfirmed bool    `bson:"isConfirmed,omitempty"`
}

type NoOwnerDeposit struct {
	TxID    string  `bson:"txid,omitempty"`
	Address string  `bson:"address,omitempty"`
	Amount  float64 `bson:"amount,omitempty"`
	Time    string  `bson:"time,omitempty"`
}

func (db *DB) cNoOwnerDeposit(session *mgo.Session) *mgo.Collection {
	return session.DB(db.database).C("no_owner_deposit")
}

func (db *DB) Deposit(address, txid string, time int64, amountF float64, isConfirmed bool) error {
	session := mgoSession.Clone()
	defer session.Close()
	user, err := db.userByAddress(session, address)
	if err != nil {
		return err
	}
	symbol := db.symbol
	if user == nil {
		db.saveNoOwnerDeposit(session, txid, address, amountF, time)
		db.TxProcessStatusDone(session, symbol, txid)
		return nil
	}
	depositQuery := &Deposit{
		TxID:    txid,
		Address: address,
	}
	c := db.cDeposit(session)
	query := c.Find(depositQuery)
	deposit := new(Deposit)
	err = query.One(deposit)
	if err != nil && err != mgo.ErrNotFound {
		return err
	}
	if err != nil && err == mgo.ErrNotFound {
		depositData := &Deposit{
			UserID:      user.UserID,
			UserName:    user.UserName,
			Amount:      amountF,
			TxID:        txid,
			Address:     address,
			Time:        utime.FormatLongTimeStrUTC(time),
			IsConfirmed: isConfirmed,
		}
		err = c.Insert(depositData)
		if err != nil {
			return err
		}

		addAmount, _ := amount.FromFloat64(amountF)
		if !isConfirmed {
			db.TxProcessExtendTime(session, symbol, txid, 60)
			err = db.userUnconfirmedAmountAddUpsert(session, user.UserID, user.UserName, addAmount)
		} else {
			db.TxProcessStatusDone(session, symbol, txid)
			err = db.UserAmountAddUpsert(session, user.UserID, user.UserName, addAmount)
		}

		if err != nil {
			db.TxProcessExtendTime(session, symbol, txid, 60)
			return err
		}
		return nil
	}

	if !deposit.IsConfirmed {
		if isConfirmed {
			err = c.Update(depositQuery,
				bson.M{
					"$set": &Deposit{
						IsConfirmed: true,
						Time:        utime.FormatLongTimeStrUTC(time),
					},
				})
			if err != nil {
				return err
			}
			amountConfirmed, _ := amount.FromFloat64(amountF)
			db.userConfirmedAmount(user.UserID, amountConfirmed)
			if err != nil {
				return err
			}
			db.TxProcessStatusDone(session, symbol, txid)
		} else {
			db.TxProcessExtendTime(session, symbol, txid, 60)
		}
	} else {
		db.TxProcessStatusDone(session, symbol, txid)
	}
	// fmt.Println("User", user)
	return nil
}

func (db *DB) saveNoOwnerDeposit(sessionUse *mgo.Session, txid, address string, amount float64, time int64) {
	session, closer := session(sessionUse)
	defer closer()
	selector := &NoOwnerDeposit{
		TxID:    txid,
		Address: address,
	}
	data := &NoOwnerDeposit{
		TxID:    txid,
		Address: address,
		Amount:  amount,
		Time:    utime.FormatLongTimeStrUTC(time),
	}
	changeInfo, err := db.cNoOwnerDeposit(session).Upsert(selector, data)
	if err != nil {
		log.Error("SaveNoOwnerDeposit Error:", err)
		return
	}
	log.Infof("SaveNoOwnerDeposit: %+v\n", changeInfo)
}
