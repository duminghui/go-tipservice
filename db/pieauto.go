// Package db provides ...
package db

import (
	"errors"
	"time"

	"github.com/duminghui/go-tipservice/amount"
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
)

type PieAuto struct {
	ID           bson.ObjectId `bson:"_id"`
	GuildID      string        `bson:"guildid"`
	GuildName    string        `bson:"guildname"`
	UserID       string        `bson:"userid"`
	UserName     string        `bson:"username"`
	Symbol       string        `bson:"symbol"`
	ChannelID    string        `bson:"channelid"`
	ChannelName  string        `bson:"channelName"`
	RoleID       string        `bson:"roleid"`
	CreateTime   time.Time     `bson:"createtime"`
	CycleTimes   int64         `bosn:"cycletimes"`
	IntervalTime time.Duration `bson:"intervaltime"`
	Amount       amount.Amount `bson:"amount"`
	AmountF      float64       `bson:"amount_f"`
	IsOnlineUser bool          `bson:"isonlineuser"`
	NextPieTime  time.Time     `bson:"nextpietime"`
	RunnedTimes  int64         `bson:"runnedtimes"`
	IsEnd        bool          `bson:"isend"`
}

type pieAutoOption struct {
	ChannelName string    `bson:"channelName,omitempty"`
	NextPieTime time.Time `bson:"nextpietime,omitempty"`
	RunnedTimes int64     `bson:"runnedtimes,omitempty"`
	IsEnd       bool      `bson:"isend,omitempty"`
}

type pieAutoSelector struct {
	Symbol string        `bson:"symbol,omitempty"`
	ID     bson.ObjectId `bson:"_id,omitempty"`
	UserID string        `bson:"userid,omitempty"`
}

func (db *DBGuild) cPieAuto(session *mgo.Session) *mgo.Collection {
	return session.DB(db.database).C("pieauto")
}

func (db *DBGuild) PieAutoAdd(guildID, guildName, userID, userName, symbol, channelID, channelName, roleID string, cycleTimes int64, intervalTime time.Duration, amount amount.Amount, isOnlineUser bool, extendTime time.Duration) (*PieAuto, error) {
	session := mgoSession.Clone()
	defer session.Close()
	id := bson.NewObjectId()
	createTime := time.Now()
	nextPieTime := createTime.Add(extendTime)
	data := &PieAuto{
		ID:           id,
		GuildID:      guildID,
		GuildName:    guildName,
		UserID:       userID,
		UserName:     userName,
		Symbol:       symbol,
		ChannelID:    channelID,
		ChannelName:  channelName,
		RoleID:       roleID,
		CreateTime:   createTime,
		CycleTimes:   cycleTimes,
		IntervalTime: intervalTime,
		Amount:       amount,
		AmountF:      amount.Float64(),
		IsOnlineUser: isOnlineUser,
		NextPieTime:  nextPieTime,
		RunnedTimes:  0,
		IsEnd:        false,
	}
	col := db.cPieAuto(session)
	err := col.Insert(data)
	if err != nil {
		log.Errorf("[%s]AutoPieAdd Error:%s[%s(%s)][%s(%s)]", symbol, err, guildID, guildName, userID, userName)
	}
	return data, nil
}

var NotIDError = errors.New("Not Bson ID Error")
var PieAutoOwnerError = errors.New("Not PieAuto owner")

func (db *DBGuild) PieAutoRemove(userID, id string) error {
	if !bson.IsObjectIdHex(id) {
		return NotIDError
	}
	session := mgoSession.Clone()
	defer session.Close()
	selector := &pieAutoSelector{
		ID:     bson.ObjectIdHex(id),
		UserID: userID,
	}
	err := db.cPieAuto(session).Remove(selector)
	if err != nil {
		if err == mgo.ErrNotFound {
			return nil
		}
		log.Errorf("PieAutoRemove Error:%s[uID:%s,id:%s]", err, userID, id)
		return err
	}
	return nil
}

func (db *DBGuild) PieAutoUpdateChannelName(id, channelName string) error {
	session := mgoSession.Clone()
	defer session.Close()
	data := &pieAutoOption{
		ChannelName: channelName,
	}
	col := db.cPieAuto(session)
	err := col.UpdateId(bson.ObjectIdHex(id), bson.M{
		"$set": data,
	})
	return err
}

func (db *DBGuild) PieAutoUpTimes(id string) error {
	session := mgoSession.Clone()
	defer session.Close()
	autopie, err := db.PieAutoByID(session, id)
	if err != nil {
		log.Errorf("AutoPieUpdateTimes Error:%s(%s)", err, id)
		return err
	}
	runnedTimes := autopie.RunnedTimes + 1
	nextPieTime := autopie.NextPieTime.Add(autopie.IntervalTime)
	isEnd := runnedTimes >= autopie.CycleTimes
	data := &pieAutoOption{
		RunnedTimes: runnedTimes,
		NextPieTime: nextPieTime,
		IsEnd:       isEnd,
	}
	col := db.cPieAuto(session)
	err = col.UpdateId(bson.ObjectIdHex(id), bson.M{
		"$set": data,
	})
	return err
}

func (db *DBGuild) PieAutoProcessLists(size int) ([]*PieAuto, error) {
	session := mgoSession.Clone()
	defer session.Close()
	autoPieList := make([]*PieAuto, 0)
	selector := bson.M{
		"isend": false,
		"nextpietime": bson.M{
			"$lte": time.Now(),
		}}
	query := db.cPieAuto(session).Find(selector)
	query.Sort("nextpietime")
	query.Limit(size)
	err := query.All(&autoPieList)
	if err != nil {
		log.Errorf("AutoPieProcessLists Errors:%s", err)
		return nil, err
	}
	return autoPieList, nil
}

func (db *DBGuild) PieAutoLists(symbol string, start, size int) ([]*PieAuto, error) {
	session := mgoSession.Clone()
	defer session.Close()
	autoPieList := make([]*PieAuto, 0)
	selector := &pieAutoSelector{
		Symbol: symbol,
	}
	query := db.cPieAuto(session).Find(selector)
	query.Sort("nextpietime")
	query.Skip(start).Limit(size)
	err := query.All(&autoPieList)
	if err != nil {
		log.Errorf("AutoPieLists Errors:%s", err)
		return nil, err
	}
	return autoPieList, nil
}

func (db *DBGuild) PieAutoByID(sessionIn *mgo.Session, id string) (*PieAuto, error) {
	if !bson.IsObjectIdHex(id) {
		return nil, NotIDError
	}
	session, closer := session(sessionIn)
	defer closer()
	autopie := new(PieAuto)
	col := db.cPieAuto(session)
	err := col.FindId(bson.ObjectIdHex(id)).One(autopie)
	if err != nil {
		if err == mgo.ErrNotFound {
			return nil, nil
		}
		log.Errorf("AutoPieByID Error:%s(%s)", err, id)
		return nil, err
	}
	return autopie, nil
}
