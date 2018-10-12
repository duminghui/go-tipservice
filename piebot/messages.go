// Package main provides ...
package main

import (
	"fmt"
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
	prefix    prefix
	contents  []string
	isManager bool
}

func (p *msgParts) channelMessageSend(msg string) {
	p.s.ChannelMessageSend(p.channel.ID, msg)
}

func (p *msgParts) channelMessageSendComplex(msg *discordgo.MessageSend) {
	p.s.ChannelMessageSendComplex(p.channel.ID, msg)
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
	help := &cmdInfo{
		name:         "help",
		channelLimit: true,
		handler:      (*guildSymbolPresenter).cmdPieHelperHandler,
	}
	cmdInfoMap[help.name] = help
	pie := &cmdInfo{
		name:         "pie",
		channelLimit: true,
		handler:      (*guildSymbolPresenter).cmdPieHandler,
	}
	cmdInfoMap[pie.name] = pie
	deposit := &cmdInfo{
		name:         "deposit",
		channelLimit: true,
		handler:      (*guildSymbolPresenter).cmdDepositHandler,
	}
	cmdInfoMap[deposit.name] = deposit
	bal := &cmdInfo{
		name:         "bal",
		channelLimit: true,
		handler:      (*guildSymbolPresenter).cmdBalHandler,
	}
	cmdInfoMap[bal.name] = bal
	withdraw := &cmdInfo{
		name:         "withdraw",
		channelLimit: true,
		handler:      (*guildSymbolPresenter).cmdWithdrawHandler,
	}
	cmdInfoMap[withdraw.name] = withdraw

	setChannel := &cmdInfo{
		name:         "channel",
		managerCmd:   true,
		channelLimit: false,
		handler:      (*guildSymbolPresenter).cmdChannelHandler,
	}
	cmdInfoMap[setChannel.name] = setChannel

	pieAutoAdd := &cmdInfo{
		name:         "pieAutoAdd",
		managerCmd:   true,
		channelLimit: false,
		handler:      (*guildSymbolPresenter).cmdPieAutoHandler,
	}
	cmdInfoMap[pieAutoAdd.name] = pieAutoAdd

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
	guildID := channel.GuildID
	guild, err := guild(s, guildID)
	if err != nil {
		log.Error("messageCreateHandler guild Error:", err)
		return
	}
	gc := guildConfigs.gc(guildID)
	isManager := gc.isBotManager(s, guild, m.Author.ID)
	msgParts := &msgParts{
		s:         s,
		m:         m,
		channel:   channel,
		guild:     guild,
		contents:  cntParts[1:],
		isManager: isManager,
	}
	botMainCmd := fmt.Sprintf("%spie", piebotConfig.Discord.Prefix)
	if cntParts[0] == botMainCmd {
		if !isManager {
			return
		}

		// todo
		gp := guildPresenters.read(guildID)
		gp.cmdMainPie(msgParts)
		return
	}
	sccm := guildSymbolCoinConfigs.sccm(guildID)
	guildName := guild.Name
	if len(sccm) == 0 {
		log.Errorf("Prefix List Empty:[%s(%s)]", guildName, guildID)
		return
	}
	var pfx prefix
	var sbl symbol
	for k, v := range sccm {
		if strings.HasPrefix(m.Content, string(v.CmdPrefix)) {
			sbl = k
			pfx = prefix(v.CmdPrefix)
			break
		}
	}
	if pfx == "" {
		log.Errorf("can't find match prefix for:%s[%s]", guildName, guildID)
		return
	}
	// just only prefix
	if strings.Compare(string(pfx), cntParts[0]) == 0 {
		return
	}
	msgParts.prefix = pfx

	_, ok := coinInfos[string(sbl)]
	if !ok {
		log.Errorf("No Coin Infos for:[%s]", sbl)
		return
	}
	cmd := strings.Replace(cntParts[0], string(pfx), "", 1)
	if cmdInfo, ok := cmdInfoMap[cmd]; ok {
		gcc := sccm[sbl]
		isInChannel := gcc.inChannels(m.ChannelID)
		if cmdInfo.channelLimit && !isInChannel {
			return
		}
		if cmdInfo.managerCmd && !isManager {
			return
		}
		gsp := guildSymbolPresenters.gsp(guildID, sbl)
		cmdInfo.handler(gsp, msgParts)
	}
}
