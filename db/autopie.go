// Package db provides ...
package db

import (
	"fmt"
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
	NextPieTime time.Time `bson:"netpietime,omitempty"`
	RunnedTimes int64     `bson:"runnedtimes,omitempty"`
	IsEnd       bool      `bson:"isend,omitempty"`
}

type pieAutoSelector struct {
	Symbol string `bson:"symbol,omitempty"`
}

func (db *DBGuild) cPieAuto(session *mgo.Session) *mgo.Collection {
	return session.DB(db.database).C("pieauto")
}

func (db *DBGuild) PieAutoAdd(guildID, guildName, userID, userName, symbol, channelID, roleID string, cycleTimes int64, intervalTime time.Duration, amount amount.Amount, isOnlineUser bool, extendTime time.Duration) (*PieAuto, error) {
	id := bson.NewObjectId()
	session := mgoSession.Clone()
	defer session.Close()
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
	// col := cPieAuto(session)
	// err := col.Insert(data)
	// if err != nil {
	// 	log.Errorf("[%s]AutoPieAdd Error:%s[%s(%s)][%s(%s)]", symbol, err, guildID, guildName, userID, userName)
	// }
	return data, nil
}

func (db *DBGuild) PieAutoRemove(id string) error {
	if !bson.IsObjectIdHex(id) {
		return fmt.Errorf("No AutoPie ID:%s", id)
	}
	session := mgoSession.Clone()
	defer session.Close()
	err := db.cPieAuto(session).RemoveId(bson.ObjectIdHex(id))
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
	session, closer := session(sessionIn)
	defer closer()
	autopie := new(PieAuto)
	col := db.cPieAuto(session)
	err := col.FindId(bson.ObjectIdHex(id)).One(autopie)
	if err != nil {
		log.Errorf("AutoPieByID Error:%s(%s)", err, id)
		return nil, err
	}
	return autopie, nil
}
