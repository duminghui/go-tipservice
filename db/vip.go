// Package db provides ...
package db

import (
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
)

func (db *DBSymbol) cVipUserPoints(s *mgo.Session) *mgo.Collection {
	return s.DB(db.database).C("vip_user_points")
}

type VipUserPoints struct {
	UserID   string `bson:"userid,omitempty"`
	Points   int64  `bson:"points,omitempty"`
	RoleName string `bson:"rolename,omitempty"`
}

type VipUserPointsPoints struct {
	Points int64 `bson:"points"`
}

func (db *DBSymbol) VipUserPointsCount() (int, error) {
	session := mgoSession.Clone()
	defer session.Close()
	count, err := db.cVipUserPoints(session).Count()
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (db *DBSymbol) VipUserPointsList(start, size int) ([]*VipUserPoints, error) {
	session := mgoSession.Clone()
	defer session.Close()
	list := make([]*VipUserPoints, 0)
	err := db.cVipUserPoints(session).Find(nil).Sort("-points").Skip(start).Limit(size).All(&list)
	return list, err
}

func (db *DBSymbol) VipUserPointsRoleName(userID, roleName string) error {
	session := mgoSession.Clone()
	defer session.Close()
	selector := &VipUserPoints{
		UserID: userID,
	}
	data := &VipUserPoints{
		RoleName: roleName,
	}
	err := db.cVipUserPoints(session).Update(selector, bson.M{"$set": data})
	if err != nil {
		log.Errorf("VipUserPointsRoleName Error:%s:%s:%s", err, userID, roleName)
		return err
	}
	return nil
}

func (db *DBSymbol) VipUserPointsChange(userID string, points int64) (*VipUserPoints, error) {
	session := mgoSession.Clone()
	defer session.Close()
	userPoints, err := db.VipUserPointsByUserID(session, userID)
	col := db.cVipUserPoints(session)
	if err != nil {
		if err == mgo.ErrNotFound {
			if points < 0 {
				points = 0
			}
			if points == 0 {
				return userPoints, nil
			}
			data := &VipUserPoints{
				UserID:   userID,
				Points:   points,
				RoleName: "Not VIP",
			}
			err := col.Insert(data)
			if err != nil {
				return nil, err
			}
			return data, nil
		}
		return nil, err
	}
	if points == 0 {
		return userPoints, nil
	}
	finalPoints := userPoints.Points + points
	if finalPoints < 0 {
		finalPoints = 0
	}
	userPoints.Points = finalPoints
	selector := &VipUserPoints{
		UserID: userID,
	}
	data := &VipUserPointsPoints{
		Points: finalPoints,
	}
	err = db.cVipUserPoints(session).Update(selector, bson.M{"$set": data})
	if err != nil {
		log.Infof("VipUserPointsChange Error:%s:%s:%d", err, userID, points)
		return nil, err
	}
	return userPoints, nil
}

func (db *DBSymbol) VipUserPointsByUserID(s *mgo.Session, userID string) (*VipUserPoints, error) {
	session, colser := session(s)
	defer colser()
	selector := &VipUserPoints{
		UserID: userID,
	}
	userPoints := new(VipUserPoints)
	err := db.cVipUserPoints(session).Find(selector).One(&userPoints)
	if err != nil {
		if err == mgo.ErrNotFound {
			return &VipUserPoints{
				UserID:   userID,
				Points:   0,
				RoleName: "Not VIP",
			}, err
		}
		log.Infof("VipUserPoints Error:%s:%s", err, userID)
		return nil, err
	}
	return userPoints, nil
}

func (db *DBSymbol) cVipRolePoints(s *mgo.Session) *mgo.Collection {
	return s.DB(db.database).C("vip_role_points")
}

type VipRolePoints struct {
	RoleID string `bson:"roleid,omitempty"`
	Points int64  `bson:"points,omitempty"`
}

func (db *DBSymbol) VipRolePointsSet(roleID string, points int64) error {
	session := mgoSession.Clone()
	defer session.Close()
	selector := &VipRolePoints{
		RoleID: roleID,
	}
	col := db.cVipRolePoints(session)
	if points == 0 {
		err := col.Remove(selector)
		if err != nil && err != mgo.ErrNotFound {
			log.Infof("VipRolePointsSet Remove Error:%s:%s:%d", err, roleID, points)
			return err
		}
		return nil
	}
	data := &VipRolePoints{
		RoleID: roleID,
		Points: points,
	}
	_, err := col.Upsert(selector, data)
	if err != nil {
		log.Infof("VipRolePointsSet Error:%s:%s:%d", err, roleID, points)
		return err
	}
	return nil
}

func (db *DBSymbol) VipRolePointsList() ([]*VipRolePoints, error) {
	session := mgoSession.Clone()
	defer session.Close()
	rolePointsList := make([]*VipRolePoints, 0)
	err := db.cVipRolePoints(session).Find(nil).Sort("points").All(&rolePointsList)
	if err != nil {
		log.Infof("VipRolePointsList Error:%s", err)
		return nil, err
	}
	return rolePointsList, nil
}

func (db *DBSymbol) cVipChannelPoints(s *mgo.Session) *mgo.Collection {
	return s.DB(db.database).C("vip_channel_points")
}

type VipChannelPoints struct {
	ChannelID string `bson:"channelid,omitempty"`
	Points    int64  `bson:"points,omitempty"`
}

func (db *DBSymbol) VipChannelPointsSet(channelID string, points int64) error {
	session := mgoSession.Clone()
	defer session.Close()
	selector := &VipChannelPoints{
		ChannelID: channelID,
	}
	col := db.cVipChannelPoints(session)
	if points == 0 {
		err := col.Remove(selector)
		if err != nil && err != mgo.ErrNotFound {
			log.Infof("VipChannelPointsSet Remove Error:%s:%s:%d", err, channelID, points)
			return err
		}
		return nil
	}
	data := &VipChannelPoints{
		ChannelID: channelID,
		Points:    points,
	}
	_, err := col.Upsert(selector, data)
	if err != nil {
		log.Infof("VipChannelPointsSet Error:%s:%s:%d", err, channelID, points)
		return err
	}
	return nil
}

func (db *DBSymbol) VipChannelPointsList() ([]*VipChannelPoints, error) {
	session := mgoSession.Clone()
	defer session.Close()
	channelPointsList := make([]*VipChannelPoints, 0)
	err := db.cVipChannelPoints(session).Find(nil).All(&channelPointsList)
	if err != nil {
		log.Errorf("VipChannelPointsList Error:%s", err)
		return nil, err
	}
	return channelPointsList, nil
}

func (db *DBSymbol) VipChannelPointsByChannelID(channelID string) (*VipChannelPoints, error) {
	session := mgoSession.Clone()
	defer session.Close()
	channelPoints := new(VipChannelPoints)
	selector := &VipChannelPoints{
		ChannelID: channelID,
	}
	err := db.cVipChannelPoints(session).Find(selector).One(channelPoints)
	if err != nil {
		if err != mgo.ErrNotFound {
			log.Errorf("VipChannelPointsByChannelID Error:%s:%s", err, channelID)
		}
		return nil, err
	}
	return channelPoints, nil
}

func (db *DBSymbol) cVipEmoji(s *mgo.Session) *mgo.Collection {
	return s.DB(db.database).C("vip_emoji")
}

type VipEmoji struct {
	Key  string `bson:"key,omitempty"`
	ID   string `bson:"id,omitempty"`
	Name string `bson:"name,omitempty"`
}

func (db *DBSymbol) VipEmojiChange(id, name string) (*VipEmoji, error) {
	session := mgoSession.Clone()
	defer session.Close()
	selector := &VipEmoji{
		Key: "1",
	}
	data := &VipEmoji{
		Key:  "1",
		ID:   id,
		Name: name,
	}
	_, err := db.cVipEmoji(session).Upsert(selector, data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (db *DBSymbol) VipEmoji() (*VipEmoji, error) {
	session := mgoSession.Clone()
	defer session.Close()
	vipEmoji := new(VipEmoji)
	err := db.cVipEmoji(session).Find(nil).One(vipEmoji)
	if err != nil {
		return nil, err
	}
	return vipEmoji, nil
}
