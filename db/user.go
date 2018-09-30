// Package db provides ...
package db

import (
	"github.com/duminghui/go-tipservice/amount"
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
)

func (db *DB) UserByAddress(sessionIn *mgo.Session, address string) (*User, error) {
	session, closer := session(sessionIn)
	defer closer()
	query := session.DB(db.database).C(colUser).Find(&User{Address: address})
	user := new(User)
	err := query.One(user)
	if err != nil {
		if err != mgo.ErrNotFound {
			return nil, err
		}
		return nil, nil
	}
	return user, nil
}

func (db *DB) UserByID(sessionIn *mgo.Session, userID string) (*User, error) {
	session, closer := session(sessionIn)
	defer closer()
	query := session.DB(db.database).C(colUser).Find(&User{UserID: userID})
	user := new(User)
	err := query.One(user)
	if err != nil {
		if err != mgo.ErrNotFound {
			return nil, err
		}
		return nil, nil
	}
	return user, nil
}

func (db *DB) UserAmountUpsert(userID, userName string, amountF float64) error {
	session := mgoSession.Clone()
	defer session.Close()
	col := session.DB(db.database).C(colUser)
	selector := &User{
		UserID: userID,
	}
	user := new(User)
	err := col.Find(selector).One(user)
	if err != nil && err != mgo.ErrNotFound {
		return err
	}
	amountAdd, _ := amount.FromFloat64(amountF)
	if err != nil && err == mgo.ErrNotFound {
		data := &User{
			UserID:            userID,
			UserName:          userName,
			Address:           "",
			Amount:            amountAdd,
			UnconfirmedAmount: amount.Zero,
		}
		err := col.Insert(data)
		if err != nil {
			return err
		}
		return nil
	}

	confirmedAmount := user.Amount.Add(amountAdd)
	data := bson.M{
		"$set": &User{
			Amount: confirmedAmount,
		},
	}
	err = col.Update(selector, data)
	return err
}

func (db *DB) UserAddressUpsert(userID, userName, address string, isInsert bool) error {
	session := mgoSession.Clone()
	defer session.Close()
	col := session.DB(db.database).C(colUser)
	if isInsert {
		data := &User{
			UserID:            userID,
			UserName:          userName,
			Address:           address,
			Amount:            amount.Zero,
			UnconfirmedAmount: amount.Zero,
		}
		err := col.Insert(data)
		if err != nil {
			return err
		}
	} else {
		selector := &User{
			UserID: userID,
		}
		data := bson.M{
			"$set": &User{
				Address: address,
			},
		}
		err := col.Update(selector, data)
		if err != nil {
			return err
		}
	}
	return nil
}
