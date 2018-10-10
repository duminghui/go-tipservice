// Package db provides ...
package db

import (
	"fmt"

	"github.com/duminghui/go-tipservice/amount"
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
)

func (db *DB) cUser(session *mgo.Session) *mgo.Collection {
	return session.DB(db.database).C("user")
}

type User struct {
	UserID             string        `bson:"user_id"`
	UserName           string        `bson:"user_name"`
	Address            string        `bson:"address"`
	Amount             amount.Amount `bson:"amount"`
	AmountF            float64       `bson:"amount_f"`
	UnconfirmedAmount  amount.Amount `bson:"unconfirmed_amount"`
	UnconfirmedAmountF float64       `bson:"unconfirmed_amount_f"`
}

type userSelector struct {
	UserID  string `bson:"user_id,omitempty"`
	Address string `bson:"address,omitempty"`
}

type userAddress struct {
	UserName string `bson:"user_name"`
	Address  string `bson:"address"`
}

type userAmount struct {
	UserName string        `bson:"user_name"`
	Amount   amount.Amount `bson:"amount"`
	AmountF  float64       `bson:"amount_f"`
}

type userUnconfirmedAmount struct {
	UserName           string        `bson:"user_name"`
	UnconfirmedAmount  amount.Amount `bson:"unconfirmed_amount"`
	UnconfirmedAmountF float64       `bson:"unconfirmed_amount_f"`
}

type userAllAmount struct {
	Amount             amount.Amount `bson:"amount"`
	AmountF            float64       `bson:"amount_f"`
	UnconfirmedAmount  amount.Amount `bson:"unconfirmed_amount"`
	UnconfirmedAmountF float64       `bson:"unconfirmed_amount_f"`
}

func (db *DB) userByAddress(sessionIn *mgo.Session, address string) (*User, error) {
	session, closer := session(sessionIn)
	defer closer()
	query := db.cUser(session).Find(&userSelector{Address: address})
	user := new(User)
	err := query.One(user)
	if err != nil {
		if err != mgo.ErrNotFound {
			log.Errorf("[%s]userByAddress Error:%s", db.symbol, err)
			return nil, err
		}
		return nil, nil
	}
	return user, nil
}

func (db *DB) UserByID(sessionIn *mgo.Session, userID string) (*User, error) {
	session, closer := session(sessionIn)
	defer closer()
	query := db.cUser(session).Find(&userSelector{UserID: userID})
	user := new(User)
	err := query.One(user)
	if err != nil {
		if err != mgo.ErrNotFound {
			log.Errorf("[%s]UserByID Error:%s", db.symbol, err)
			return nil, err
		}
		return nil, nil
	}
	return user, nil
}

func (db *DB) UserAmountSub(sessionIn *mgo.Session, userID, userName string, amountSub amount.Amount) error {
	session, closer := session(sessionIn)
	defer closer()
	user, err := db.UserByID(session, userID)
	if err != nil {
		return err
	}
	if user == nil {
		log.Errorf("UserAmountSub NoUserByID:%s", userID)
		return fmt.Errorf("UserAmountSub NoUserByID:%s", userID)
	}
	col := db.cUser(session)
	confirmedAmount := user.Amount.Sub(amountSub)
	selector := &userSelector{
		UserID: userID,
	}
	data := bson.M{
		"$set": &userAmount{
			UserName: userName,
			Amount:   confirmedAmount,
			AmountF:  confirmedAmount.Float64(),
		},
	}
	err = col.Update(selector, data)
	if err != nil {
		log.Errorf("[%s]UserAmountAddUpsert Error:%s", db.symbol, err)
		return err
	}
	return nil
}

