// Package main provides ...
package main

import (
	"sync"

	"github.com/duminghui/go-tipservice/db"
)

var dbGuild = db.NewDBGuild()

type guildPresenter struct {
	guildID   string
	guildName string
	dbGuild   *db.DBGuild
}

func (p *guildPresenter) updatePrefix(sbl symbol, oldPrefix, newPrefix prefix) error {
	guildID := p.guildID
	err := p.dbGuild.GuildCoinUpdateCmdPrefix(guildID, p.guildName, string(sbl), string(newPrefix))
	if err != nil {
		log.Errorf("updatePrefix err:%s,%s,%s,%s", err, guildID, sbl, newPrefix)
		return err
	}
	gccDB, err := p.dbGuild.GuildCoinConfigBySymbol(nil, guildID, string(sbl))
	if err != nil {
		log.Errorf("update prefix GuildCoinConfigBySymbol error:%s,%s:%s", err, p.guildID, sbl)
		return err
	}
	guildSymbolCoinConfigs.update(guildID, sbl, gccDB)
	return nil
}

func (p *guildPresenter) guildManagerUpdate(operator string, users, roles []string) ([]string, []string, error) {
	var err error
	guildID := p.guildID
	guildName := p.guildName
	if operator == "add" {
		err = p.dbGuild.GuildConfigManagerAdd(guildID, guildName, users, roles)
	} else {
		err = p.dbGuild.GuildConfigManagerRemove(guildID, guildName, users, roles)
	}
	if err != nil {
		log.Errorf("guildManagerUpdate Error:%s,%s:%s", err, guildID, operator)
		return nil, nil, err
	}
	gmDB, err := dbGuild.GuildConfigManagerByGuildID(nil, guildID)
	if err != nil {
		log.Errorf("guildManagerUpdate read from db Error:%s,%s", err, guildID)
		return nil, nil, err
	}
	guildConfigs.update(guildID, gmDB)
	return gmDB.Managers, gmDB.ManagerRoles, nil
}

func (p *guildPresenter) guildExcludeUpdate(operator string, roles []string) ([]string, error) {
	guildID := p.guildID
	guildName := p.guildName
	var err error
	if operator == "add" {
		err = p.dbGuild.GuildConfigExcludeRolesAdd(guildID, guildName, roles)
	} else {
		err = p.dbGuild.GuildConfigExcludeRolesRemove(guildID, guildName, roles)
	}
	if err != nil {
		log.Errorf("guildExcludeUpdate Error:%s,%s:%s", err, guildID, operator)
		return nil, err
	}
	gmDB, err := dbGuild.GuildConfigManagerByGuildID(nil, guildID)
	if err != nil {
		log.Errorf("guildExcludeUpdate read from db Error:%s,%s", err, p.guildID)
		return nil, err
	}
	guildConfigs.update(guildID, gmDB)
	return gmDB.ExcludeRoles, nil
}

type guildPresenterMap map[string]*guildPresenter

var guildPresenters = make(guildPresenterMap)

var gpMU sync.RWMutex

func (gpm guildPresenterMap) read(guildID string) *guildPresenter {
	gpMU.RLock()
	gp, ok := gpm[guildID]
	if ok {
		gpMU.RUnlock()
		return gp
	}
	gpMU.RUnlock()
	gpMU.Lock()
	gp = &guildPresenter{
		guildID: guildID,
		dbGuild: dbGuild,
	}
	gpMU.Unlock()
	return gp
}
