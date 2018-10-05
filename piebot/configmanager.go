// Package main provides ...
package main

import (
	"errors"
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

type guildConfigManager struct {
	prefixSymbolMap map[prefixWrap]symbolWrap
	symbolPrefixMap map[symbolWrap]prefixWrap
	// key symbol
	guildCoinConfig map[symbolWrap]*db.GuildCoinConfig
	guildManager    *db.GuildManager
}

type guildConfigManagerMap map[string]*guildConfigManager

var guildConfigManagers = make(guildConfigManagerMap)

func initGuildConfig() {
	guildConfigList, err := db.GuildConfigList()
	if err != nil {
		return
	}
	for _, v := range guildConfigList {
		guildCfgMge, ok := guildConfigManagers[v.GuildID]
		if !ok {
			guildCfgMge = guildConfigManagers.initGuildConfigManager(v.GuildID)
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
		guildCfgMge, ok := guildConfigManagers[v.GuildID]
		if !ok {
			guildCfgMge = guildConfigManagers.initGuildConfigManager(v.GuildID)
		}
		guildCfgMge.guildManager = v
	}

}

func (cfg guildConfigManagerMap) initGuildConfigManager(guildID string) *guildConfigManager {
	guildCfgMge := new(guildConfigManager)
	guildCfgMge.prefixSymbolMap = make(map[prefixWrap]symbolWrap)
	guildCfgMge.symbolPrefixMap = make(map[symbolWrap]prefixWrap)
	guildCfgMge.guildCoinConfig = make(map[symbolWrap]*db.GuildCoinConfig)
	cfg[guildID] = guildCfgMge
	return guildCfgMge
}

func (cfg guildConfigManagerMap) prefixList(guildID string) []prefixWrap {
	guildCfgMge, ok := guildConfigManagers[guildID]
	if !ok {
		return []prefixWrap{}
	}
	prefixs := make([]prefixWrap, 0, len(guildCfgMge.prefixSymbolMap))
	for k := range guildCfgMge.prefixSymbolMap {
		prefixs = append(prefixs, k)
	}
	return prefixs
}

func (cfg guildConfigManagerMap) coinConfig(guildID string, symbol symbolWrap) (*db.GuildCoinConfig, error) {
	guildCfgMge, ok := cfg[guildID]
	if !ok {
		return nil, errors.New("coinConfig no GuildConfigManager")
	}
	coinConfig, ok := guildCfgMge.guildCoinConfig[symbol]
	if !ok {
		coinConfig = new(db.GuildCoinConfig)
		guildCfgMge.guildCoinConfig[symbol] = coinConfig
	}
	return coinConfig, nil
}

func (cfg guildConfigManagerMap) symbolByPrefix(guildID string, pfx prefixWrap) (symbolWrap, error) {
	guildCfgMge, ok := cfg[guildID]
	if !ok {
		return "", errors.New("symbolByPrefix no GuildConfigManager")
	}
	sbl, ok := guildCfgMge.prefixSymbolMap[pfx]
	if !ok {
		errMsg := fmt.Sprintf("symbolByPrefix No Symbol:%s %s", guildID, pfx)
		return "", errors.New(errMsg)
	}
	return sbl, nil
}

func (cfg guildConfigManagerMap) updatePrefix(guildID string, sbl symbolWrap, newPrefix prefixWrap) error {
	guildCfgMge, ok := cfg[guildID]
	if !ok {
		guildCfgMge = cfg.initGuildConfigManager(guildID)
		guildCoinConfig := new(db.GuildCoinConfig)
		guildCfgMge.guildCoinConfig[sbl] = guildCoinConfig

	}
	oldPfx, ok := guildCfgMge.symbolPrefixMap[sbl]
	if ok {
		delete(guildCfgMge.prefixSymbolMap, oldPfx)
	}
	guildCfgMge.prefixSymbolMap[newPrefix] = sbl
	guildCfgMge.symbolPrefixMap[sbl] = newPrefix
	return nil
}

func (cfg guildConfigManagerMap) guildConfigBySymbol(guildID string, sbl symbolWrap) (*db.GuildCoinConfig, error) {
	guildCfgMge, ok := cfg[guildID]
	if !ok {
		return nil, errors.New("symbolByPrefix no GuildConfigManager")
	}
	guildCfg, ok := guildCfgMge.guildCoinConfig[sbl]
	if !ok {
		errMsg := fmt.Sprintf("getGuildConfigBySymbol No GuildConfig:%s %s", guildID, sbl)
		return nil, errors.New(errMsg)
	}
	return guildCfg, nil
}

func (cfg guildConfigManagerMap) guildConfigByPrefix(guildID string, prefix prefixWrap) (*db.GuildCoinConfig, error) {
	guildCfgMge, ok := cfg[guildID]
	if !ok {
		return nil, errors.New("symbolByPrefix no GuildConfigManager")
	}
	sbl, ok := guildCfgMge.prefixSymbolMap[prefix]
	if !ok {
		errMsg := fmt.Sprintf("getGuildConfigByPrefix No Symbol:%s %s", guildID, prefix)
		return nil, errors.New(errMsg)
	}
	guildCfg, ok := guildCfgMge.guildCoinConfig[sbl]
	if !ok {
		errMsg := fmt.Sprintf("getGuildConfigByPrefix No GuildConfig:%s %s", guildID, prefix)
		return nil, errors.New(errMsg)
	}
	return guildCfg, nil
}

func (gcm guildConfigManagerMap) guildChannelUpdate(guildID string, symbol symbolWrap, channel []string) {
	guildCfgMge, _ := gcm[guildID]
	guildCfgMge.guildCoinConfig[symbol].ChannelIDs = channel
}

func (gcm guildConfigManagerMap) guildManagerUpdate(guildID string, users, roles []string) {
	guildCfgMge, ok := guildConfigManagers[guildID]
	if !ok {
		guildCfgMge = guildConfigManagers.initGuildConfigManager(guildID)
	}
	if guildCfgMge.guildManager == nil {
		guildCfgMge.guildManager = &db.GuildManager{
			GuildID:      guildID,
			Managers:     users,
			ManagerRoles: roles,
		}
	} else {
		guildCfgMge.guildManager.Managers = users
		guildCfgMge.guildManager.ManagerRoles = roles
	}
}
