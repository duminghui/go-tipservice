// Package db provides ...
package db

import (
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
)

const (
	colGuildManager = "guildconfig"
	dbGuildConfig   = "guild_config"
)

type GuildConfig struct {
	GuildID      string   `bson:"guildid"`
	GuildName    string   `bson:"guildname"`
	Managers     []string `bson:"managers"`
	ManagerRoles []string `bson:"managerroles"`
	ExcludeRoles []string `bson:"excluderoles"`
}

type guildConfigSelector struct {
	GuildID string `bson:"guildid"`
}

type guildConfigManagers struct {
	GuildName    string   `bson:"guildname"`
	Managers     []string `bson:"managers"`
	ManagerRoles []string `bson:"managerroles"`
}

type guildConfigExcludeRoles struct {
	GuildName    string   `bson:"guildname"`
	ExcludeRoles []string `bson:"excluderoles"`
}

func GuildConfigExcludeRolesRemove(guildID, guildName string, excluderoles []string) error {
	session := mgoSession.Clone()
	defer session.Close()
	manager, err := GuildConfigManagerByGuildID(session, guildID)
	if err != nil {
		return err
	}
	if manager == nil {
		return nil
	}
	if len(manager.ExcludeRoles) == 0 {
		return nil
	}
	rolesMap := make(map[string]int)
	for _, role := range manager.ExcludeRoles {
		rolesMap[role] = 1
	}
	for _, role := range excluderoles {
		delete(rolesMap, role)
	}
	finalRoles := make([]string, 0)
	for k := range rolesMap {
		finalRoles = append(finalRoles, k)
	}
	selector := &guildConfigSelector{
		GuildID: guildID,
	}
	data := bson.M{
		"$set": &guildConfigExcludeRoles{
			GuildName:    guildName,
			ExcludeRoles: finalRoles,
		},
	}
	col := session.DB(dbGuildConfig).C(colGuildManager)
	err = col.Update(selector, data)
	return err
}

func GuildConfigExcludeRolesAdd(guildID, guildName string, excluderoles []string) error {
	session := mgoSession.Clone()
	defer session.Close()
	manager, err := GuildConfigManagerByGuildID(session, guildID)
	if err != nil {
		return err
	}
	col := session.DB(dbGuildConfig).C(colGuildManager)
	if manager == nil {
		data := &GuildConfig{
			GuildID:      guildID,
			GuildName:    guildName,
			ExcludeRoles: excluderoles,
		}
		err := col.Insert(data)
		return err
	}
	rolesMap := make(map[string]int)
	for _, role := range manager.ExcludeRoles {
		rolesMap[role] = 1
	}
	for _, role := range excluderoles {
		rolesMap[role] = 1
	}
	finalRoles := make([]string, 0)
	for k := range rolesMap {
		finalRoles = append(finalRoles, k)
	}
	selector := &guildConfigSelector{
		GuildID: guildID,
	}
	data := bson.M{
		"$set": &guildConfigExcludeRoles{
			GuildName:    guildName,
			ExcludeRoles: finalRoles,
		},
	}
	err = col.Update(selector, data)
	return err
}

func GuildConfigManagerRemove(guildID, guildName string, users, roles []string) error {
	session := mgoSession.Clone()
	defer session.Close()
	manager, err := GuildConfigManagerByGuildID(session, guildID)
	if err != nil {
		return err
	}
	if manager == nil {
		return nil
	}
	if len(manager.Managers) == 0 && len(manager.ManagerRoles) == 0 {
		return nil
	}
	managers := make([]string, 0)
	if len(manager.Managers) != 0 {
		managerMap := make(map[string]int)
		for _, exitUser := range manager.Managers {
			managerMap[exitUser] = 1
		}
		for _, user := range users {
			delete(managerMap, user)
		}
		for k := range managerMap {
			managers = append(managers, k)
		}
	}
	updateRoles := make([]string, 0)
	if len(manager.ManagerRoles) != 0 {
		roleMap := make(map[string]int)
		for _, existRole := range manager.ManagerRoles {
			roleMap[existRole] = 1
		}
		for _, role := range roles {
			delete(roleMap, role)
		}
		for k := range roleMap {
			updateRoles = append(updateRoles, k)
		}
	}
	data := bson.M{
		"$set": &guildConfigManagers{
			GuildName:    guildName,
			Managers:     managers,
			ManagerRoles: updateRoles,
		},
	}
	selector := &guildConfigSelector{
		GuildID: guildID,
	}
	col := session.DB(dbGuildConfig).C(colGuildManager)
	err = col.Update(selector, data)
	return err
}

