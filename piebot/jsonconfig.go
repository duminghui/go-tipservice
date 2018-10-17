// Package main provides ...
package main

import (
	"fmt"

	"github.com/duminghui/go-tipservice/config"
	"github.com/duminghui/go-util/ulog"
	"github.com/duminghui/go-util/umgo"
)

type Discord struct {
	Token                    string `json:"token"`
	SuperManagerIDs          string `json:"supermanagerids"`
	Prefix                   string `json:"prefix"`
	EachPieMsgReceiversLimit int    `json:"eachPieMsgReceiversLimit"`
}

type PieBotConfig struct {
	Discord       *Discord     `json:"discord"`
	WorkDir       string       `json:"workDir"`
	PidFile       string       `json:"pidFile"`
	Log           *ulog.Config `json:"log"`
	DBConfigFile  string       `json:"dbConfigFile"`
	CoinInfosFile string       `json:"coinInfosFile"`
}

var piebotConfig *PieBotConfig
var dbConfig *umgo.ConnConfig
var coinInfos map[string]*config.CoinInfo

func readConfig(file string) error {
	log.Infof("PieBot Config File:%s", file)
	piebotConfig = new(PieBotConfig)
	_, err := config.FromFile(file, piebotConfig)
	if err != nil {
		return fmt.Errorf("PieBotConfig:%s", err)
	}

	dbConfigFile := piebotConfig.DBConfigFile
	log.Infof("DB Config File:%s", dbConfigFile)
	dbConfig = new(umgo.ConnConfig)
	_, err = config.FromFile(dbConfigFile, dbConfig)
	if err != nil {
		return fmt.Errorf("DBConfig:%s[%s]", err, dbConfigFile)
	}

	coinInfosFile := piebotConfig.CoinInfosFile
	log.Infof("Coin Infos Config File:%s", coinInfosFile)
	_, err = config.FromFile(coinInfosFile, &coinInfos)
	if err != nil {
		return fmt.Errorf("CoinInfosConfig:%s[%s]", err, coinInfosFile)
	}
	return nil
}

func hasSymbol(symbol string) bool {
	_, ok := coinInfos[symbol]
	if ok {
		return true
	}
	return false
}
