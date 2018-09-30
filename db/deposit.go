// Package db provides ...
package db

import (
	"github.com/duminghui/go-tipservice/amount"
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
)

func (db *DB) Deposit(address, txid string, amountF float64, isConfirmed bool) error {
	session := mgoSession.Clone()
	defer session.Close()
	user, err := db.UserByAddress(session, address)
	if err != nil {
		return err
	}
	if user == nil {
		db.saveNoOwnerDeposit(session, txid, address, amountF)
		db.TxProcessUpdate(session, txid, 1, 0)
		return nil
	}
	depositQuery := &Deposit{
		TxID:    txid,
		Address: address,
	}
	query := session.DB(db.database).C(colDeposit).Find(depositQuery)
	deposit := new(Deposit)
	err = query.One(deposit)
	if err != nil && err != mgo.ErrNotFound {
		return err
	}
	if err != nil && err == mgo.ErrNotFound {
		depositData := &Deposit{
			UserID:      user.UserID,
			Amount:      amountF,
			TxID:        txid,
			Address:     address,
			IsConfirmed: isConfirmed,
		}
		err = session.DB(db.database).C(colDeposit).Insert(depositData)
		if err != nil {
			return err
		}

		selector := &User{Address: address}
		var update *User
		addAmount, _ := amount.FromFloat64(amountF)
		if !isConfirmed {
			unconfirmedAmount := user.UnconfirmedAmount.Add(addAmount)
			update = &User{UnconfirmedAmount: unconfirmedAmount}
		} else {
			confirmedAmount := user.Amount.Add(addAmount)
			update = &User{Amount: confirmedAmount}
		}

		err = session.DB(db.database).C(colUser).Update(selector, bson.M{
			"$set": update,
		})
		if err != nil {
			return err
		}
		if !isConfirmed {
			db.TxProcessUpdate(session, txid, 0, 60)
		} else {
			db.TxProcessUpdate(session, txid, 1, 0)
		}
		return nil
	}

	if !deposit.IsConfirmed {
		if isConfirmed {
			err = session.DB(db.database).C(colDeposit).Update(depositQuery,
				bson.M{
					"$set": &Deposit{IsConfirmed: true},
				})
			if err != nil {
				return err
			}
			addAmount, _ := amount.FromFloat64(amountF)
			confirmedAmount := user.Amount.Add(addAmount)
			unconfirmedAmount := user.UnconfirmedAmount.Sub(addAmount)
			err = session.DB(db.database).C(colUser).Update(
				&User{
					UserID: user.UserID,
				},
				bson.M{"$set": &User{
					Amount:            confirmedAmount,
					UnconfirmedAmount: unconfirmedAmount,
				}})
			if err != nil {
				return err
			}
			db.TxProcessUpdate(session, txid, 1, 0)
		} else {
			db.TxProcessUpdate(session, txid, 0, 60)
		}
	} else {
		db.TxProcessUpdate(session, txid, 1, 0)
	}
	// fmt.Println("User", user)
	return nil
}

func (db *DB) saveNoOwnerDeposit(sessionUse *mgo.Session, txid, address string, amount float64) {
	session, closer := session(sessionUse)
	defer closer()
	data := &NoOwnerDeposit{
		TxID:    txid,
		Address: address,
		Amount:  amount,
	}
	changeInfo, err := session.DB(db.database).C(colNoOwnerDeposit).Upsert(data, data)
	if err != nil {
		log.Error("SaveNoOwnerDeposit Error:", err)
		return
	}
	log.Infof("SaveNoOwnerDeposit: %+v\n", changeInfo)
}
