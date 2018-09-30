// Package db provides ...
package db

import (
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
)

const (
	colGuildManager = "guildmanager"
	dbmain          = "main"
)

type GuildManager struct {
	GuildID      string   `bson:"guildid,omitempty"`
	Managers     []string `bson:"managers"`
	ManagerRoles []string `bson:"managerroles"`
}

func GuildManagerRemove(guildID string, users, roles []string) ([]string, []string, error) {
	session := mgoSession.Clone()
	defer session.Close()
	col := session.DB(dbmain).C(colGuildManager)
	selector := bson.M{
		"guildid": guildID,
	}
	manager := new(GuildManager)
	err := col.Find(selector).One(manager)
	if err != nil {
		return nil, nil, err
	}
	if len(manager.Managers) == 0 && len(manager.ManagerRoles) == 0 {
		return users, roles, nil
	}
	managers := make([]string, 0)
	if len(manager.Managers) != 0 {
		managerMap := make(map[string]int)
		for _, existUser := range manager.Managers {
			for _, user := range users {
				if existUser != user {
					managerMap[existUser]++
				}
			}
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
	return managers, updateRoles, err
}

func GuildManagerAdd(guildID string, users, roles []string) ([]string, []string, error) {
	session := mgoSession.Clone()
	defer session.Close()
	col := session.DB(dbmain).C(colGuildManager)
	selector := bson.M{
		"guildid": guildID,
	}
	manager := new(GuildManager)
	err := col.Find(selector).One(manager)
	if err != nil && err != mgo.ErrNotFound {
		return nil, nil, err
	}
	if err != nil && err == mgo.ErrNotFound {
		data := &GuildManager{
			GuildID:      guildID,
			Managers:     users,
			ManagerRoles: roles,
		}
		err := col.Insert(data)
		return users, roles, err
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
	return managers, addRoles, err
}

func (gc *GuildManager) IsManager(userID string) bool {
	for _, member := range gc.Managers {
		if member == userID {
			return true
		}
	}
	return false
}

func (gc *GuildManager) InManagerRoles(userRoles []string) bool {
	for _, role := range gc.ManagerRoles {
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
	col := session.DB(dbmain).C(colGuildManager)
	guildManagers := make([]*GuildManager, 0)
	err := col.Find(nil).All(&guildManagers)
	return guildManagers, err
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

type GuildCoinConfigChannel struct {
	ChannelIDs []string `bson:"channelids"`
}

func GuildChannelRemove(guildID, symbol string, channels []string) ([]string, error) {
	session := mgoSession.Clone()
	defer session.Close()
	col := session.DB(dbmain).C(colGuildCoinConfig)
	selector := &GuildCoinConfig{
		GuildID: guildID,
		Symbol:  symbol,
	}
	coinConfig := new(GuildCoinConfig)
	err := col.Find(selector).One(coinConfig)
	if err != nil {
		return nil, err
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
		"$set": &GuildCoinConfigChannel{
			ChannelIDs: addChannels,
		},
	}
	err = col.Update(selector, data)
	return addChannels, err
}

func GuildChannelAdd(guildID, symbol string, channels []string) ([]string, error) {
	session := mgoSession.Clone()
	defer session.Close()
	col := session.DB(dbmain).C(colGuildCoinConfig)
	selector := &GuildCoinConfig{
		GuildID: guildID,
		Symbol:  symbol,
	}
	coinConfig := new(GuildCoinConfig)
	err := col.Find(selector).One(coinConfig)
	if err != nil {
		return nil, err
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
	return addChannels, err
}

func GuildConfigList() ([]*GuildCoinConfig, error) {
	session := mgoSession.Clone()
	defer session.Close()
	col := session.DB(dbmain).C(colGuildCoinConfig)
	guildConfigs := make([]*GuildCoinConfig, 0)
	err := col.Find(nil).All(&guildConfigs)
	return guildConfigs, err
}

func GuildUpdateCmdPrefix(guildID, symbol, cmdPrefix string) error {
	session := mgoSession.Clone()
	defer session.Close()
	col := session.DB(dbmain).C(colGuildCoinConfig)
	selector := &GuildCoinConfig{
		GuildID: guildID,
		Symbol:  symbol,
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
			GuildID:   guildID,
			Symbol:    symbol,
			CmdPrefix: cmdPrefix,
		}
		err := col.Insert(data)
		return err
	}
	return nil
}