func GuildConfigManagerAdd(guildID, guildName string, users, roles []string) error {
	session := mgoSession.Clone()
	defer session.Close()
	manager, err := GuildConfigManagerByGuildID(session, guildID)
	if err != nil {
		return err
	}
	col := session.DB(dbGuildConfig).C(colGuildManager)
	if manager == nil {
		data := &GuildConfig{
			GuildID:      guildID,
			GuildName:    guildName,
			Managers:     users,
			ManagerRoles: roles,
		}
		err := col.Insert(data)
		return err
	}
	//Update
	managerMap := make(map[string]int)
	for _, user := range users {
		managerMap[user] = 1
	}
	for _, user := range manager.Managers {
		managerMap[user] = 1
	}
	managers := make([]string, 0)
	for k := range managerMap {
		managers = append(managers, k)
	}
	roleMap := make(map[string]int)
	for _, role := range roles {
		roleMap[role] = 1
	}
	for _, role := range manager.ManagerRoles {
		roleMap[role] = 1
	}
	addRoles := make([]string, 0)
	for k := range roleMap {
		addRoles = append(addRoles, k)
	}
	data := bson.M{
		"$set": &guildConfigManagers{
			GuildName:    guildName,
			Managers:     managers,
			ManagerRoles: addRoles,
		},
	}
	selector := &guildConfigSelector{
		GuildID: guildID,
	}
	err = col.Update(selector, data)
	return err
}

func GuildConfigManagerList() ([]*GuildConfig, error) {
	session := mgoSession.Clone()
	defer session.Close()
	col := session.DB(dbGuildConfig).C(colGuildManager)
	guildManagers := make([]*GuildConfig, 0)
	err := col.Find(nil).All(&guildManagers)
	return guildManagers, err
}

func GuildConfigManagerByGuildID(sessionIn *mgo.Session, guildID string) (*GuildConfig, error) {
	session, closer := session(sessionIn)
	defer closer()
	col := session.DB(dbGuildConfig).C(colGuildManager)
	selector := &guildConfigSelector{
		GuildID: guildID,
	}
	manager := new(GuildConfig)
	err := col.Find(selector).One(manager)
	if err != nil {
		if err == mgo.ErrNotFound {
			return nil, nil
		}
		return nil, err
	}
	return manager, nil
}

const colGuildCoinConfig = "guildcoinconfig"

type guildCoinConfigSelector struct {
	GuildID string `bson:"guildid"`
	Symbol  string `bson:"symbol"`
}

type GuildCoinConfig struct {
	GuildID    string   `bson:"guildid"`
	GuildName  string   `bson:"guildname"`
	Symbol     string   `bson:"symbol"`
	CmdPrefix  string   `bson:"cmdprefix"`
	ChannelIDs []string `bson:"channelids"`
}

type guildCoinConfigCmdPrefix struct {
	GuildName string `bson:"guildname"`
	CmdPrefix string `bson:"cmdprefix"`
}

type guildCoinConfigChannel struct {
	GuildName  string   `bson:"guildname"`
	ChannelIDs []string `bson:"channelids"`
}

