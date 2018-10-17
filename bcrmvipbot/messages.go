// Package main provides ...
package main

import (
	"strings"

	"github.com/bwmarrin/discordgo"
)

type cmdGuildSymbolHandler func(*guildSymbolPresenter, *msgParts)

// type cmdHandler func(*discordgo.Session, *discordgo.MessageCreate, *msgParts)

type msgParts struct {
	s         *discordgo.Session
	m         *discordgo.MessageCreate
	channel   *discordgo.Channel
	guild     *discordgo.Guild
	contents  []string
	isManager bool
	prefix    string
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
	guildLimit   bool
	handler      cmdGuildSymbolHandler
}

var cmdInfoMap = make(map[string]*cmdInfo)

func registerCmds() {
	registerSymbolVipCmds()
}

func registerSymbolCmd(name string, channelLimit, managerCmd, guildLimit bool, handler cmdGuildSymbolHandler) {
	cmdInfo := &cmdInfo{
		name:         name,
		channelLimit: channelLimit,
		managerCmd:   managerCmd,
		handler:      handler,
	}
	cmdInfoMap[name] = cmdInfo
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
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
	isManager := gc.isBotManager(s, guild, m.Author.ID)
	msgParts := &msgParts{
		s:         s,
		m:         m,
		channel:   channel,
		guild:     guild,
		contents:  cntParts[1:],
		isManager: isManager,
	}
	gcc := gcc(guild.ID)
	if gcc == nil {
		return
	}
	pfx := gcc.CmdPrefix
	if !strings.HasPrefix(m.Content, pfx) {
		return
	}
	if strings.Compare(pfx, cntParts[0]) == 0 {
		return
	}
	msgParts.prefix = pfx

	cmd := strings.Replace(cntParts[0], string(pfx), "", 1)
	presenter.guildID = guild.ID
	if cmdInfo, ok := cmdInfoMap[cmd]; ok {
		isInChannel := gcc.inChannels(m.ChannelID)
		if cmdInfo.channelLimit && !isInChannel {
			return
		}
		if cmdInfo.managerCmd && !isManager {
			return
		}
		cmdInfo.handler(presenter, msgParts)
	}
}
