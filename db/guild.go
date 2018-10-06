// Package db provides ...
package db

import (
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
)

const (
	colGuildManager = "guildmanager"
	dbGuildConfig   = "guild_config"
)

type GuildManager struct {
	GuildID      string   `bson:"guildid,omitempty"`
	Managers     []string `bson:"managers"`
	ManagerRoles []string `bson:"managerroles"`
}

func (gm *GuildManager) ManagerRemove(users, roles []string) error {
	session := mgoSession.Clone()
	defer session.Close()
	col := session.DB(dbGuildConfig).C(colGuildManager)
	selector := bson.M{
		"guildid": gm.GuildID,
	}
	manager := new(GuildManager)
	err := col.Find(selector).One(manager)
	if err != nil {
		return err
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
		"$set": &GuildManager{
			Managers:     managers,
			ManagerRoles: updateRoles,
		},
	}
	err = col.Update(selector, data)
	return err
}

func (gm *GuildManager) ManagerAdd(users, roles []string) error {
	session := mgoSession.Clone()
	defer session.Close()
	col := session.DB(dbGuildConfig).C(colGuildManager)
	selector := bson.M{
		"guildid": gm.GuildID,
	}
	manager := new(GuildManager)
	err := col.Find(selector).One(manager)
	if err != nil && err != mgo.ErrNotFound {
		return err
	}
	if err != nil && err == mgo.ErrNotFound {
		data := &GuildManager{
			GuildID:      gm.GuildID,
			Managers:     users,
			ManagerRoles: roles,
		}
		err := col.Insert(data)
		return err
	}
	//Update
	managers := make([]string, len(manager.Managers))
	copy(managers, manager.Managers)
	for _, user := range users {
		add := true
		for _, existUser := range manager.Managers {
			if user == existUser {
				add = false
				break
			}
		}
		if add {
			managers = append(managers, user)
		}
	}
	addRoles := make([]string, 0)
	roleMap := make(map[string]int)
	for _, role := range roles {
		roleMap[role] = 1
	}
	for _, role := range manager.ManagerRoles {
		roleMap[role] = 1
	}
	for k := range roleMap {
		addRoles = append(addRoles, k)
	}
	data := bson.M{
		"$set": &GuildManager{
			Managers:     managers,
			ManagerRoles: addRoles,
		},
	}
	err = col.Update(selector, data)
	return err
}

func (gm *GuildManager) IsManager(userID string) bool {
	for _, member := range gm.Managers {
		if member == userID {
			return true
		}
	}
	return false
}

func (gm *GuildManager) InManagerRoles(userRoles []string) bool {
	for _, role := range gm.ManagerRoles {
		for _, userRole := range userRoles {
			if role == userRole {
				return true
			}
		}
	}
	return false
}

func GuildManagerList() ([]*GuildManager, error) {
	session := mgoSession.Clone()
	defer session.Close()
	col := session.DB(dbGuildConfig).C(colGuildManager)
	guildManagers := make([]*GuildManager, 0)
	err := col.Find(nil).All(&guildManagers)
	return guildManagers, err
}

func GuildManagerByGuildID(guildID string) (*GuildManager, error) {
	session := mgoSession.Clone()
	defer session.Close()
	col := session.DB(dbGuildConfig).C(colGuildManager)
	selector := bson.M{
		"guildid": guildID,
	}
	manager := new(GuildManager)
	err := col.Find(selector).One(manager)
	return manager, err
}

const colGuildCoinConfig = "guildcoinconfig"

type GuildCoinConfig struct {
	GuildID    string   `bson:"guildid,omitempty"`
	Symbol     string   `bson:"symbol,omitempty"`
	CmdPrefix  string   `bson:"cmdprefix,omitempty"`
	ChannelIDs []string `bson:"channelids,omitempty"`
}

func (gc *GuildCoinConfig) InChannels(channelID string) bool {
	if len(gc.ChannelIDs) == 0 {
		return true
	}
	for _, channel := range gc.ChannelIDs {
		if channel == channelID {
			return true
		}
	}
	return false
}

type guildCoinConfigChannel struct {
	ChannelIDs []string `bson:"channelids"`
}

func (gc *GuildCoinConfig) ChannelRemove(channels []string) error {
	session := mgoSession.Clone()
	defer session.Close()
	col := session.DB(dbGuildConfig).C(colGuildCoinConfig)
	selector := &GuildCoinConfig{
		GuildID: gc.GuildID,
		Symbol:  gc.Symbol,
	}
	coinConfig := new(GuildCoinConfig)
	err := col.Find(selector).One(coinConfig)
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
	data := bson.M{
		"$set": &guildCoinConfigChannel{
			ChannelIDs: addChannels,
		},
	}
	err = col.Update(selector, data)
	return err
}

func (gc *GuildCoinConfig) ChannelAdd(channels []string) error {
	session := mgoSession.Clone()
	defer session.Close()
	col := session.DB(dbGuildConfig).C(colGuildCoinConfig)
	selector := &GuildCoinConfig{
		GuildID: gc.GuildID,
		Symbol:  gc.Symbol,
	}
	coinConfig := new(GuildCoinConfig)
	err := col.Find(selector).One(coinConfig)
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
	data := bson.M{
		"$set": &GuildCoinConfig{
			ChannelIDs: addChannels,
		},
	}
	err = col.Update(selector, data)
	return err
}

func (gc *GuildCoinConfig) UpdateCmdPrefix(cmdPrefix string) error {
	session := mgoSession.Clone()
	defer session.Close()
	col := session.DB(dbGuildConfig).C(colGuildCoinConfig)
	selector := &GuildCoinConfig{
		GuildID: gc.GuildID,
		Symbol:  gc.Symbol,
	}
	data := bson.M{
		"$set": &GuildCoinConfig{
			CmdPrefix: cmdPrefix,
		},
	}
	err := col.Update(selector, data)
	if err != nil && err != mgo.ErrNotFound {
		return nil
	}
	if err != nil && err == mgo.ErrNotFound {
		data := &GuildCoinConfig{
			GuildID:   gc.GuildID,
			Symbol:    gc.Symbol,
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

func GuildCoinConfigBySymbol(guildID, symbol string) (*GuildCoinConfig, error) {
	session := mgoSession.Clone()
	defer session.Close()
	col := session.DB(dbGuildConfig).C(colGuildCoinConfig)
	guildConfig := new(GuildCoinConfig)
	selector := &GuildCoinConfig{
		GuildID: guildID,
		Symbol:  symbol,
	}
	err := col.Find(selector).One(guildConfig)
	return guildConfig, err
}
