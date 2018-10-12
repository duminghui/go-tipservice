// Package main provides ...
package main

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/duminghui/go-tipservice/amount"
	"github.com/duminghui/go-tipservice/db"
	"github.com/duminghui/go-util/utime"
)

func (p *guildSymbolPresenter) cmdPieAutoHandler(parts *msgParts) {
	symbol := string(p.symbol)
	userMention := parts.m.Author.Mention()
	prefix := parts.prefix
	contents := parts.contents
	min := p.coinInfo.Pie.Min
	usageTmplValue := tmplValueMap{
		"IsShowUsageHint": true,
		"CmdName":         "pieAutoAdd",
		"UserMention":     userMention,
		"Prefix":          prefix,
		"PieMin":          min,
		"Symbol":          symbol,
	}
	if len(contents) == 0 {
		msg := msgFromTmpl("pieAutoAddUsage", usageTmplValue)
		parts.channelMessageSend(msg)
		return
	}

	//1
	str := strings.Join(contents, "")
	exp := regexp.MustCompile(`<#(\d{18})>`)
	result := exp.FindAllStringSubmatch(str, -1)
	resultLen := len(result)
	if resultLen > 1 {
		msg := msgFromTmpl("pieAutoAddChannelErr", usageTmplValue)
		parts.channelMessageSend(msg)
		return
	}
	startIndex := 0
	mustContentLen := 6
	channel := parts.channel.ID
	if resultLen == 1 {
		channel = result[0][1]
		startIndex = 1
		mustContentLen = 7
	}

	if len(contents) < mustContentLen-1 {
		msg := msgFromTmpl("pieAutoAddParamsLenErr", usageTmplValue)
		parts.channelMessageSend(msg)
		return
	}

	//2
	m := parts.m
	roles := m.MentionRoles
	if len(roles) != 1 {
		msg := msgFromTmpl("pieAutoAddRoleErr", usageTmplValue)
		parts.channelMessageSend(msg)
		return
	}
	role := roles[0]

	//3
	status := contents[startIndex+1]
	if status != "online" && status != "all" {
		msg := msgFromTmpl("pieAutoAddStatusErr", usageTmplValue)
		parts.channelMessageSend(msg)
		return
	}
	isOnline := false
	if status == "online" {
		isOnline = true
	}

	//4
	interval := timeStrToDuration(contents[startIndex+2])
	if interval < 180*time.Second {
		msg := msgFromTmpl("pieAutoAddTimeErr", usageTmplValue)
		parts.channelMessageSend(msg)
		return
	}

	//5
	amt, err := amount.FromNumString(contents[startIndex+3])
	if err != nil {
		msg := msgFromTmpl("pieAutoAddAmountErr", usageTmplValue)
		parts.channelMessageSend(msg)
		return
	}
	if amt.Cmp(amount.Zero) <= 0 {
		msg := msgFromTmpl("pieAutoAddAmountErr", usageTmplValue)
		parts.channelMessageSend(msg)
		return
	}

	//6
	cycleTimesStr := contents[startIndex+4]
	cycleTimes, err := strconv.ParseInt(cycleTimesStr, 0, 64)
	if err != nil {
		msg := msgFromTmpl("pieAutoAddCycleTimeErr", usageTmplValue)
		parts.channelMessageSend(msg)
		return
	}
	if cycleTimes <= 0 {
		msg := msgFromTmpl("pieAutoAddCycleTimeErr", usageTmplValue)
		parts.channelMessageSend(msg)
		return
	}

	//7
	extendTime := timeStrToDuration(contents[startIndex+5])

	guildID := parts.guild.ID
	guildName := parts.guild.Name
	userID := parts.m.Author.ID
	userName := parts.m.Author.Username
	// func AutoPieAdd(guildID, guildName, userID, userName, symbol, channelID, roleID string, cycleTimes int64, intervalTime time.Duration, amount amount.Amount, isOnlineUser bool, extendTime time.Duration) error {
	ap, err := dbGuild.PieAutoAdd(guildID, guildName, userID, userName, symbol, channel, role, cycleTimes, interval, amt, isOnline, extendTime)
	if err != nil {
		msg := msgFromTmpl("pieAutoAddErr", nil)
		parts.channelMessageSend(msg)
		return
	}
	// pieAuto, err := db.PieAutoByID(nil, ap.ID.Hex())
	// if err != nil {
	// 	msg := msgFromTmpl("pieAutoAddErr", nil)
	// 	parts.channelMessageSend(msg)
	// 	return
	// }
	// msg := msgFromTmpl("pieAutoInfo", ap)
	// parts.channelMessageSend(msg)

	embed := pieAuto2Embed(parts.s, ap)
	content := msgFromTmpl("pieAutoAddSuccess", userMention)
	send := &discordgo.MessageSend{
		Content: content,
		Embed:   embed,
	}
	_, err = parts.s.ChannelMessageSendComplex(parts.m.ChannelID, send)
	if err != nil {
		log.Error(err)
	}
}

