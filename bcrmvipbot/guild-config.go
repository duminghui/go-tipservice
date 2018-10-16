// Package main provides ...
package main

import (
	"strings"

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
	if strings.Contains(bcrmVipConfig.Discord.SuperManagerIDs, userID) {
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

func gc(guildID string) *guildConfig {
	gcdb, _ := dbGuild.GuildConfigByID(nil, guildID)
	return (*guildConfig)(gcdb)
}
