// Package main provides ...
package main

import (
	"github.com/duminghui/go-tipservice/config"
	"github.com/duminghui/go-tipservice/db"
)

type guildSymbolPresenter struct {
	symbol   string
	guildID  string
	dbGuild  *db.DBGuild
	dbSymbol *db.DBSymbol
	coinInfo *config.CoinInfo
}

var dbBcrm = db.NewDBSymbol("BCRM", "bcrm")
var dbGuild = db.NewDBGuild()

var presenter *guildSymbolPresenter

func initPresenter() {
	presenter = &guildSymbolPresenter{
		symbol:   "BCRM",
		dbGuild:  dbGuild,
		dbSymbol: dbBcrm,
		coinInfo: coinInfo,
	}
}
