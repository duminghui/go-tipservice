// Package main provides ...
package main

import "github.com/duminghui/go-tipservice/db"

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
func gcc(guildID string) *guildCoinConfig {
	gccdb, _ := dbGuild.GuildCoinConfigBySymbol(nil, guildID, "BCRM")
	return (*guildCoinConfig)(gccdb)
}
