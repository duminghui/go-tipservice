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
	serverConfig = new(ServerConfig)
	_, err := config.FromFile(file, serverConfig)
	if err != nil {
		return fmt.Errorf("ServerConfig:%s[%s]", err, file)
	}
	dbConfig = new(umgo.ConnConfig)
	_, err = config.FromFile(serverConfig.DBConfigFile, dbConfig)
	if err != nil {
		return fmt.Errorf("DBConfig:%s[%s]", err, serverConfig.DBConfigFile)
	}
	_, err = config.FromFile(serverConfig.CoinInfosFile, &coinInfos)
	if err != nil {
		return fmt.Errorf("CoinInfosConfig:%s[%s]", err, serverConfig.CoinInfosFile)
	}
	return nil
}
