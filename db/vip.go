// Package db provides ...
package db

import (
	"github.com/globalsign/mgo"
)

func (db *DBSymbol) cVipUserPoints(s *mgo.Session) *mgo.Collection {
	return s.DB(db.database).C("vip_user_points")
}

type VipUserPoints struct {
	UserID   string `bson:"userid,omitempty"`
	Points   int64  `bson:"points,omitempty"`
	RoleName string `bosn:"rolename,omitempty"`
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
	err := db.cVipUserPoints(session).Update(selector, data)
	if err != nil {
		log.Infof("VipUserPointsRoleName Error:%s:%s:%s", err, userID, roleName)
		return err
	}
	return nil
}

func (db *DBSymbol) VipUserPointsChange(userID string, points int64) (*VipUserPoints, error) {
	session := mgoSession.Clone()
	defer session.Close()
	userPoints, err := db.VipUserPoints(session, userID)
	if points == 0 {
		return userPoints, err
	}
	finalPoints := points
	if userPoints != nil {
		finalPoints = userPoints.Points + points
	}
	if finalPoints < 0 {
		finalPoints = 0
	}
	selector := &VipUserPoints{
		UserID: userID,
	}
	data := &VipUserPoints{
		UserID: userID,
		Points: finalPoints,
	}
	_, err = db.cVipUserPoints(session).Upsert(selector, data)
	if err != nil {
		log.Infof("VipUserPointsChange Error:%s:%s:%d", err, userID, points)
		return nil, err
	}
	return data, nil
}

func (db *DBSymbol) VipUserPoints(s *mgo.Session, userID string) (*VipUserPoints, error) {
	session, colser := session(s)
	defer colser()
	selector := &VipUserPoints{
		UserID: userID,
	}
	userPoints := new(VipUserPoints)
	err := db.cVipUserPoints(session).Find(selector).One(userPoints)
	if err != nil {
		if err == mgo.ErrNotFound {
			return nil, nil
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
	ChannelID string `bson:"channel,omitempty"`
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
		log.Infof("VipChannelPointsList Error:%s", err)
		return nil, err
	}
	return channelPointsList, nil
}
