// Package main provides ...
package main

import "sync"

type symbol string
type prefix string

type guildSymbolPresenter struct {
	symbol symbol
	*guildPresenter
	*coinPresenter
}

func (p *guildSymbolPresenter) guildChannelUpdate(sbl symbol, operator string, channels []string) ([]string, error) {
	var err error
	guildID := p.guildID
	guildName := p.guildName
	if operator == "add" {
		err = p.dbGuild.GuildCoinChannelAdd(guildID, guildName, string(sbl), channels)
	} else {
		err = p.dbGuild.GuildCoinChannelRemove(guildID, guildName, string(sbl), channels)
	}
	if err != nil {
		log.Errorf("GuildChannelUpdate Error:%s,%s:%s:%s", err, guildID, sbl, operator)
		return nil, err
	}
	gccDB, err := p.dbGuild.GuildCoinConfigBySymbol(nil, guildID, string(sbl))
	if err != nil {
		log.Errorf("update channel GuildCoinConfigBySymbol error:%s,%s:%s", err, p.guildID, sbl)
		return nil, err
	}
	guildSymbolCoinConfigs.update(guildID, sbl, gccDB)
	return gccDB.ChannelIDs, nil
}

type guildSymbolPresenterMap map[string]map[symbol]*guildSymbolPresenter

var guildSymbolPresenters = make(guildSymbolPresenterMap)

var gspmMU sync.RWMutex

func (gspm guildSymbolPresenterMap) gsp(guildID string, sbl symbol) *guildSymbolPresenter {
	gspmMU.RLock()
	sspm, ok := gspm[guildID]
	if !ok {
		gspmMU.RUnlock()
		gspmMU.Lock()
		sspm = make(map[symbol]*guildSymbolPresenter)
		gspm[guildID] = sspm
		gspmMU.Unlock()
		gspmMU.RLock()
	}
	gsp, ok := sspm[sbl]
	if ok {
		gspmMU.RUnlock()
		return gsp
	}
	gspmMU.RUnlock()
	gspmMU.Lock()
	gsp = new(guildSymbolPresenter)
	gsp.guildPresenter = guildPresenters.read(guildID)
	gsp.symbol = sbl
	gsp.coinPresenter = coinPresenters[sbl]
	gspm[guildID][sbl] = gsp
	gspmMU.Unlock()
	return gsp
}
