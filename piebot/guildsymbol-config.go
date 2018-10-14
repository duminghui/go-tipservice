// Package main provides ...
package main

import (
	"sync"

	"github.com/duminghui/go-tipservice/db"
)

type guildCoinConfig db.GuildCoinConfig

func (p *guildCoinConfig) inChannels(channelID string) bool {
	if len(p.ChannelIDs) == 0 {
		return true
	}
	for _, channel := range p.ChannelIDs {
		if channel == channelID {
			return true
		}
	}
	return false
}

type symbolCoinConfigMap map[symbol]*guildCoinConfig
type guildSymbolCoinConfigMap map[string]symbolCoinConfigMap

var guildSymbolCoinConfigs = make(guildSymbolCoinConfigMap)

func readGuildSymbolCoinConfigsFromDB() {
	guildConfigList, err := dbGuild.GuildCoinConfigList()
	if err != nil {
		return
	}
	for _, v := range guildConfigList {
		glv1, ok := guildSymbolCoinConfigs[v.GuildID]
		if !ok {
			glv1 = guildSymbolCoinConfigs.initSymbolCoinConfigMap(v.GuildID)
		}
		sbl := symbol(v.Symbol)
		glv1[sbl] = (*guildCoinConfig)(v)
	}
}

func (gccm guildSymbolCoinConfigMap) initSymbolCoinConfigMap(guildID string) symbolCoinConfigMap {
	sccm := make(symbolCoinConfigMap)
	for k := range coinInfos {
		gcc := new(guildCoinConfig)
		gcc.ChannelIDs = []string{}
		sccm[symbol(k)] = gcc
	}
	gccm[guildID] = sccm
	return sccm
}

var mu sync.RWMutex

func (gccm guildSymbolCoinConfigMap) sccm(guildID string) symbolCoinConfigMap {
	mu.RLock()
	sccm, ok := gccm[guildID]
	if ok {
		mu.RUnlock()
		return sccm
	}
	mu.RUnlock()
	mu.Lock()
	sccm = gccm.initSymbolCoinConfigMap(guildID)
	mu.Unlock()
	return sccm
}

func (gccm guildSymbolCoinConfigMap) update(guildID string, sbl symbol, gcc *db.GuildCoinConfig) {
	mu.Lock()
	defer mu.Unlock()
	gccm[guildID][sbl] = (*guildCoinConfig)(gcc)
}

func (sccm symbolCoinConfigMap) symbolByPrefix(pfx prefix) (symbol, bool) {
	mu.RLock()
	defer mu.RUnlock()
	for k, v := range sccm {
		if v.CmdPrefix == string(pfx) {
			return k, true
		}
	}
	return "", false
}
