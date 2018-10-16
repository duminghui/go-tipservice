// Package db provides ...
package db

import (
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
)

func (db *DBGuild) cGuildConfig(session *mgo.Session) *mgo.Collection {
	return session.DB(db.database).C("guildconfig")
}

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

func (db *DBGuild) GuildConfigExcludeRolesRemove(guildID, guildName string, excluderoles []string) error {
	session := mgoSession.Clone()
	defer session.Close()
	config, err := db.GuildConfigByID(session, guildID)
	if err != nil {
		return err
	}
	if config == nil {
		return nil
	}
	if len(config.ExcludeRoles) == 0 {
		return nil
	}
	rolesMap := make(map[string]int)
	for _, role := range config.ExcludeRoles {
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
	col := db.cGuildConfig(session)
	err = col.Update(selector, data)
	return err
}

func (db *DBGuild) GuildConfigExcludeRolesAdd(guildID, guildName string, excluderoles []string) error {
	session := mgoSession.Clone()
	defer session.Close()
	config, err := db.GuildConfigByID(session, guildID)
	if err != nil {
		return err
	}
	col := db.cGuildConfig(session)
	if config == nil {
		data := &GuildConfig{
			GuildID:      guildID,
			GuildName:    guildName,
			ExcludeRoles: excluderoles,
		}
		err := col.Insert(data)
		return err
	}
	rolesMap := make(map[string]int)
	for _, role := range config.ExcludeRoles {
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

func (db *DBGuild) GuildConfigManagerRemove(guildID, guildName string, users, roles []string) error {
	session := mgoSession.Clone()
	defer session.Close()
	config, err := db.GuildConfigByID(session, guildID)
	if err != nil {
		return err
	}
	if config == nil {
		return nil
	}
	if len(config.Managers) == 0 && len(config.ManagerRoles) == 0 {
		return nil
	}
	managers := make([]string, 0)
	if len(config.Managers) != 0 {
		managerMap := make(map[string]int)
		for _, exitUser := range config.Managers {
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
	if len(config.ManagerRoles) != 0 {
		roleMap := make(map[string]int)
		for _, existRole := range config.ManagerRoles {
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
	col := db.cGuildConfig(session)
	err = col.Update(selector, data)
	return err
}

func (db *DBGuild) GuildConfigManagerAdd(guildID, guildName string, users, roles []string) error {
	session := mgoSession.Clone()
	defer session.Close()
	config, err := db.GuildConfigByID(session, guildID)
	if err != nil {
		return err
	}
	col := db.cGuildConfig(session)
	if config == nil {
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
	for _, user := range config.Managers {
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
	for _, role := range config.ManagerRoles {
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

func (db *DBGuild) GuildConfigList() ([]*GuildConfig, error) {
	session := mgoSession.Clone()
	defer session.Close()
	col := db.cGuildConfig(session)
	guildConfigs := make([]*GuildConfig, 0)
	err := col.Find(nil).All(&guildConfigs)
	return guildConfigs, err
}

func (db *DBGuild) GuildConfigByID(sessionIn *mgo.Session, guildID string) (*GuildConfig, error) {
	session, closer := session(sessionIn)
	defer closer()
	col := db.cGuildConfig(session)
	selector := &guildConfigSelector{
		GuildID: guildID,
	}
	config := new(GuildConfig)
	err := col.Find(selector).One(config)
	if err != nil {
		if err == mgo.ErrNotFound {
			return nil, nil
		}
		return nil, err
	}
	return config, nil
}

func (db *DBGuild) cGuildCoinConfig(session *mgo.Session) *mgo.Collection {
	return session.DB(db.database).C("guildcoinconfig")
}

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

func (db *DBGuild) GuildCoinChannelRemove(guildID, guildName, symbol string, channels []string) error {
	session := mgoSession.Clone()
	defer session.Close()
	coinConfig, err := db.GuildCoinConfigBySymbol(session, guildID, symbol)
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
	col := db.cGuildCoinConfig(session)
	err = col.Update(selector, data)
	return err
}

func (db *DBGuild) GuildCoinChannelAdd(guildID, guildName, symbol string, channels []string) error {
	session := mgoSession.Clone()
	defer session.Close()
	coinConfig, err := db.GuildCoinConfigBySymbol(session, guildID, symbol)
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
	col := db.cGuildCoinConfig(session)
	err = col.Update(selector, data)
	return err
}

func (db *DBGuild) GuildCoinUpdateCmdPrefix(guildID, guildName, symbol string, cmdPrefix string) error {
	session := mgoSession.Clone()
	defer session.Close()
	col := db.cGuildCoinConfig(session)
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

func (db *DBGuild) GuildCoinConfigList() ([]*GuildCoinConfig, error) {
	session := mgoSession.Clone()
	defer session.Close()
	col := db.cGuildCoinConfig(session)
	guildConfigs := make([]*GuildCoinConfig, 0)
	err := col.Find(nil).All(&guildConfigs)
	return guildConfigs, err
}

func (db *DBGuild) GuildCoinConfigBySymbol(sessionIn *mgo.Session, guildID, symbol string) (*GuildCoinConfig, error) {
	session, closer := session(sessionIn)
	defer closer()
	col := db.cGuildCoinConfig(session)
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
