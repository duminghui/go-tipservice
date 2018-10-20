// Package main provides ...
package main

import (
	"flag"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/duminghui/go-tipservice/amount"
	"github.com/duminghui/go-tipservice/db"
	"github.com/duminghui/go-util/utime"
)

func registerSymbolPieAutoCmd() {
	registerSymbolCmd("pieAuto", false, true, false, (*guildSymbolPresenter).cmdPieAutoHandler)
}

func (p *guildSymbolPresenter) cmdPieAutoHandler(parts *msgParts) {
	sbl := p.symbol
	userMention := parts.m.Author.Mention()
	prefix := parts.prefix
	contents := parts.contents
	min := p.coinInfo.Pie.Min
	amountMin, _ := amount.FromFloat64(min)
	usageTmplValue := tmplValueMap{
		"IsShowUsageHint": true,
		"CmdName":         "pieAuto",
		"UserMention":     userMention,
		"Prefix":          prefix,
		"PieMin":          min,
		"Symbol":          sbl,
	}
	msgUsage := msgFromTmpl("pieAutoUsage", usageTmplValue)
	if len(contents) == 0 {
		parts.channelMessageSend(msgUsage)
		return
	}

	m := parts.m
	roles := m.MentionRoles
	role := ""
	roleCount := len(roles)
	if roleCount > 1 {
		msg := msgFromTmpl("pieAutoRoleErr", usageTmplValue)
		parts.channelMessageSend(msg)
		return
	}
	if roleCount == 1 {
		role = roles[0]
	}

	cmdFlag := flag.NewFlagSet("pieAuto", flag.ContinueOnError)
	statusField := cmdFlag.String("s", "online", "")
	intervalField := cmdFlag.String("i", "180s", "")
	amountField := cmdFlag.Float64("a", 0, "")
	cycleTimesField := cmdFlag.Int64("c", 10, "")
	extendTimeField := cmdFlag.String("e", "30s", "")
	err := cmdFlag.Parse(parts.contents)
	if err != nil {
		log.Errorf("PieAutoUsage parse Error:%s", err)
		parts.channelMessageSend(msgUsage)
		return
	}

	if cmdFlag.NArg() > 2 {
		log.Errorf("PieAutoUsage parse Error:%s", err)
		parts.channelMessageSend(msgUsage)
		return
	}

	//1
	str := strings.Join(contents, "")
	exp := regexp.MustCompile(`<#(\d{18})>`)
	result := exp.FindAllStringSubmatch(str, -1)
	resultLen := len(result)
	if resultLen > 1 {
		msg := msgFromTmpl("pieAutoChannelErr", usageTmplValue)
		parts.channelMessageSend(msg)
		return
	}

	channelID := parts.channel.ID
	if resultLen == 1 {
		channelID = result[0][1]
	}

	if *statusField != "online" && *statusField != "all" {
		msg := msgFromTmpl("pieAutoStatusErr", usageTmplValue)
		parts.channelMessageSend(msg)
		return
	}
	isOnline := false
	if *statusField == "online" {
		isOnline = true
	}

	//4
	interval := timeStrToDuration(*intervalField)
	if interval < 180*time.Second {
		msg := msgFromTmpl("pieAutoIntervalErr", usageTmplValue)
		parts.channelMessageSend(msg)
		return
	}

	//6
	if *cycleTimesField <= 0 {
		msg := msgFromTmpl("pieAutoCycleTimeErr", usageTmplValue)
		parts.channelMessageSend(msg)
		return
	}

	//5
	sendAmount, err := amount.FromFloat64(*amountField)
	if err != nil {
		msg := msgFromTmpl("pieAutoAmountErr", usageTmplValue)
		parts.channelMessageSend(msg)
		return
	}

	if sendAmount.Cmp(amountMin) < 0 {
		msg := msgFromTmpl("pieAutoAmountMinErr", usageTmplValue)
		parts.channelMessageSend(msg)
		return
	}

	userID := parts.m.Author.ID

	pieer, err := p.dbSymbol.UserByID(nil, userID)
	if err != nil {
		log.Errorf("[CMD]pie UserByID Error:%s", err)
		return
	}
	userAmount := amount.Zero
	if pieer != nil {
		userAmount = pieer.Amount
	}
	cycleTimesTmp, _ := amount.FromInt64(*cycleTimesField)
	sumSendAmount := sendAmount.Mul(cycleTimesTmp)
	userName := parts.m.Author.Username
	if userAmount.Cmp(sumSendAmount) < 0 {
		msg := msgFromTmpl("pieAutoAmountNotEnoughErr", tmplValueMap{
			"UserMention": userMention,
			"SumAmount":   sumSendAmount,
			"Symbol":      sbl,
		})
		balInfoEmbed := balInfoEmbed(pieer, userName, sbl)
		parts.channelMessageSendComplex(msg, balInfoEmbed)
		return
	}
	//7
	extendTime := timeStrToDuration(*extendTimeField)
	if extendTime < 30*time.Second {
		msg := msgFromTmpl("pieAutoExtendTimeErr", usageTmplValue)
		parts.channelMessageSend(msg)
		return
	}

	guildID := parts.guild.ID
	guildName := parts.guild.Name
	// func AutoPieAdd(guildID, guildName, userID, userName, symbol, channelID, roleID string, cycleTimes int64, intervalTime time.Duration, amount amount.Amount, isOnlineUser bool, extendTime time.Duration) error {

	targetChannel, err := channel(parts.s, channelID)
	targetChannelName := ""
	if err == nil {
		targetChannelName = targetChannel.Name
	}
	ap, err := dbGuild.PieAutoAdd(guildID, guildName, userID, userName, string(sbl), channelID, targetChannelName, role, *cycleTimesField, interval, sendAmount, isOnline, extendTime)
	if err != nil {
		msg := msgFromTmpl("pieAutoErr", nil)
		parts.channelMessageSend(msg)
		return
	}
	ap, err = dbGuild.PieAutoByID(nil, ap.ID.Hex())
	if err != nil {
		msg := msgFromTmpl("pieAutoErr", nil)
		parts.channelMessageSend(msg)
		return
	}

	embed := pieAuto2Embed(parts.s, ap)
	content := msgFromTmpl("pieAutoSuccess", userMention)
	channelMsg, err := parts.channelMessageSendComplex(content, embed)
	if err != nil {
		log.Error(err)
	}
	dmContent := msgFromTmpl("pieAutoSuccessDM", userMention)
	dmMsg, err := directMessageComplx(parts.s, userID, dmContent, embed)
	if err != nil {
		log.Error(err)
	}
	p.dbGuild.PieAutoMsgAdd(dmMsg.ID, userID, ap.ID.Hex())
	err = parts.s.MessageReactionAdd(dmMsg.ChannelID, dmMsg.ID, reactionStop)
	if err != nil {
		log.Errorf("MessageReactionAddError:%s", err)
	}
	err = parts.s.MessageReactionAdd(dmMsg.ChannelID, dmMsg.ID, reactionRefresh)
	if err != nil {
		log.Errorf("MessageReactionAddError:%s", err)
	}
	p.dbGuild.PieAutoMsgAdd(channelMsg.ID, userID, ap.ID.Hex())
	err = parts.s.MessageReactionAdd(channelMsg.ChannelID, channelMsg.ID, reactionStop)
	err = parts.s.MessageReactionAdd(channelMsg.ChannelID, channelMsg.ID, reactionRefresh)
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
	desc := ""
	if p.IsEnd {
		desc = "This task is finished"
	}
	embed := &discordgo.MessageEmbed{
		Title:       title,
		Description: desc,
		Color:       0x00ff00,
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
	roleField := "All Roles"
	if p.RoleID != "" {
		roleField = p.RoleID
		role, err := role(s, p.GuildID, p.RoleID)
		if err == nil {
			roleField = fmt.Sprintf("@%s", role.Name)
		}
	}
	cycleTimes := fmt.Sprintf("%d", p.CycleTimes)
	intervalTime := p.IntervalTime.String()
	peoples := "All Peoples"
	if p.IsOnlineUser {
		peoples = "Only Online Peoples"
	}
	cycleTimesAmount, _ := amount.FromInt64(p.CycleTimes)
	sendAmount := p.Amount.Mul(cycleTimesAmount)
	sendAmountField := fmt.Sprintf("%s %s", sendAmount, p.Symbol)
	amount := fmt.Sprintf("%s %s", p.Amount, p.Symbol)
	runnedTimes := fmt.Sprintf("%d", p.RunnedTimes)
	nextPieTime := utime.FormatTimeStrUTC(p.NextPieTime)
	isEnd := fmt.Sprintf("%v", p.IsEnd)
	fields := make([]*discordgo.MessageEmbedField, 0)
	fields = append(fields, mef("Guild", guildField, true))
	fields = append(fields, mef("Channel", channelField, true))
	fields = append(fields, mef("Role", roleField, true))
	fields = append(fields, mef("Cycle Times", cycleTimes, true))
	fields = append(fields, mef("Already Run Times", runnedTimes, true))
	fields = append(fields, mef("Interval", intervalTime, true))
	fields = append(fields, mef("Each Pie Amount", amount, true))
	fields = append(fields, mef("Total Amount", sendAmountField, true))
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
