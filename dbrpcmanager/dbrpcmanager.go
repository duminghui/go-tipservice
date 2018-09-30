// Package member provides ...
package dbrpcmanager

import (
	rpcclient "github.com/duminghui/go-rpcclient"
	"github.com/duminghui/go-tipservice/config"
	"github.com/duminghui/go-tipservice/db"
	"github.com/globalsign/mgo"
	"github.com/sirupsen/logrus"
)

type DBRpcManager struct {
	DB  *db.DB
	RPC *rpcclient.Client
}

func (m *DBRpcManager) Init(log *logrus.Logger, session *mgo.Session, config *config.CoinInfo) {
	db.SetLog(log)
	db.SetSession(session)
	m.DB = db.New(config.Symbol, config.Database)
	rpcclient.SetLog(log)
	m.RPC = rpcclient.New(config.RPC)
}
