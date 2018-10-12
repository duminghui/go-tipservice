// Package db provides ...
package db

import (
	"github.com/globalsign/mgo"
	"github.com/sirupsen/logrus"
)

var log *logrus.Logger
var mgoSession *mgo.Session

func SetLog(logUse *logrus.Logger) {
	log = logUse
}

func SetSession(session *mgo.Session) {
	mgoSession = session
}

type DBSymbol struct {
	symbol   string
	database string
}

func NewDBSymbol(symbol, database string) *DBSymbol {
	return &DBSymbol{
		symbol:   symbol,
		database: database,
	}
}

type DBGuild struct {
	database string
}

func NewDBGuild() *DBGuild {
	return &DBGuild{
		database: "guild_config",
	}
}

func session(sessionUse *mgo.Session) (*mgo.Session, func()) {
	if sessionUse == nil {
		session := mgoSession.Clone()
		return session, func() { session.Close() }
	}
	return sessionUse, func() {}
}
