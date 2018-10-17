// Package main provides ...
package main

import (
	"regexp"

	"github.com/bwmarrin/discordgo"
)

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

func member(s *discordgo.Session, guildID, userID string) (*discordgo.Member, error) {
	member, err := s.State.Member(guildID, userID)
	if err != nil {
		member, err = s.GuildMember(guildID, userID)
		if err != nil {
			log.Errorf("Get Member Error:%s:%s:%s", err, guildID, userID)
			return nil, err
		}
		return member, nil
	}

	return member, nil
}

// UserChannelPermissions
func userChannelPermissions(s *discordgo.Session, userID, channelID string) (int, error) {
	permission, err := s.State.UserChannelPermissions(userID, channelID)
	if err != nil {
		log.Errorf("userChannelPermissions Error:%s,uid:%s,cid:%s", err, userID, channelID)
		// s.UserChannelPermissions same as above
	}
	return permission, err
}

func presence(s *discordgo.Session, guildID, userID string) (*discordgo.Presence, error) {
	presence, err := s.State.Presence(guildID, userID)
	if err != nil {
		log.Errorf("presence Error:%s,uid:%s,gid:%s", err, userID, guildID)
	}
	return presence, err
}

func message(s *discordgo.Session, channelID, messageID string) (*discordgo.Message, error) {
	message, err := s.State.Message(channelID, messageID)
	if err != nil {
		message, err = s.ChannelMessage(channelID, messageID)
		if err != nil {
			log.Errorf("message Error:%s,cid:%s,mid:%s", err, channelID, messageID)
			return nil, err
		}
		return message, nil
	}
	return message, nil
}

//Direct Message
func directMessageComplx(s *discordgo.Session, userID, content string, embed *discordgo.MessageEmbed) (*discordgo.Message, error) {
	channel, err := s.UserChannelCreate(userID)
	if err != nil {
		return nil, err
	}
	msgSend := &discordgo.MessageSend{
		Content: content,
		Embed:   embed,
	}
	msg, err := s.ChannelMessageSendComplex(channel.ID, msgSend)
	if err != nil {
		log.Errorf("directMessageComplx Error:%s[uID:%s]", err, userID)
		return nil, err
	}
	return msg, nil
}

func channelMessageEditComplx(s *discordgo.Session, channelID, msgID, content string, embed *discordgo.MessageEmbed) (*discordgo.Message, error) {
	msgEdit := discordgo.NewMessageEdit(channelID, msgID)
	if content != "" {
		msgEdit.SetContent(content)
	}
	msgEdit.SetEmbed(embed)
	msg, err := s.ChannelMessageEditComplex(msgEdit)
	if err != nil {
		log.Errorf("channelMessageEditComplx Error:%s[cID:%s,mID:%s]", err, channelID, msgID)
		return nil, err
	}
	return msg, nil
}

func channelIDsFromContent(content string) []string {
	exp := regexp.MustCompile(`<#(\d{18})>`)
	result := exp.FindAllStringSubmatch(content, -1)
	channels := make([]string, 0, len(result))
	for _, v := range result {
		channels = append(channels, v[1])
	}
	return channels
}
