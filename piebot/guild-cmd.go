// Package main provides ...
package main

import "fmt"

func (p *guildPresenter) cmdMainPie(parts *msgParts) {
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
	case cntParts[0] == "exclude":
		parts.contents = cntParts[1:]
		p.cmdPieMainExcludeHandler(parts)
	default:
		p.cmdPieMainHelpHandler(parts)
	}
}

func (p *guildPresenter) cmdPieMainExcludeHandler(parts *msgParts) {
	userMention := parts.m.Author.Mention()
	msg := msgFromTmpl("pieMainExcludeUsage", tmplValueMap{
		"UserMention":     userMention,
		"IsShowUsageHint": true,
		"CmdName":         "manager",
		"BotPrefix":       piebotConfig.Discord.Prefix,
	})
	if len(parts.contents) < 2 {
		parts.channelMessageSend(msg)
		return
	}
	operator := parts.contents[0]
	if operator != "add" && operator != "remove" {
		parts.channelMessageSend(msg)
		return
	}
	if len(parts.m.MentionRoles) == 0 {
		parts.channelMessageSend(msg)
		return
	}
	finalRoles, err := p.guildExcludeUpdate(operator, parts.m.MentionRoles)
	if err != nil {
		log.Error("cmdPieMainExcludeHandler Error:", err)
		return
	}
	roles := make([]string, 0, len(finalRoles))
	for _, roleID := range finalRoles {
		role, _ := role(parts.s, p.guildID, roleID)
		roles = append(roles, role.Name)
	}
	msg = msgFromTmpl("pieMainExcludeInfo", tmplValueMap{
		"UserMention":  userMention,
		"ExcludeRoles": roles,
	})
	parts.channelMessageSend(msg)
}

func (p *guildPresenter) cmdPieMainInfoHandler(parts *msgParts) {
	coinConfigs := make([]*tmplValueMap, 0)
	guildID := p.guildID
	scc := guildSymbolCoinConfigs[guildID]

	for k, v := range scc {
		channels := make([]string, 0)
		for _, channelID := range v.ChannelIDs {
			channels = append(channels, fmt.Sprintf("<#%s>", channelID))
		}
		coinConfig := &tmplValueMap{
			"Prefix":   v.CmdPrefix,
			"Symbol":   k,
			"Channels": channels,
		}
		coinConfigs = append(coinConfigs, coinConfig)
	}
	gc := guildConfigs.gc(guildID)
	managers := make([]string, 0, len(gc.Managers))
	for _, manager := range gc.Managers {
		member, _ := member(parts.s, p.guildID, manager)
		managers = append(managers, member.User.Username)
	}
	roles := make([]string, 0, len(gc.ManagerRoles))
	for _, roleID := range gc.ManagerRoles {
		role, _ := role(parts.s, p.guildID, roleID)
		roles = append(roles, role.Name)
	}

	excludeRoles := make([]string, 0, len(gc.ExcludeRoles))
	for _, roleID := range gc.ExcludeRoles {
		role, _ := role(parts.s, p.guildID, roleID)
		excludeRoles = append(roles, role.Name)
	}

	msg := msgFromTmpl("pieMainInfo", tmplValueMap{
		"UserMention":  parts.m.Author.Mention(),
		"CoinConfigs":  coinConfigs,
		"Managers":     managers,
		"ManagerRoles": roles,
		"ExcludeRoles": excludeRoles,
	})
	parts.channelMessageSend(msg)
}

func (p *guildPresenter) cmdPieMainManagerHandler(parts *msgParts) {
	userMention := parts.m.Author.Mention()
	msg := msgFromTmpl("pieMainManagerUsage", tmplValueMap{
		"UserMention":     userMention,
		"IsShowUsageHint": true,
		"CmdName":         "manager",
		"BotPrefix":       piebotConfig.Discord.Prefix,
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
		member, _ := member(parts.s, p.guildID, manager)
		managers = append(managers, member.User.Username)
	}
	roles := make([]string, 0, len(updateRoles))
	for _, roleID := range updateRoles {
		role, _ := role(parts.s, p.guildID, roleID)
		roles = append(roles, role.Name)
	}
	msg = msgFromTmpl("pieMainManagerInfo", tmplValueMap{
		"UserMention":  userMention,
		"Managers":     managers,
		"ManagerRoles": roles,
	})
	parts.channelMessageSend(msg)
}

func (p *guildPresenter) cmdPieMainPrefixHandler(parts *msgParts) {
	contents := parts.contents
	userMention := parts.m.Author.Mention()
	if len(contents) != 2 {
		msg := msgFromTmpl("pieMainPrefixUsage", tmplValueMap{
			"UserMention":     userMention,
			"IsShowUsageHint": true,
			"CmdName":         "prefix",
			"BotPrefix":       piebotConfig.Discord.Prefix,
		})
		parts.channelMessageSend(msg)
		return
	}
	sbl := symbol(contents[0])
	pfx := prefix(contents[1])
	botPrefix := piebotConfig.Discord.Prefix
	if pfx == prefix(botPrefix) {
		msg := msgFromTmpl("pieMainPrefixErr", tmplValueMap{
			"UserMention": userMention,
			"BotPrefix":   botPrefix,
		})
		parts.channelMessageSend(msg)
		return
	}
	if !hasSymbol(string(sbl)) {
		msg := msgFromTmpl("pieMainPrefixSymbolNotExistErr", tmplValueMap{
			"UserMention": userMention,
			"Symbol":      sbl,
		})
		parts.channelMessageSend(msg)
		return
	}

	sccm := guildSymbolCoinConfigs.sccm(p.guildID)
	symbolTmp, _ := sccm.symbolByPrefix(pfx)

	if symbolTmp != "" {
		msg := msgFromTmpl("pieMainPrefixExistErr", tmplValueMap{
			"UserMention": userMention,
			"Prefix":      pfx,
			"Symbol":      symbolTmp,
		})
		parts.channelMessageSend(msg)
		return
	}

	err := p.updatePrefix(sbl, parts.prefix, pfx)
	if err != nil {
		log.Errorf("cmdPieSetPrefix Cache error:%s", err)
		return
	}
	msg := msgFromTmpl("pieMainPrefixSuccess", tmplValueMap{
		"UserMention": userMention,
		"Prefix":      pfx,
		"Symbol":      sbl,
	})
	parts.channelMessageSend(msg)
}

func (p *guildPresenter) cmdPieManListHandler(parts *msgParts) {
	symbols := make([]string, 0, len(coinInfos))
	for k := range coinInfos {
		symbols = append(symbols, k)
	}
	msg := msgFromTmpl("pieMainListInfo", tmplValueMap{
		"UserMention": parts.m.Author.Mention(),
		"Symbols":     symbols,
	})
	parts.channelMessageSend(msg)
}

func (p *guildPresenter) cmdPieMainHelpHandler(parts *msgParts) {
	msg := msgFromTmpl("pieMainHelpUsage", tmplValueMap{
		"UserMention": parts.m.Author.Mention(),
		"BotPrefix":   piebotConfig.Discord.Prefix,
	})
	parts.channelMessageSend(msg)
}
