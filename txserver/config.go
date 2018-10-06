// Package main provides ...
package main

import (
	"fmt"

	"github.com/duminghui/go-tipservice/config"
	"github.com/duminghui/go-util/ulog"
	"github.com/duminghui/go-util/umgo"
)

type ServerConfig struct {
	ListenerAddr  string       `json:"listenerAddr"`
	WorkDir       string       `json:"workDir"`
	PidFile       string       `json:"pidFile"`
	Log           *ulog.Config `json:"log"`
	DBConfigFile  string       `json:"dbConfigFile"`
	CoinInfosFile string       `json:"coinInfosFile"`
}

var serverConfig *ServerConfig
var dbConfig *umgo.ConnConfig
var coinInfos map[string]*config.CoinInfo

func readConfig(file string) error {
	log.Infof("Server Config File:%s", file)
	serverConfig = new(ServerConfig)
	_, err := config.FromFile(file, serverConfig)
	if err != nil {
		return fmt.Errorf("ServerConfig:%s[%s]", err, file)
	}

	dbConfigFile := serverConfig.DBConfigFile
	log.Infof("DB Config File:%s", dbConfigFile)
	dbConfig = new(umgo.ConnConfig)
	_, err = config.FromFile(dbConfigFile, dbConfig)
	if err != nil {
		return fmt.Errorf("DBConfig:%s[%s]", err, dbConfigFile)
	}

	coinInfosFile := serverConfig.CoinInfosFile
	log.Infof("Coin Infos Config File:%s", dbConfigFile)
	_, err = config.FromFile(coinInfosFile, &coinInfos)
	if err != nil {
		return fmt.Errorf("CoinInfosConfig:%s[%s]", err, coinInfosFile)
	}
	return nil
}
