// Package main provides ...
package main

import (
	"fmt"

	"github.com/duminghui/go-tipservice/config"
	"github.com/duminghui/go-util/ulog"
	"github.com/duminghui/go-util/umgo"
)

type Discord struct {
	Token          string `json:"token"`
	SuperManagerID string `json:"supermanagerid"`
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
	piebotConfig = new(PieBotConfig)
	_, err := config.FromFile(file, piebotConfig)
	if err != nil {
		return fmt.Errorf("PieBotConfig:%s", err)
	}
	dbConfig = new(umgo.ConnConfig)
	_, err = config.FromFile(piebotConfig.DBConfigFile, dbConfig)
	if err != nil {
		return fmt.Errorf("DBConfig:%s", err)
	}
	_, err = config.FromFile(piebotConfig.CoinInfosFile, &coinInfos)
	if err != nil {
		return fmt.Errorf("CoinInfosConfig:%s", err)
	}
	return nil
}
