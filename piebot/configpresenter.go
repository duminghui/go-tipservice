// Package main provides ...
package main

import (
	"fmt"

	"github.com/duminghui/go-tipservice/db"
)

func hasSymbol(symbol symbolWrap) bool {
	_, ok := coinPresenters[symbol]
	if ok {
		return true
	}
	return false
}

type symbolWrap string
type prefixWrap string

type guildConfigPresenter struct {
	prefixSymbolMap map[prefixWrap]symbolWrap
	symbolPrefixMap map[symbolWrap]prefixWrap
	// key symbol
	guildCoinConfig map[symbolWrap]*db.GuildCoinConfig
	guildManager    *db.GuildManager
}

// key guildid
type guildConfigPresenterMap map[string]*guildConfigPresenter

var guildConfigPresenters = make(guildConfigPresenterMap)

func initGuildConfig() {
	guildConfigList, err := db.GuildCoinConfigList()
	if err != nil {
		return
	}
	for _, v := range guildConfigList {
		guildCfgMge, ok := guildConfigPresenters[v.GuildID]
		if !ok {
			guildCfgMge = guildConfigPresenters.initGuildConfigManager(v.GuildID)
		}
		guildCfgMge.prefixSymbolMap[prefixWrap(v.CmdPrefix)] = symbolWrap(v.Symbol)
		guildCfgMge.symbolPrefixMap[symbolWrap(v.Symbol)] = prefixWrap(v.CmdPrefix)
		guildCfgMge.guildCoinConfig[symbolWrap(v.Symbol)] = v
	}
	guildManagerList, err := db.GuildManagerList()
	if err != nil {
		return
	}
	for _, v := range guildManagerList {
		guildCfgMge, ok := guildConfigPresenters[v.GuildID]
		if !ok {
			guildCfgMge = guildConfigPresenters.initGuildConfigManager(v.GuildID)
		}
		guildCfgMge.guildManager = v
	}
}

func (p guildConfigPresenterMap) initGuildConfigManager(guildID string) *guildConfigPresenter {
	guildCfgMge := new(guildConfigPresenter)
	guildCfgMge.prefixSymbolMap = make(map[prefixWrap]symbolWrap)
	guildCfgMge.symbolPrefixMap = make(map[symbolWrap]prefixWrap)
	guildCfgMge.guildCoinConfig = make(map[symbolWrap]*db.GuildCoinConfig)
	p[guildID] = guildCfgMge
	return guildCfgMge
}

func (p guildConfigPresenterMap) prefixList(guildID string) []prefixWrap {
	guildCfgMge, ok := guildConfigPresenters[guildID]
	if !ok {
		return []prefixWrap{}
	}
	prefixs := make([]prefixWrap, 0, len(guildCfgMge.prefixSymbolMap))
	for k := range guildCfgMge.prefixSymbolMap {
		prefixs = append(prefixs, k)
	}
	return prefixs
}

func (p guildConfigPresenterMap) symbolByPrefix(guildID string, pfx prefixWrap) (symbolWrap, error) {
	guildCfgMge, ok := p[guildID]
	if !ok {
		return "", fmt.Errorf("symbolByPrefix no GuildConfigMapager:%s:%s", guildID, pfx)
	}
	sbl, ok := guildCfgMge.prefixSymbolMap[pfx]
	if !ok {
		return "", fmt.Errorf("symbolByPrefix No Symbol:%s %s", guildID, pfx)
	}
	return sbl, nil
}

func (p guildConfigPresenterMap) updatePrefix(guildID string, sbl symbolWrap, newPrefix prefixWrap) error {
	guildCfgMge, ok := p[guildID]
	if !ok {
		guildCfgMge = p.initGuildConfigManager(guildID)
	}
	guildCoinConfig, ok := guildCfgMge.guildCoinConfig[sbl]
	if !ok {
		guildCoinConfig = &db.GuildCoinConfig{
			GuildID: guildID,
			Symbol:  string(sbl),
		}
	}
	guildCoinConfig.UpdateCmdPrefix(string(newPrefix))
	oldPfx, ok := guildCfgMge.symbolPrefixMap[sbl]
	if ok {
		delete(guildCfgMge.prefixSymbolMap, oldPfx)
	}
	guildCfgMge.prefixSymbolMap[newPrefix] = sbl
	guildCfgMge.symbolPrefixMap[sbl] = newPrefix
	guildCoinConfigDB, err := db.GuildCoinConfigBySymbol(guildID, string(sbl))
	if err != nil {
		log.Errorf("update prefix GuildCoinConfigBySymbol error:%s,%s:%s", err, guildID, sbl)
		return err
	}
	guildCfgMge.guildCoinConfig[sbl] = guildCoinConfigDB
	return nil
}

func (p guildConfigPresenterMap) guildCoinConfigBySymbol(guildID string, symbol symbolWrap) (*db.GuildCoinConfig, error) {
	guildCfgMge, ok := p[guildID]
	if !ok {
		return nil, fmt.Errorf("coinConfig no GuildConfigManager:%s:%s", guildID, symbol)
	}
	coinConfig, ok := guildCfgMge.guildCoinConfig[symbol]
	if !ok {
		coinConfig = new(db.GuildCoinConfig)
		coinConfig.GuildID = guildID
		coinConfig.Symbol = string(symbol)
		guildCfgMge.guildCoinConfig[symbol] = coinConfig
	}
	return coinConfig, nil
}

func (p guildConfigPresenterMap) guildChannelUpdate(guildID string, symbol symbolWrap, operator string, channel []string) ([]string, error) {
	guildCfgMge := p[guildID]
	guildCoinConfig := guildCfgMge.guildCoinConfig[symbol]
	var err error
	if operator == "add" {
		err = guildCoinConfig.ChannelAdd(channel)
	} else {
		err = guildCoinConfig.ChannelRemove(channel)
	}
	if err != nil {
		log.Errorf("GuildChannelUpdate Error:%s,%s:%s:%s", err, guildID, symbol, operator)
		return nil, err
	}
	guildCoinConfigDB, err := db.GuildCoinConfigBySymbol(guildID, string(symbol))
	if err != nil {
		log.Errorf("update channel GuildCoinConfigBySymbol error:%s,%s:%s", err, guildID, symbol)
		return nil, err
	}
	guildCfgMge.guildCoinConfig[symbol] = guildCoinConfigDB
	return guildCoinConfigDB.ChannelIDs, nil
}

func (p guildConfigPresenterMap) guildManagerUpdate(guildID string, operator string, users, roles []string) ([]string, []string, error) {
	guildCfgMge, ok := guildConfigPresenters[guildID]
	if !ok {
		guildCfgMge = guildConfigPresenters.initGuildConfigManager(guildID)
	}
	gm := guildCfgMge.guildManager
	if gm == nil {
		gm = &db.GuildManager{
			GuildID: guildID,
		}
	}
	var err error
	if operator == "add" {
		err = gm.ManagerAdd(users, roles)
	} else {
		err = gm.ManagerRemove(users, roles)
	}
	if err != nil {
		log.Errorf("guildManagerUpdate Error:%s,%s:%s", err, guildID, operator)
		return nil, nil, err
	}
	gmDB, err := db.GuildManagerByGuildID(guildID)
	if err != nil {
		log.Errorf("guildManagerUpdate read from db Error:%s,%s", err, guildID)
		return nil, nil, err
	}
	guildCfgMge.guildManager = gmDB
	return gmDB.Managers, gmDB.ManagerRoles, nil
}
