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

type GuildCoinConfig struct {
	prefix     prefixWrap
	symbol     symbolWrap
	channelIDs []string
}

func (g *GuildCoinConfig) inChannels(channelID string) bool {
	if len(g.channelIDs) == 0 {
		return true
	}
	for _, channel := range g.channelIDs {
		if channel == channelID {
			return true
		}
	}
	return false
}

type symbolWrap string
type prefixWrap string

type guildConfigPresenter struct {
	guildID      string
	guildName    string
	prefixSymbol map[prefixWrap]symbolWrap
	gccMap       map[symbolWrap]*GuildCoinConfig
	managers     []string
	managerRoles []string
	excludeRoles []string
	// symbolPrefixMap map[symbolWrap]prefixWrap
	// key symbol
	// guildCoinConfig map[symbolWrap]GuildCoinConfig
}

func (p *guildConfigPresenter) isManager(userID string) bool {
	for _, member := range p.managers {
		if member == userID {
			return true
		}
	}
	return false
}

func (p *guildConfigPresenter) inManagerRoles(userRoles []string) bool {
	for _, role := range p.managerRoles {
		for _, userRole := range userRoles {
			if role == userRole {
				return true
			}
		}
	}
	return false
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
		p, ok := guildConfigPresenters[v.GuildID]
		if !ok {
			p = guildConfigPresenters.initGuildConfigPresenter(v.GuildID)
		}
		gcc := &GuildCoinConfig{
			prefix:     prefixWrap(v.CmdPrefix),
			symbol:     symbolWrap(v.Symbol),
			channelIDs: v.ChannelIDs,
		}
		p.prefixSymbol[gcc.prefix] = gcc.symbol
		p.gccMap[symbolWrap(v.Symbol)] = gcc
	}
	guildManagerList, err := db.GuildConfigManagerList()
	if err != nil {
		return
	}
	for _, v := range guildManagerList {
		p, ok := guildConfigPresenters[v.GuildID]
		if !ok {
			p = guildConfigPresenters.initGuildConfigPresenter(v.GuildID)
		}
		p.managers = v.Managers
		p.managerRoles = v.ManagerRoles
	}
}

func (g guildConfigPresenterMap) initGuildConfigPresenter(guildID string) *guildConfigPresenter {
	p := new(guildConfigPresenter)
	p.guildID = guildID
	p.prefixSymbol = make(map[prefixWrap]symbolWrap)
	p.gccMap = make(map[symbolWrap]*GuildCoinConfig)
	p.managers = []string{}
	p.managerRoles = []string{}
	g[guildID] = p
	return p
}

func (p *guildConfigPresenter) prefixList() []prefixWrap {
	prefixs := make([]prefixWrap, 0, len(p.prefixSymbol))
	for k := range p.prefixSymbol {
		prefixs = append(prefixs, k)
	}
	return prefixs
}

func (p *guildConfigPresenter) symbolByPrefix(pfx prefixWrap) (symbolWrap, error) {
	symbol, ok := p.prefixSymbol[pfx]
	if !ok {
		return "", fmt.Errorf("symbolByPrefix No Symbol:%s %s", p.guildID, pfx)
	}
	return symbol, nil
}

func (p *guildConfigPresenter) updatePrefix(sbl symbolWrap, oldPrefix, newPrefix prefixWrap) error {
	err := db.GuildCoinUpdateCmdPrefix(p.guildID, p.guildName, string(sbl), string(newPrefix))
	if err != nil {
		log.Errorf("updatePrefix err:%s,%s,%s,%s", err, p.guildID, sbl, newPrefix)
		return err
	}
	delete(p.prefixSymbol, oldPrefix)
	p.prefixSymbol[newPrefix] = sbl
	gccDB, err := db.GuildCoinConfigBySymbol(nil, p.guildID, string(sbl))
	if err != nil {
		log.Errorf("update prefix GuildCoinConfigBySymbol error:%s,%s:%s", err, p.guildID, sbl)
		return err
	}
	gcc := &GuildCoinConfig{
		prefix:     prefixWrap(gccDB.CmdPrefix),
		symbol:     symbolWrap(gccDB.Symbol),
		channelIDs: gccDB.ChannelIDs,
	}
	p.gccMap[sbl] = gcc
	return nil
}

func (p *guildConfigPresenter) guildChannelUpdate(sbl symbolWrap, operator string, channels []string) ([]string, error) {
	var err error
	if operator == "add" {
		err = db.GuildCoinChannelAdd(p.guildID, p.guildName, string(sbl), channels)
	} else {
		err = db.GuildCoinChannelRemove(p.guildID, p.guildName, string(sbl), channels)
	}
	if err != nil {
		log.Errorf("GuildChannelUpdate Error:%s,%s:%s:%s", err, p.guildID, sbl, operator)
		return nil, err
	}
	gccDB, err := db.GuildCoinConfigBySymbol(nil, p.guildID, string(sbl))
	if err != nil {
		log.Errorf("update channel GuildCoinConfigBySymbol error:%s,%s:%s", err, p.guildID, sbl)
		return nil, err
	}
	gcc := &GuildCoinConfig{
		prefix:     prefixWrap(gccDB.CmdPrefix),
		symbol:     symbolWrap(gccDB.Symbol),
		channelIDs: gccDB.ChannelIDs,
	}
	p.gccMap[sbl] = gcc
	return gccDB.ChannelIDs, nil
}

func (p *guildConfigPresenter) guildManagerUpdate(operator string, users, roles []string) ([]string, []string, error) {
	var err error
	if operator == "add" {
		err = db.GuildConfigManagerAdd(p.guildID, p.guildName, users, roles)
	} else {
		err = db.GuildConfigManagerRemove(p.guildID, p.guildName, users, roles)
	}
	if err != nil {
		log.Errorf("guildManagerUpdate Error:%s,%s:%s", err, p.guildID, operator)
		return nil, nil, err
	}
	gmDB, err := db.GuildConfigManagerByGuildID(nil, p.guildID)
	if err != nil {
		log.Errorf("guildManagerUpdate read from db Error:%s,%s", err, p.guildID)
		return nil, nil, err
	}
	p.managers = gmDB.Managers
	p.managerRoles = gmDB.ManagerRoles
	return gmDB.Managers, gmDB.ManagerRoles, nil
}

func (p *guildConfigPresenter) guildExcludeUpdate(operator string, roles []string) ([]string, error) {
	var err error
	if operator == "add" {
		err = db.GuildConfigExcludeRolesAdd(p.guildID, p.guildName, roles)
	} else {
		err = db.GuildConfigExcludeRolesRemove(p.guildID, p.guildName, roles)
	}
	if err != nil {
		log.Errorf("guildExcludeUpdate Error:%s,%s:%s", err, p.guildID, operator)
		return nil, err
	}
	gmDB, err := db.GuildConfigManagerByGuildID(nil, p.guildID)
	if err != nil {
		log.Errorf("guildExcludeUpdate read from db Error:%s,%s", err, p.guildID)
		return nil, err
	}
	p.excludeRoles = gmDB.ExcludeRoles
	return gmDB.ExcludeRoles, nil
}