func (db *DB) UserAmountAddUpsert(sessionIn *mgo.Session, userID, userName string, amountAdd amount.Amount) error {
	session, closer := session(sessionIn)
	defer closer()
	user, err := db.UserByID(session, userID)
	if err != nil {
		return err
	}
	col := db.cUser(session)
	if user == nil {
		data := &User{
			UserID:             userID,
			UserName:           userName,
			Address:            "",
			Amount:             amountAdd,
			AmountF:            amountAdd.Float64(),
			UnconfirmedAmount:  amount.Zero,
			UnconfirmedAmountF: 0.0,
		}
		err := col.Insert(data)
		if err != nil {
			log.Errorf("[%s]UserAmountAddUpsert Error:%s", db.symbol, err)
			return err
		}
		return nil
	}

	confirmedAmount := user.Amount.Add(amountAdd)
	data := bson.M{
		"$set": &userAmount{
			UserName: userName,
			Amount:   confirmedAmount,
			AmountF:  confirmedAmount.Float64(),
		},
	}
	selector := &userSelector{
		UserID: userID,
	}
	err = col.Update(selector, data)
	if err != nil {
		log.Errorf("[%s]UserAmountAddUpsert Error:%s", db.symbol, err)
		return err
	}
	return nil
}

func (db *DB) userUnconfirmedAmountAddUpsert(sessionIn *mgo.Session, userID, userName string, amountAdd amount.Amount) error {
	session, closer := session(sessionIn)
	defer closer()
	user, err := db.UserByID(session, userID)
	if err != nil {
		return err
	}
	col := db.cUser(session)
	if user == nil {
		data := &User{
			UserID:             userID,
			UserName:           userName,
			Address:            "",
			Amount:             amount.Zero,
			AmountF:            0.0,
			UnconfirmedAmount:  amountAdd,
			UnconfirmedAmountF: amountAdd.Float64(),
		}
		err := col.Insert(data)
		if err != nil {
			log.Errorf("[%s]UserUnconfirmedAmountAddUpsert Error:%s", db.symbol, err)
			return err
		}
		return nil
	}

	unconfirmedAmount := user.UnconfirmedAmount.Add(amountAdd)
	data := bson.M{
		"$set": &userUnconfirmedAmount{
			UserName:           userName,
			UnconfirmedAmount:  unconfirmedAmount,
			UnconfirmedAmountF: unconfirmedAmount.Float64(),
		},
	}
	selector := &userSelector{
		UserID: userID,
	}
	err = col.Update(selector, data)
	if err != nil {
		log.Errorf("[%s]UserUnconfirmedAmountAddUpsert Error:%s", db.symbol, err)
		return err
	}
	return nil
}

func (db *DB) userConfirmedAmount(userID string, amountCfm amount.Amount) error {
	session := mgoSession.Clone()
	defer session.Close()
	user, err := db.UserByID(session, userID)
	if err != nil {
		return err
	}
	if user == nil {
		return fmt.Errorf("userConfirmedAmount no user:%s", userID)
	}
	confirmedAmount := user.Amount.Add(amountCfm)
	unconfirmedAmount := user.UnconfirmedAmount.Sub(amountCfm)
	col := db.cUser(session)
	err = col.Update(
		&userSelector{
			UserID: user.UserID,
		},
		bson.M{"$set": &userAllAmount{
			Amount:             confirmedAmount,
			AmountF:            confirmedAmount.Float64(),
			UnconfirmedAmount:  unconfirmedAmount,
			UnconfirmedAmountF: unconfirmedAmount.Float64(),
		}})
	if err != nil {
		return err
	}
	return nil
}

func (db *DB) UserAddressUpsert(userID, userName, address string, isInsert bool) error {
	session := mgoSession.Clone()
	defer session.Close()
	col := db.cUser(session)
	if isInsert {
		data := &User{
			UserID:             userID,
			UserName:           userName,
			Address:            address,
			Amount:             amount.Zero,
			AmountF:            0.0,
			UnconfirmedAmount:  amount.Zero,
			UnconfirmedAmountF: 0.0,
		}
		err := col.Insert(data)
		if err != nil {
			log.Errorf("[%s]UserAmountAddUpsert Error:%s", db.symbol, err)
			return err
		}
	} else {
		selector := &userSelector{
			UserID: userID,
		}
		data := bson.M{
			"$set": &userAddress{
				UserName: userName,
				Address:  address,
			},
		}
		err := col.Update(selector, data)
		if err != nil {
			return err
		}
	}
	return nil
}
