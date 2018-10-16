// Package main provides ...
package main

import (
	"fmt"

	"github.com/duminghui/go-tipservice/config"
	"github.com/duminghui/go-util/ulog"
	"github.com/duminghui/go-util/umgo"
)

type Discord struct {
	Token           string `json:"token"`
	SuperManagerIDs string `json:"supermanagerids"`
	Prefix          string `json:"prefix"`
	VipFlagEmojiID  string `json:"vipflagemojiid"`
}

type BcrmVipConfig struct {
	Discord      *Discord     `json:"discord"`
	WorkDir      string       `json:"workDir"`
	PidFile      string       `json:"pidFile"`
	Log          *ulog.Config `json:"log"`
	DBConfigFile string       `json:"dbConfigFile"`
}

var bcrmVipConfig *BcrmVipConfig
var dbConfig *umgo.ConnConfig

func readConfig(file string) error {
	log.Infof("BcrmVipBot Config File:%s", file)
	bcrmVipConfig = new(BcrmVipConfig)
	_, err := config.FromFile(file, bcrmVipConfig)
	if err != nil {
		return fmt.Errorf("BcrmVipBotConfig:%s", err)
	}

	dbConfigFile := bcrmVipConfig.DBConfigFile
	log.Infof("DB Config File:%s", dbConfigFile)
	dbConfig = new(umgo.ConnConfig)
	_, err = config.FromFile(dbConfigFile, dbConfig)
	if err != nil {
		return fmt.Errorf("DBConfig:%s[%s]", err, dbConfigFile)
	}

	return nil
}
