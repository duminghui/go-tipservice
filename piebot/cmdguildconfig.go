// Package main provides ...
package main

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/bwmarrin/discordgo"
)

func (p *guildConfigPresenter) cmdMainPie(parts *msgParts) {
	// if !isBotManager(s, guild, m.Author.ID) {
	// return
	// }
	cntParts := parts.contents
	switch {
	case len(cntParts) == 0:
		p.cmdPieMainHelpHandler(parts)
	case cntParts[0] == "list":
		p.cmdPieManListHandler(parts)
	case cntParts[0] == "prefix":
		parts.contents = cntParts[1:]
		p.cmdPieMainPrefixHandler(parts)
	case cntParts[0] == "manager":
		parts.contents = cntParts[1:]
		p.cmdPieMainManagerHandler(parts)
	case cntParts[0] == "info":
		p.cmdPieMainInfoHandler(parts)
	default:
		p.cmdPieMainHelpHandler(parts)
	}
}

func (p *guildConfigPresenter) cmdPieMainInfoHandler(parts *msgParts) {
	coinConfigs := make([]*tmplValueMap, 0)
	for k, v := range p.prefixSymbol {
		channels := make([]string, 0)
		for _, channelID := range p.gccMap[v].channelIDs {
			channels = append(channels, fmt.Sprintf("<#%s>", channelID))
		}
		coinConfig := &tmplValueMap{
			"Prefix":   k,
			"Symbol":   v,
			"Channels": channels,
		}
		coinConfigs = append(coinConfigs, coinConfig)
	}
	managers := make([]string, 0, len(p.managers))
	for _, manager := range p.managers {
		member, _ := parts.s.State.Member(p.guildID, manager)
		managers = append(managers, member.User.Username)
	}
	roles := make([]string, 0, len(p.managerRoles))
	for _, roleID := range p.managerRoles {
		role, _ := parts.s.State.Role(p.guildID, roleID)
		roles = append(roles, role.Name)
	}
	msg := msgFromTmpl("pieMainInfo", tmplValueMap{
		"UserMention":  parts.m.Author.Mention(),
		"CoinConfigs":  coinConfigs,
		"Managers":     managers,
		"ManagerRoles": roles,
	})
	parts.channelMessageSend(msg)
}

func (p *guildConfigPresenter) cmdPieMainManagerHandler(parts *msgParts) {
	userMention := parts.m.Author.Mention()
	msg := msgFromTmpl("pieMainManagerUsage", tmplValueMap{
		"UserMention":     userMention,
		"IsShowUsageHint": true,
		"CmdName":         "manager",
	})
	contents := parts.contents
	if len(contents) < 2 {
		parts.channelMessageSend(msg)
		return
	}
	operator := contents[0]
	if operator != "add" && operator != "remove" {
		parts.channelMessageSend(msg)
		return
	}
	if len(parts.m.MentionRoles) == 0 && len(parts.m.Mentions) == 0 {
		parts.channelMessageSend(msg)
		return
	}
	users := make([]string, 0, len(parts.m.Mentions))
	for _, user := range parts.m.Mentions {
		users = append(users, user.ID)
	}
	updateUsers, updateRoles, err := p.guildManagerUpdate(operator, users, parts.m.MentionRoles)
	if err != nil {
		log.Error("cmdPieSetManagerHandler Error:", err)
		return
	}

	managers := make([]string, 0, len(updateUsers))
	for _, manager := range updateUsers {
		member, _ := parts.s.State.Member(p.guildID, manager)
		managers = append(managers, member.User.Username)
	}
	roles := make([]string, 0, len(updateRoles))
	for _, roleID := range updateRoles {
		role, _ := parts.s.State.Role(p.guildID, roleID)
		roles = append(roles, role.Name)
	}
	msg = msgFromTmpl("pieMainManagerInfo", tmplValueMap{
		"UserMention":  userMention,
		"Managers":     managers,
		"ManagerRoles": roles,
	})
	parts.channelMessageSend(msg)
}

func (p *guildConfigPresenter) isBotManager(s *discordgo.Session, guild *discordgo.Guild, userID string) bool {
	if userID == piebotConfig.Discord.SuperManagerID {
		return true
	}
	if userID == guild.OwnerID {
		return true
	}
	member, err := s.State.Member(guild.ID, userID)
	if err != nil {
		log.Error("isBotManager Error:", err)
		return false
	}
	if p.isManager(userID) || p.inManagerRoles(member.Roles) {
		return true
	}
	return false
}

