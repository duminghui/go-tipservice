// Package main provides ...
package main

import (
	"bytes"
	"fmt"
	"regexp"
	"strings"

	"github.com/bwmarrin/discordgo"
)

func cmdPieSet(s *discordgo.Session, m *discordgo.MessageCreate, msgParts []string) {
	if !isBotManager(s, m) {
		return
	}
	switch {
	case len(msgParts) == 0:
		cmdPieSetHelpHandler(s, m)
	case msgParts[0] == "list":
		cmdPieSetListHandler(s, m)
	case msgParts[0] == "prefix":
		cmdPieSetPrefixHandler(s, m, msgParts[1:])
	case msgParts[0] == "manager":
		cmdPieSetManagerHandler(s, m, msgParts[1:])
	case msgParts[0] == "info":
		cmdPieSetInfoHandler(s, m)
	default:
		cmdPieSetHelpHandler(s, m)
	}
}

func cmdPieSetInfoHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	channel, err := channel(s, m.ChannelID)
	if err != nil {
		log.Error("cmdPieSetInfoHander Error:", err)
		return
	}
	gcm, ok := guildConfigPresenters[channel.GuildID]
	if !ok {
		msg := fmt.Sprintf("%s no config for this server", m.Author.Mention())
		s.ChannelMessageSend(m.ChannelID, msg)
		return
	}
	buf := new(bytes.Buffer)
	buf.WriteString(m.Author.Mention())
	buf.WriteString(" Server Config:\n")
	symbolPrefix := gcm.symbolPrefixMap
	if symbolPrefix != nil {
		for k, v := range symbolPrefix {
			buf.WriteString("**")
			buf.WriteString(string(k))
			buf.WriteString("**\n  **Prefix: **")
			buf.WriteString(string(v))
			buf.WriteString("\n  **Active channels: **")
			coinConfig, err := guildConfigPresenters.guildCoinConfigBySymbol(channel.GuildID, k)
			if err != nil {
				log.Info("cmdPieSetInfoHandler Error:", err)
				buf.WriteString("\n")
				continue
			}
			channels := coinConfig.ChannelIDs
			if len(channels) == 0 {
				buf.WriteString("All channel")
			} else {
				for _, channelID := range channels {
					buf.WriteString(fmt.Sprintf("<#%s>", channelID))
				}
			}
			buf.WriteString("\n")
		}
	}
	buf.WriteString("**Manager**\n")
	guildManager := gcm.guildManager
	if guildManager != nil {
		buf.WriteString("  **User:**")
		if len(guildManager.Managers) > 0 {
			buf.WriteString("\n    - ")
		}
		managers := make([]string, 0, len(guildManager.Managers))
		for _, manager := range guildManager.Managers {
			member, _ := s.State.Member(channel.GuildID, manager)
			managers = append(managers, member.User.Username)
		}
		buf.WriteString(strings.Join(managers, ", "))

		buf.WriteString("\n  **Role:**")
		if len(guildManager.ManagerRoles) > 0 {
			buf.WriteString("\n    - ")
		}
		roles := make([]string, 0, len(guildManager.ManagerRoles))
		for _, roleID := range guildManager.ManagerRoles {
			role, _ := s.State.Role(channel.GuildID, roleID)
			roles = append(roles, role.Name)
		}
		buf.WriteString(strings.Join(roles, ", "))
	}
	s.ChannelMessageSend(m.ChannelID, buf.String())
}

