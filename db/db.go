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

type DB struct {
	symbol   string
	database string
}

func New(symbol, database string) *DB {
	return &DB{
		symbol:   symbol,
		database: database,
	}
}

func session(sessionUse *mgo.Session) (*mgo.Session, func()) {
	if sessionUse == nil {
		session := mgoSession.Clone()
		return session, func() { session.Close() }
	}
	return sessionUse, func() {}
}