func GuildCoinChannelRemove(guildID, guildName, symbol string, channels []string) error {
	session := mgoSession.Clone()
	defer session.Close()
	coinConfig, err := GuildCoinConfigBySymbol(session, guildID, symbol)
	if err != nil {
		return err
	}
	channelMap := make(map[string]int)
	for _, channel := range coinConfig.ChannelIDs {
		channelMap[channel] = 1
	}
	for _, channel := range channels {
		delete(channelMap, channel)
	}
	addChannels := make([]string, 0)
	for k := range channelMap {
		addChannels = append(addChannels, k)
	}
	selector := &guildCoinConfigSelector{
		GuildID: guildID,
		Symbol:  symbol,
	}
	data := bson.M{
		"$set": &guildCoinConfigChannel{
			GuildName:  guildName,
			ChannelIDs: addChannels,
		},
	}
	col := session.DB(dbGuildConfig).C(colGuildCoinConfig)
	err = col.Update(selector, data)
	return err
}

func GuildCoinChannelAdd(guildID, guildName, symbol string, channels []string) error {
	session := mgoSession.Clone()
	defer session.Close()
	coinConfig, err := GuildCoinConfigBySymbol(session, guildID, symbol)
	if err != nil {
		return err
	}
	addChannelMap := make(map[string]int)
	for _, channel := range coinConfig.ChannelIDs {
		addChannelMap[channel] = 1
	}
	for _, channel := range channels {
		addChannelMap[channel] = 1
	}
	addChannels := make([]string, 0)
	for k := range addChannelMap {
		addChannels = append(addChannels, k)
	}
	selector := &guildCoinConfigSelector{
		GuildID: guildID,
		Symbol:  symbol,
	}
	data := bson.M{
		"$set": &guildCoinConfigChannel{
			GuildName:  guildName,
			ChannelIDs: addChannels,
		},
	}
	col := session.DB(dbGuildConfig).C(colGuildCoinConfig)
	err = col.Update(selector, data)
	return err
}

func GuildCoinUpdateCmdPrefix(guildID, guildName, symbol string, cmdPrefix string) error {
	session := mgoSession.Clone()
	defer session.Close()
	col := session.DB(dbGuildConfig).C(colGuildCoinConfig)
	selector := &guildCoinConfigSelector{
		GuildID: guildID,
		Symbol:  symbol,
	}
	data := bson.M{
		"$set": &guildCoinConfigCmdPrefix{
			GuildName: guildName,
			CmdPrefix: cmdPrefix,
		},
	}
	err := col.Update(selector, data)
	if err != nil && err != mgo.ErrNotFound {
		return err
	}
	if err != nil && err == mgo.ErrNotFound {
		data := &GuildCoinConfig{
			GuildID:   guildID,
			GuildName: guildName,
			Symbol:    symbol,
			CmdPrefix: cmdPrefix,
		}
		err := col.Insert(data)
		return err
	}
	return nil
}

func GuildCoinConfigList() ([]*GuildCoinConfig, error) {
	session := mgoSession.Clone()
	defer session.Close()
	col := session.DB(dbGuildConfig).C(colGuildCoinConfig)
	guildConfigs := make([]*GuildCoinConfig, 0)
	err := col.Find(nil).All(&guildConfigs)
	return guildConfigs, err
}

func GuildCoinConfigBySymbol(sessionIn *mgo.Session, guildID, symbol string) (*GuildCoinConfig, error) {
	session, closer := session(sessionIn)
	defer closer()
	col := session.DB(dbGuildConfig).C(colGuildCoinConfig)
	guildCoinConfig := new(GuildCoinConfig)
	selector := &guildCoinConfigSelector{
		GuildID: guildID,
		Symbol:  symbol,
	}
	err := col.Find(selector).One(guildCoinConfig)
	if err != nil && err != mgo.ErrNotFound {
		return nil, err
	}
	if err != nil && err == mgo.ErrNotFound {
		return nil, nil
	}
	return guildCoinConfig, nil
}