func pieAuto2Embed(s *discordgo.Session, p *db.PieAuto) *discordgo.MessageEmbed {

	fields := pieAuto2EmbedFields(s, p)
	createTime := fmt.Sprintf("CreateTime(UTC):%s", utime.FormatTimeStrUTC(p.CreateTime))
	coinConfig, ok := coinInfos[p.Symbol]
	if !ok {
		return nil
	}
	coinName := coinConfig.Name
	iconURL := coinConfig.IconURL
	website := coinConfig.Website
	title := fmt.Sprintf("PieAuto-Task#%s", p.ID.Hex())
	embed := &discordgo.MessageEmbed{
		Title: title,
		Color: 0xff000,
		Footer: &discordgo.MessageEmbedFooter{
			Text: createTime,
		},
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: iconURL,
		},
		Author: embedAuthor(coinName, website, ""),
		Fields: fields,
	}
	return embed
}

func pieAuto2EmbedFields(s *discordgo.Session, p *db.PieAuto) []*discordgo.MessageEmbedField {
	channel, err := channel(s, p.ChannelID)
	channelField := p.ChannelID
	if err == nil {
		channelField = channel.Name
	}
	guild, err := guild(s, p.GuildID)
	guildField := p.GuildID
	if err == nil {
		guildField = guild.Name
	}
	role, err := role(s, p.GuildID, p.RoleID)
	roleField := p.RoleID
	if err == nil {
		roleField = role.Name
	}
	cycleTimes := fmt.Sprintf("%d", p.CycleTimes)
	intervalTime := p.IntervalTime.String()
	peoples := "All Peoples"
	if p.IsOnlineUser {
		peoples = "Only Online Peoples"
	}
	amount := fmt.Sprintf("%s %s", p.Amount.String(), p.Symbol)
	runnedTimes := fmt.Sprintf("%d", p.RunnedTimes)
	nextPieTime := utime.FormatTimeStrUTC(p.NextPieTime)
	isEnd := fmt.Sprintf("%v", p.IsEnd)
	fields := make([]*discordgo.MessageEmbedField, 0)
	fields = append(fields, mef("Guild", guildField, true))
	fields = append(fields, mef("Channel", channelField, true))
	fields = append(fields, mef("Role", roleField, true))
	fields = append(fields, mef("Each Pie Amount", amount, true))
	fields = append(fields, mef("Cycle Times", cycleTimes, true))
	fields = append(fields, mef("Already Run Times", runnedTimes, true))
	fields = append(fields, mef("Interval", intervalTime, true))
	fields = append(fields, mef("Peoples", peoples, true))
	fields = append(fields, mef("Next Pie Time(UTC)", nextPieTime, true))
	fields = append(fields, mef("Task Finish", isEnd, true))

	return fields
}

func timeStrToDuration(s string) time.Duration {
	if strings.HasPrefix(s, "-") {
		return 0
	}
	td, err := time.ParseDuration(s)
	if err != nil {
		return 0
	}
	return td
}