func (p *guildConfigPresenter) cmdPieMainPrefixHandler(parts *msgParts) {
	contents := parts.contents
	userMention := parts.m.Author.Mention()
	if len(contents) != 2 {
		msg := msgFromTmpl("pieMainPrefixUsage", tmplValueMap{
			"UserMention":     userMention,
			"IsShowUsageHint": true,
			"CmdName":         "prefix",
		})
		parts.channelMessageSend(msg)
		return
	}
	symbol := symbolWrap(contents[0])
	prefix := prefixWrap(contents[1])
	if prefix == "?" {
		msg := msgFromTmpl("pieMainPrefixErr", userMention)
		parts.channelMessageSend(msg)
		return
	}
	if !hasSymbol(symbol) {
		msg := msgFromTmpl("pieMainPrefixSymbolNotExistErr", tmplValueMap{
			"UserMention": userMention,
			"Symbol":      symbol,
		})
		parts.channelMessageSend(msg)
		return
	}
	symbolTmp, _ := p.symbolByPrefix(prefix)

	if symbolTmp != "" {
		msg := msgFromTmpl("pieMainPrefixExistErr", tmplValueMap{
			"UserMention": userMention,
			"Prefix":      prefix,
			"Symbol":      symbolTmp,
		})
		parts.channelMessageSend(msg)
		return
	}

	err := p.updatePrefix(symbol, parts.prefix, prefix)
	if err != nil {
		log.Errorf("cmdPieSetPrefix Cache error:%s", err)
		return
	}
	msg := msgFromTmpl("pieMainPrefixSuccess", tmplValueMap{
		"UserMention": userMention,
		"Prefix":      prefix,
		"Symbol":      symbol,
	})
	parts.channelMessageSend(msg)
}

func (p *guildConfigPresenter) cmdPieManListHandler(parts *msgParts) {
	symbols := make([]string, 0, len(coinPresenters))
	for k := range coinPresenters {
		symbols = append(symbols, string(k))
	}
	msg := msgFromTmpl("pieMainListInfo", tmplValueMap{
		"UserMention": parts.m.Author.Mention(),
		"Symbols":     symbols,
	})
	parts.channelMessageSend(msg)
}

func (p *guildConfigPresenter) cmdPieMainHelpHandler(parts *msgParts) {
	msg := msgFromTmpl("pieMainHelpUsage", parts.m.Author.Mention())
	parts.channelMessageSend(msg)
}

func (p *guildConfigPresenter) cmdChannelHandler(parts *msgParts) {
	cmdPrefix := parts.prefix
	userMention := parts.m.Author.Mention()
	symbol := parts.symbol
	cmdUsage := &cmdUsageInfo{
		tmplName:        "channelUsage",
		IsShowUsageHint: true,
		CmdName:         "channel",
		UserMention:     userMention,
		Prefix:          string(cmdPrefix),
		Symbol:          string(symbol),
	}
	contents := parts.contents
	if len(contents) < 2 {
		parts.channelMessageSend(cmdUsage.String())
		return
	}
	operator := contents[0]
	if operator != "add" && operator != "remove" {
		parts.channelMessageSend(cmdUsage.String())
		return
	}
	str := strings.Join(contents[1:], "")
	exp := regexp.MustCompile(`<#(\d{18})>`)
	result := exp.FindAllStringSubmatch(str, -1)
	channels := make([]string, 0, len(result))
	for _, v := range result {
		channels = append(channels, v[1])
	}
	if len(channels) == 0 {
		parts.channelMessageSend(cmdUsage.String())
		return
	}
	finalChannels, err := p.guildChannelUpdate(symbol, operator, channels)
	if err != nil {
		log.Error("cmdSetChannelHandler Error:", err)
		return
	}
	msg := msgFromTmpl("channelOperatorSuccess", tmplValueMap{
		"UserMention": userMention,
		"Operator":    operator,
		"Symbol":      symbol,
		"Channels":    finalChannels,
	})
	parts.channelMessageSend(msg)
}
