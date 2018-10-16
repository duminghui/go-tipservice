// Package main provides ...
package main

import (
	"strings"
	"sync"

	"github.com/bwmarrin/discordgo"
	"github.com/duminghui/go-tipservice/db"
)

type guildConfig db.GuildConfig

func (c *guildConfig) isManager(userID string) bool {
	for _, member := range c.Managers {
		if member == userID {
			return true
		}
	}
	return false
}

func (c *guildConfig) inManagerRoles(userRoles []string) bool {
	for _, role := range c.ManagerRoles {
		for _, userRole := range userRoles {
			if role == userRole {
				return true
			}
		}
	}
	return false
}

func (c *guildConfig) isBotManager(s *discordgo.Session, guild *discordgo.Guild, userID string) bool {
	if strings.Contains(piebotConfig.Discord.SuperManagerIDs, userID) {
		return true
	}
	if userID == guild.OwnerID {
		return true
	}
	member, err := member(s, guild.ID, userID)
	if err != nil {
		log.Error("isBotManager Error:", err)
		return false
	}
	if c.isManager(userID) || c.inManagerRoles(member.Roles) {
		return true
	}
	return false
}

type guildConfigMap map[string]*guildConfig

var guildConfigs = make(guildConfigMap)

func readGuildConfigsFromDB() {
	guildConfigList, err := dbGuild.GuildConfigList()
	if err != nil {
		return
	}
	for _, v := range guildConfigList {
		guildConfigs[v.GuildID] = (*guildConfig)(v)
	}
}

func (gcm guildConfigMap) initGuildConfig(guildID string) *guildConfig {
	c := new(guildConfig)
	c.Managers = []string{}
	c.ManagerRoles = []string{}
	c.ExcludeRoles = []string{}
	gcm[guildID] = c
	return c
}

var gcmMU sync.RWMutex

func (gcm guildConfigMap) update(guildID string, gc *db.GuildConfig) {
	gcmMU.Lock()
	defer gcmMU.Unlock()
	gcm[guildID] = (*guildConfig)(gc)
}

func (gcm guildConfigMap) gc(guildID string) *guildConfig {
	gcmMU.RLock()
	gc, ok := guildConfigs[guildID]
	if ok {
		gcmMU.RUnlock()
		return gc
	}
	gcmMU.RUnlock()
	gcmMU.Lock()
	gc = guildConfigs.initGuildConfig(guildID)
	gcmMU.Unlock()
	return gc
}
