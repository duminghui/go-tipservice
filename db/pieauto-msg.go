// Package db provides ...
package db

import "github.com/globalsign/mgo"

type PieAutoMsg struct {
	MsgID     string `bson:"msgID,omitempty"`
	UserID    string `bson:"userID,omitempty"`
	PieAutoID string `bson:"pieautoID,omitempty"`
}

func (db *DBGuild) cPieAutoMsg(session *mgo.Session) *mgo.Collection {
	return session.DB(db.database).C("pieautomsg")
}

func (db *DBGuild) PieAutoMsgAdd(msgID, userID, pieAutoID string) error {
	session := mgoSession.Clone()
	defer session.Clone()
	data := &PieAutoMsg{
		MsgID:     msgID,
		UserID:    userID,
		PieAutoID: pieAutoID,
	}
	err := db.cPieAutoMsg(session).Insert(data)
	if err != nil {
		log.Errorf("PieAutoMsgAdd Error:%s[mID:%s,uID:%s,pieAutoID:%s]", err, msgID, userID, pieAutoID)
	}
	return err
}

func (db *DBGuild) PieAutoMsg(msgID, userID string) (*PieAutoMsg, error) {
	session := mgoSession.Clone()
	defer session.Clone()
	selector := &PieAutoMsg{
		MsgID:  msgID,
		UserID: userID,
	}
	pieAutoMsg := new(PieAutoMsg)
	err := db.cPieAutoMsg(session).Find(selector).One(pieAutoMsg)
	if err != nil {
		return nil, err
	}
	return pieAutoMsg, nil
}

func (db *DBGuild) PieAutoMsgRemove(msgID, pieAutoID string) error {
	session := mgoSession.Clone()
	defer session.Clone()
	selector := &PieAutoMsg{
		MsgID:     msgID,
		PieAutoID: pieAutoID,
	}
	err := db.cPieAutoMsg(session).Remove(selector)
	return err
}