func cmdPieSetManagerHandler(s *discordgo.Session, m *discordgo.MessageCreate, msgParts []string) {
	msg := fmt.Sprintf("%s manager command usage:\n  **manager <add|remove> <@user|@role>**", m.Author.Mention())
	if len(msgParts) < 2 {
		s.ChannelMessageSend(m.ChannelID, msg)
		return
	}
	operator := msgParts[0]
	if operator != "add" && operator != "remove" {
		s.ChannelMessageSend(m.ChannelID, msg)
		return
	}
	if len(m.MentionRoles) == 0 && len(m.Mentions) == 0 {
		s.ChannelMessageSend(m.ChannelID, msg)
		return
	}
	channel, err := channel(s, m.ChannelID)
	if err != nil {
		log.Error("cmdPieSetManagerHandler Error:", err)
		return
	}
	users := make([]string, 0, len(m.Mentions))
	for _, user := range m.Mentions {
		users = append(users, user.ID)
	}
	updateUsers, updateRoles, err := guildConfigPresenters.guildManagerUpdate(channel.GuildID, operator, users, m.MentionRoles)
	if err != nil {
		log.Error("cmdPieSetManagerHandler Error:", err)
		return
	}
	buf := new(bytes.Buffer)
	buf.WriteString(m.Author.Mention())
	buf.WriteString(" now server manager:")
	buf.WriteString("\n**Manager**\n")
	buf.WriteString("  **User:**")
	if len(updateUsers) > 0 {
		buf.WriteString("\n    - ")
	}
	managers := make([]string, 0, len(updateUsers))
	for _, manager := range updateUsers {
		member, _ := s.State.Member(channel.GuildID, manager)
		managers = append(managers, member.User.Username)
	}
	buf.WriteString(strings.Join(managers, ", "))

	buf.WriteString("\n  **Role:**")
	if len(updateRoles) > 0 {
		buf.WriteString("\n    - ")
	}
	roles := make([]string, 0, len(updateRoles))
	for _, roleID := range updateRoles {
		role, _ := s.State.Role(channel.GuildID, roleID)
		roles = append(roles, role.Name)
	}
	buf.WriteString(strings.Join(roles, ", "))
	s.ChannelMessageSend(m.ChannelID, buf.String())

}

func isBotManager(s *discordgo.Session, m *discordgo.MessageCreate) bool {
	userID := m.Author.ID
	if userID == piebotConfig.Discord.SuperManagerID {
		return true
	}
	channel, err := channel(s, m.ChannelID)
	if err != nil {
		log.Error("cmdPieSet Error:", err)
		return false
	}
	guild, err := guild(s, channel.GuildID)
	if err != nil {
		log.Error("cmdPieSet Error:", err)
		return false
	}
	if userID == guild.OwnerID {
		return true
	}
	guildConfigMge, ok := guildConfigPresenters[guild.ID]
	if !ok {
		return false
	}
	guildManager := guildConfigMge.guildManager
	if guildManager == nil {
		return false
	}
	member, err := s.State.Member(channel.GuildID, userID)
	if err != nil {
		log.Error("isBotManager Error:", err)
		return false
	}
	if guildManager.IsManager(userID) || guildManager.InManagerRoles(member.Roles) {
		return true
	}
	return false
}

func cmdPieSetPrefixHandler(s *discordgo.Session, m *discordgo.MessageCreate, msgParts []string) {
	if len(msgParts) != 2 {
		msg := fmt.Sprintf("%s prefix command usage:\n  **?pieconfig prefix <symbol> <prefix>**", m.Author.Mention())
		s.ChannelMessageSend(m.ChannelID, msg)
		return
	}
	symbol := symbolWrap(msgParts[0])
	prefix := prefixWrap(msgParts[1])
	if prefix == "?" {
		msg := fmt.Sprintf("%s prefix `?` is used by bot", m.Author.Mention())
		s.ChannelMessageSend(m.ChannelID, msg)
		return
	}
	if !hasSymbol(symbol) {
		msg := fmt.Sprintf("%s don't have this coin's symbol `%s`", m.Author.Mention(), symbol)
		s.ChannelMessageSend(m.ChannelID, msg)
		return
	}
	// oldPrefix := ""
	channel, err := channel(s, m.ChannelID)
	if err != nil {
		log.Errorf("cmdPieSetPrefix error:%s", err)
		return
	}
	symbolTmp, _ := guildConfigPresenters.symbolByPrefix(channel.GuildID, prefix)

	if symbolTmp != "" {
		msg := fmt.Sprintf("%s command prefix `%s` is config to `%s`", m.Author.Mention(), prefix, symbolTmp)
		s.ChannelMessageSend(m.ChannelID, msg)
		return
	}

	err = guildConfigPresenters.updatePrefix(channel.GuildID, symbol, prefix)
	if err != nil {
		log.Errorf("cmdPieSetPrefix Cache error:%s", err)
		return
	}
	msg := fmt.Sprintf("%s `%s` set command prefix `%s`", m.Author.Mention(), symbol, prefix)
	s.ChannelMessageSend(m.ChannelID, msg)

}

func cmdPieSetListHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	buf := new(bytes.Buffer)
	buf.WriteString(m.Author.Mention())
	buf.WriteString(" All configureable coin's symbols:\n")
	for k := range coinPresenters {
		buf.WriteString(" -- **")
		buf.WriteString(string(k))
		buf.WriteString("**\n")
	}
	s.ChannelMessageSend(m.ChannelID, buf.String())
}

func cmdPieSetHelpHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	msg := fmt.Sprintf("%s You can use these subcommands:\n"+
		"  **info**\n"+
		"    --show config info\n"+
		"  **list**\n"+
		"    --List all configurable coin's symbols\n"+
		"  **prefix <symbol> <prefix>**\n"+
		"    -- config <symbol>'s command prefix\n"+
		"  **manager <add|remove> <@user|@role>**\n"+
		"    -- add or remove manager for PieBot",
		m.Author.Mention())
	s.ChannelMessageSend(m.ChannelID, msg)
}

func cmdChannelHandler(s *discordgo.Session, m *discordgo.MessageCreate, msgParts *msgParts) {
	if !isBotManager(s, m) {
		return
	}
	cmdPrefix := msgParts.prefix
	userMention := m.Author.Mention()
	channelID := m.ChannelID
	channel, err := channel(s, channelID)
	if err != nil {
		log.Error("cmdSetChannelHandler Error:", err)
		return
	}
	symbol, err := guildConfigPresenters.symbolByPrefix(channel.GuildID, cmdPrefix)
	if err != nil {
		log.Error("cmdSetChannelHandler Error:", err)
		return
	}
	usage := fmt.Sprintf(msgParts.cmdInfo.usage, symbol)
	msg := fmt.Sprintf("%s %s%s command usage:\n%s", userMention, cmdPrefix, msgParts.cmdInfo.name, usage)
	parts := msgParts.parts
	if len(parts) < 2 {
		s.ChannelMessageSend(channelID, msg)
		return
	}
	operator := parts[0]
	if operator != "add" && operator != "remove" {
		s.ChannelMessageSend(channelID, msg)
		return
	}
	str := strings.Join(parts[1:], "")
	exp := regexp.MustCompile(`<#(\d{18})>`)
	result := exp.FindAllStringSubmatch(str, -1)
	channels := make([]string, 0, len(result))
	for _, v := range result {
		channels = append(channels, v[1])
	}
	finalChannels, err := guildConfigPresenters.guildChannelUpdate(channel.GuildID, symbol, operator, channels)
	if err != nil {
		log.Error("cmdSetChannelHandler Error:", err)
		return
	}
	if operator == "add" {
		msg = fmt.Sprintf("%s Add Success.\n", userMention)
	} else {
		msg = fmt.Sprintf("%s Remove Success.\n", userMention)
	}
	buf := new(bytes.Buffer)
	buf.WriteString(msg)
	if len(finalChannels) == 0 {
		buf.WriteString(fmt.Sprintf(" `%s`'s commands now active all channels", symbol))
	} else {
		buf.WriteString(fmt.Sprintf(" `%s`' commands  now active in these channel:\n", symbol))
		for _, channelID := range finalChannels {
			buf.WriteString(fmt.Sprintf("<#%s>", channelID))
		}
	}
	s.ChannelMessageSend(m.ChannelID, buf.String())

}
