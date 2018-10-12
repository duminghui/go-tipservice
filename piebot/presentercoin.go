// Package main provides ...
package main

import (
	rpcclient "github.com/duminghui/go-rpcclient"
	"github.com/duminghui/go-tipservice/config"
	"github.com/duminghui/go-tipservice/db"
)

type coinPresenter struct {
	dbSymbol *db.DBSymbol
	rpc      *rpcclient.Client
	coinInfo *config.CoinInfo
}

var coinPresenters = make(map[symbol]*coinPresenter)

func initCoinPresenters() {
	for k, v := range coinInfos {
		sbl := symbol(k)
		p := new(coinPresenter)
		p.dbSymbol = db.NewDBSymbol(v.Symbol, v.Database)
		p.rpc = rpcclient.New(v.RPC)
		p.coinInfo = v
		coinPresenters[sbl] = p
	}
}
