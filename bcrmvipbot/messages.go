// Package main provides ...
package main

import (
	"strings"

	"github.com/bwmarrin/discordgo"
)

type cmdGuildSymbolHandler func(*presenter, *msgParts)

type msgParts struct {
	s         *discordgo.Session
	m         *discordgo.MessageCreate
	channel   *discordgo.Channel
	guild     *discordgo.Guild
	contents  []string
	isManager bool
}

func (p *msgParts) channelMessageSend(msg string) (*discordgo.Message, error) {
	return p.s.ChannelMessageSend(p.channel.ID, msg)
}

func (p *msgParts) channelMessageSendEmbed(embed *discordgo.MessageEmbed) (*discordgo.Message, error) {
	return p.s.ChannelMessageSendEmbed(p.channel.ID, embed)
}

func (p *msgParts) channelMessageSendComplex(content string, embed *discordgo.MessageEmbed) (*discordgo.Message, error) {
	msg := &discordgo.MessageSend{
		Content: content,
		Embed:   embed,
	}
	return p.s.ChannelMessageSendComplex(p.channel.ID, msg)
}

type cmdInfo struct {
	name         string
	managerCmd   bool
	channelLimit bool
	handler      cmdGuildSymbolHandler
}

var cmdInfoMap = make(map[string]*cmdInfo)

var cmdChannel = "channel"

func reigsterBotCmdHandler() {
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	log.Info(m.Content)
	pfx := bcrmVipConfig.Discord.Prefix
	if !strings.HasPrefix(m.Content, pfx) {
		return
	}
	if m.Type != discordgo.MessageTypeDefault {
		return
	}
	if m.Author.Bot || m.Author.ID == s.State.User.ID {
		return
	}
	cntParts := strings.Fields(m.Content)
	if len(cntParts) == 0 {
		return
	}
	channel, err := channel(s, m.ChannelID)
	if err != nil {
		log.Error("messageCreateHandler channel Error:", err)
		return
	}
	if channel.Type == discordgo.ChannelTypeDM {
		return
	}
	guildID := channel.GuildID
	guild, err := guild(s, guildID)
	if err != nil {
		log.Error("messageCreateHandler guild Error:", err)
		return
	}
	gc := gc(guildID)
	if gc == nil {
		return
	}
	isManager := gc.isBotManager(s, guild, m.Author.ID)
	if !isManager {
		return
	}
	msgParts := &msgParts{
		s:         s,
		m:         m,
		channel:   channel,
		guild:     guild,
		contents:  cntParts[1:],
		isManager: isManager,
	}
	cmd := strings.Replace(cntParts[0], pfx, "", 1)
	if cmdInfo, ok := cmdInfoMap[cmd]; ok {
		cmdInfo.handler(psn, msgParts)
	}
}
