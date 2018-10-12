// Package main provides ...
package main

import "github.com/bwmarrin/discordgo"

func channel(s *discordgo.Session, channelID string) (*discordgo.Channel, error) {
	channel, err := s.State.Channel(channelID)
	if err != nil {
		channel, err = s.Channel(channelID)
		if err != nil {
			log.Errorf("Get Channel Error:%s:%s", err, channelID)
			return nil, err
		}
	}
	return channel, err
}

func guild(s *discordgo.Session, guildID string) (*discordgo.Guild, error) {
	guild, err := s.State.Guild(guildID)
	if err != nil {
		guild, err = s.Guild(guildID)
		if err != nil {
			log.Errorf("Get guild Error:%s:%s", err, guildID)
			return nil, err
		}
	}
	return guild, err
}

func role(s *discordgo.Session, guildID, roldID string) (*discordgo.Role, error) {
	role, err := s.State.Role(guildID, roldID)
	if err != nil {
		roles, err := s.GuildRoles(guildID)
		if err != nil {
			log.Errorf("Get Role Error:%s:%s:%s", err, guildID, roldID)
			return nil, err
		}
		for _, roleTmp := range roles {
			if roleTmp.ID == roldID {
				return roleTmp, nil
			}
		}
		return nil, nil
	}
	return role, nil
}
