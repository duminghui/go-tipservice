// Package main provides ...
package main

import (
	"fmt"
	"math"
	"regexp"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/duminghui/go-tipservice/amount"
	"github.com/duminghui/go-tipservice/db"
)

func registerSymbolCmds() {
	registerSymbolCmd("help", true, false, false, (*guildSymbolPresenter).cmdPieHelperHandler)
	registerSymbolCmd("pie", true, false, false, (*guildSymbolPresenter).cmdPieHandler)
	registerSymbolCmd("deposit", true, false, false, (*guildSymbolPresenter).cmdDepositHandler)
	registerSymbolCmd("bal", true, false, false, (*guildSymbolPresenter).cmdBalHandler)
	registerSymbolCmd("withdraw", true, false, false, (*guildSymbolPresenter).cmdWithdrawHandler)
	registerSymbolCmd("channel", false, true, false, cmdChannelHandler)
}

func (p *guildSymbolPresenter) cmdWithdrawHandler(parts *msgParts) {
	cmdPrefix := parts.prefix
	userID := parts.m.Author.ID
	username := parts.m.Author.Username
	userMention := parts.m.Author.Mention()
	sbl := p.symbol
	withdrawMinAmount, _ := amount.FromFloat64(p.coinInfo.Withdraw.Min)
	minTxFee, _ := amount.FromFloat64(p.coinInfo.Withdraw.TxFee)
	txFeePercent := p.coinInfo.Withdraw.TxFeePercent
	tmplValue := &tmplValueMap{
		"IsShowUsageHint": true,
		"CmdName":         "withdraw",
		"UserMention":     userMention,
		"Prefix":          string(cmdPrefix),
		"Symbol":          string(sbl),
		"WithdrawMin":     withdrawMinAmount,
		"TxFeePercent":    txFeePercent * 100,
		"TxFeeMin":        minTxFee,
	}
	cmdUsage := msgFromTmpl("withdrawUsage", tmplValue)
	if len(parts.contents) != 2 {
		parts.channelMessageSend(cmdUsage)
		return
	}

	withdrawAmount, err := strconv.ParseFloat(parts.contents[1], 64)
	if err != nil {
		parts.channelMessageSend(cmdUsage)
		return
	}

	if withdrawAmount < withdrawMinAmount.Float64() {
		msg := msgFromTmpl("withdrawMinAmountErr", tmplValueMap{
			"UserMention": userMention,
			"Min":         withdrawMinAmount,
			"Symbol":      sbl,
		})
		parts.channelMessageSend(msg)
		return
	}

	address := parts.contents[0]
	validateAddress, err := p.rpc.ValidateAddress(address)
	if err != nil {
		log.Error("[CMD]withdraw ValidateAddress Error:", err)
		msg := msgFromTmpl("walletMaintenance", userMention)
		parts.channelMessageSend(msg)
		return
	}
	if !validateAddress.IsValid {
		msg := msgFromTmpl("withdrawValidateAddrErr", tmplValueMap{
			"UserMention": userMention,
			"Addr":        address,
			"Symbol":      sbl,
		})
		parts.channelMessageSend(msg)
		return
	}
	if validateAddress.IsMine {
		msg := msgFromTmpl("withdrawBotAddrErr", tmplValueMap{
			"UserMention": userMention,
			"Addr":        address,
			"Prefix":      cmdPrefix,
			"Symbol":      sbl,
		})
		parts.channelMessageSend(msg)
		return
	}
	pieer, err := p.dbSymbol.UserByID(nil, userID)
	if err != nil {
		log.Errorf("[CMD]pie UserByID Error:%s", err)
		return
	}
	userAmount := amount.Zero
	if pieer != nil {
		userAmount = pieer.Amount
	}
	if userAmount.CmpFloat(withdrawAmount) < 0 {
		msg := msgFromTmpl("withdrawAmountNotEnoughErr", tmplValueMap{
			"UserMention": userMention,
		})
		balInfoEmbed := balInfoEmbed(pieer, username, sbl)
		parts.channelMessageSendComplex(msg, balInfoEmbed)
		return
	}
	txfee, _ := amount.FromFloat64(withdrawAmount * txFeePercent)

	if txfee.Cmp(minTxFee) == -1 {
		txfee = minTxFee
	}
	withdrawAmountProxy, _ := amount.FromFloat64(withdrawAmount)
	finalWithdrawAmount := withdrawAmountProxy.Sub(txfee)

	withdrawTxID, err := p.rpc.SendToAddress(address, finalWithdrawAmount.Float64())
	if err != nil {
		msg := msgFromTmpl("walletMaintenance", userMention)
		parts.channelMessageSend(msg)
		return
	}

	err = p.dbSymbol.UserAmountSub(nil, userID, username, withdrawAmountProxy)
	if err != nil {
		log.Errorf("[%s] Withdraw Amount Update Error:%s[%s][%s][%s][%.8f]", sbl, err, userID, username, withdrawTxID, withdrawAmount)
	}
	p.dbSymbol.SaveWithdraw(userID, username, address, withdrawTxID, withdrawAmount)
	pieer, err = p.dbSymbol.UserByID(nil, userID)
	if err != nil {
		log.Errorf("[CMD]pie UserByID Error:%s", err)
		return
	}
	msg := msgFromTmpl("withdrawSuccess", tmplValueMap{
		"UserMention": userMention,
		"Amount":      withdrawAmountProxy,
		"Symbol":      sbl,
		"Addr":        address,
		"TxFee":       txfee,
		"TxExpUrl":    p.coinInfo.TxExplorerURL,
		"TxID":        withdrawTxID,
	})
	balInfoEmbed := balInfoEmbed(pieer, parts.m.Author.Username, sbl)
	parts.channelMessageSendComplex(msg, balInfoEmbed)
}

func (p *guildSymbolPresenter) cmdBalHandler(parts *msgParts) {
	userID := parts.m.Author.ID
	sbl := p.symbol
	user, err := p.dbSymbol.UserByID(nil, userID)
	if err != nil {
		log.Errorf("[%s] Deposit UserByID Error:%s", sbl, err)
		return
	}
	embed := balInfoEmbed(user, parts.m.Author.Username, sbl)
	content := msgFromTmpl("balAmount", parts.m.Author.Mention())
	parts.channelMessageSendComplex(content, embed)
	if err != nil {
		log.Error(err)
	}
}

func balInfoEmbed(user *db.User, username string, sbl symbol) *discordgo.MessageEmbed {
	confirmed := amount.Zero
	unconfirmed := amount.Zero
	if user != nil {
		confirmed = user.Amount
		unconfirmed = user.UnconfirmedAmount
	}
	coinConfig := coinInfos[string(sbl)]
	iconURL := coinConfig.IconURL
	coinName := coinConfig.Name
	website := coinConfig.Website
	authorInfo := fmt.Sprintf("%s' Balance", username)
	fields := make([]*discordgo.MessageEmbedField, 0, 2)
	fields = append(fields, mef("Amount", fmt.Sprintf("%s %s", confirmed, sbl), false))
	fields = append(fields, mef("Unconfirmed Amount", fmt.Sprintf("%s %s", unconfirmed, sbl), true))
	embed := &discordgo.MessageEmbed{
		Title: authorInfo,
		Color: 0xff000,
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: iconURL,
		},
		Author: &discordgo.MessageEmbedAuthor{
			URL:  website,
			Name: coinName,
		},
		Fields: fields,
	}
	return embed

}
func (p *guildSymbolPresenter) cmdDepositHandler(parts *msgParts) {
	userID := parts.m.Author.ID
	userMention := parts.m.Author.Mention()
	sbl := p.symbol
	user, err := p.dbSymbol.UserByID(nil, userID)
	if err != nil {
		log.Errorf("[%s] Deposit UserByID Error:%s", sbl, err)
		return
	}
	if user != nil && user.Address != "" {
		msg := msgFromTmpl("depositInfo", tmplValueMap{
			"UserMention": userMention,
			"Symbol":      sbl,
			"Addr":        user.Address,
		})
		parts.channelMessageSend(msg)
		return
	}
	address, err := p.rpc.GetNewAddress(userID)
	if err != nil {
		msg := msgFromTmpl("walletMaintenance", userMention)
		parts.channelMessageSend(msg)
		return
	}
	err = p.dbSymbol.UserAddressUpsert(userID, parts.m.Author.Username, address, user == nil)
	if err != nil {
		log.Errorf("[%s] Deposit UserAddressUpsert Error:%s", sbl, err)
		return
	}
	msg := msgFromTmpl("depositInfo", tmplValueMap{
		"UserMention": userMention,
		"Symbol":      sbl,
		"Addr":        address,
	})
	parts.channelMessageSend(msg)
}

func (p *guildSymbolPresenter) cmdPieHelperHandler(parts *msgParts) {
	cmdPrefix := parts.prefix
	sbl := p.symbol
	coinConfig := p.coinInfo
	isManager := parts.isManager
	pieMinAmount, _ := amount.FromFloat64(coinConfig.Pie.Min)
	withdrawMinAmount, _ := amount.FromFloat64(coinConfig.Withdraw.Min)
	minTxFee, _ := amount.FromFloat64(coinConfig.Withdraw.TxFee)
	txFeePercent := coinConfig.Withdraw.TxFeePercent
	tmplValue := &tmplValueMap{
		"IsShowUsageHint": false,
		"CmdName":         "help",
		"UserMention":     parts.m.Author.Mention(),
		"Prefix":          cmdPrefix,
		"Symbol":          sbl,
		"PieMin":          pieMinAmount,
		"WithdrawMin":     withdrawMinAmount,
		"TxFeePercent":    txFeePercent * 100,
		"TxFeeMin":        minTxFee,
		"IsManager":       isManager,
	}
	helpInfo := msgFromTmpl("helpUsage", tmplValue)
	parts.channelMessageSend(helpInfo)
}

type pieReceiverGeneratorSymbol struct {
	s            *discordgo.Session
	guild        *discordgo.Guild
	channelID    string
	pieerID      string
	isEveryone   bool
	mentionRoles []string
	mentions     []*discordgo.User
}

func (r *pieReceiverGeneratorSymbol) Receivers() ([]*discordgo.User, error) {
	receivers := make([]*discordgo.User, 0)
	roles := r.mentionRoles
	rolesStr := strings.Join(roles, "|")
	guild := r.guild
	guildID := guild.ID
	guildName := guild.Name
	channelID := r.channelID
	gc := guildConfigs.gc(guildID)
	excludeRoles := strings.Join(gc.ExcludeRoles, "|")
	s := r.s
	for _, member := range guild.Members {
		userID := member.User.ID
		switch {
		case member.User.Bot:
			fallthrough
		case userID == r.pieerID:
			continue
		}
		isInUsers := false
		for _, user := range r.mentions {
			if userID == user.ID {
				isInUsers = true
				break
			}
		}
		if isInUsers {
			receivers = append(receivers, member.User)
			continue
		}

		userPermission, err := userChannelPermissions(s, userID, channelID)
		if err != nil {
			log.Errorf("PieReceivers get permission Error:%s[%s(%s)][%s(%s)]", err, member.User.Username, userID, guildName, guildID)
			continue
		}
		if (userPermission & discordgo.PermissionReadMessages) != discordgo.PermissionReadMessages {
			continue
		}
		isInExcludeRoles := false
		isInRoles := false
		for _, role := range member.Roles {
			if strings.Contains(excludeRoles, role) {
				isInExcludeRoles = true
				break
			}
			if strings.Contains(rolesStr, role) {
				isInRoles = true
			}
		}
		if isInExcludeRoles {
			continue
		}

		if isInRoles {
			receivers = append(receivers, member.User)
			continue
		}

		isOnline := false
		presence, err := presence(s, guildID, userID)
		if err != nil {
			log.Errorf("PieReceivers get Presence Error:%s[%s(%s)][%s(%s)]", err, member.User.Username, userID, guild.Name, guild.ID)
			continue
		}
		if presence.Status == discordgo.StatusOnline || presence.Status == discordgo.StatusIdle {
			isOnline = true
		}
		isAdd := false
		if r.isEveryone && isOnline {
			isAdd = true
		}
		if isAdd {
			receivers = append(receivers, member.User)
		}
	}
	return receivers, nil
}

func (p *guildSymbolPresenter) cmdPieHandler(parts *msgParts) {
	userID := parts.m.Author.ID
	userMention := parts.m.Author.Mention()
	cmdPrefix := parts.prefix
	sbl := p.symbol
	coinConfig := p.coinInfo
	pieMinAmount, err := amount.FromFloat64(coinConfig.Pie.Min)
	tmplValue := &tmplValueMap{
		"IsShowUsageHint": true,
		"CmdName":         "pie",
		"UserMention":     userMention,
		"Prefix":          string(cmdPrefix),
		"Symbol":          string(sbl),
		"PieMin":          pieMinAmount,
	}
	cmdUsage := msgFromTmpl("pieUsage", tmplValue)
	partLen := len(parts.contents)
	if partLen == 0 {
		parts.channelMessageSend(cmdUsage)
		return
	}

	sendAmount, err := amount.FromNumString(parts.contents[partLen-1])
	if err != nil {
		parts.channelMessageSend(cmdUsage)
		return
	}

	isEveryone := false
	if partLen == 1 {
		isEveryone = true
	} else if parts.m.MentionEveryone {
		isEveryone = true
	}
	receiverGenerator := &pieReceiverGeneratorSymbol{
		s:            parts.s,
		guild:        parts.guild,
		channelID:    parts.channel.ID,
		pieerID:      userID,
		isEveryone:   isEveryone,
		mentionRoles: parts.m.MentionRoles,
		mentions:     parts.m.Mentions,
	}
	pie := &pie{
		symbol:            sbl,
		userID:            userID,
		userName:          parts.m.Author.Username,
		amount:            sendAmount,
		receiverGenerator: receiverGenerator,
	}

	report, err := pie.pie()
	log.Infof("%#v", report)
	receiverCount := 0
	if report != nil {
		receiverCount = report.receiverCount
	}
	tmplValue = &tmplValueMap{
		"UserMention":   userMention,
		"Min":           pieMinAmount,
		"Symbol":        sbl,
		"Prefix":        cmdPrefix,
		"SendAmount":    sendAmount,
		"Amount":        sendAmount,
		"ReceiverCount": receiverCount,
	}
	if err != nil {
		switch err {
		case errPieUserNotExists,
			errPieNoSymbol:
			return
		case errPieAmountMin:
			msg := msgFromTmpl("pieAmountMinErr", tmplValue)
			parts.channelMessageSend(msg)
			return
		case errPieNotEnoughAmount:
			msg := msgFromTmpl("pieNotEnoughAmountErr", tmplValue)
			balInfoEmbed := balInfoEmbed(report.pieer, parts.m.Author.Username, sbl)
			parts.channelMessageSendComplex(msg, balInfoEmbed)
		case errPieNoReceiver:
			msg := msgFromTmpl("pieNoPeopleErr", userMention)
			parts.channelMessageSend(msg)
			return
		case errPieNotEnoughEachAmount:
			msg := msgFromTmpl("pieNotEnoughEachErr", tmplValue)
			parts.channelMessageSend(msg)
			return
		}
		log.Error("Pie Error:", err)
		return
	}

	eachMsgReceiverNum := piebotConfig.Discord.EachPieMsgReceiversLimit
	receivers := report.receivers
	receiversMap := make(map[int][]string)
	for i, receiver := range receivers {
		//msg index
		index := int(math.Floor(float64(i) / float64(eachMsgReceiverNum)))
		receiversMap[index] = append(receiversMap[index], receiver.Mention())
	}

	if receiverCount > eachMsgReceiverNum {
		sendCountMsg := msgFromTmpl("pieSendCountHint", tmplValue)
		parts.channelMessageSend(sendCountMsg)
	}

	for _, receivers := range receiversMap {
		msg := msgFromTmpl("pieSuccess", tmplValueMap{
			"CoinName":      coinConfig.Name,
			"AmountEach":    report.eachAmount,
			"Symbol":        sbl,
			"Receivers":     receivers,
			"ReceiverCount": receiverCount,
			"ShowAllPeople": receiverCount > eachMsgReceiverNum,
		})
		parts.channelMessageSend(msg)
	}
	userName := parts.m.Author.Username
	userDmr := parts.m.Author.Discriminator
	log.Infof("[%s:pie]%s#%s(%s) send %s to %d peoples in [%s(%s)]", sbl, userName, userDmr, userID, sendAmount, receiverCount, parts.guild.Name, parts.guild.ID)
}

var cmdChannelHandler = (*guildSymbolPresenter).cmdChannelHandler

func (p *guildSymbolPresenter) cmdChannelHandler(parts *msgParts) {
	cmdPrefix := parts.prefix
	userMention := parts.m.Author.Mention()
	sbl := p.symbol
	tmplValue := &tmplValueMap{
		"IsShowUsageHint": true,
		"CmdName":         "channel",
		"UserMention":     userMention,
		"Prefix":          string(cmdPrefix),
		"Symbol":          string(sbl),
	}
	cmdUsage := msgFromTmpl("channelUsage", tmplValue)
	contents := parts.contents
	if len(contents) < 2 {
		parts.channelMessageSend(cmdUsage)
		return
	}
	operator := contents[0]
	if operator != "add" && operator != "remove" {
		parts.channelMessageSend(cmdUsage)
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
		parts.channelMessageSend(cmdUsage)
		return
	}
	finalChannels, err := p.guildChannelUpdate(sbl, operator, channels)
	if err != nil {
		log.Error("cmdSetChannelHandler Error:", err)
		return
	}
	msg := msgFromTmpl("channelOperatorSuccess", tmplValueMap{
		"UserMention": userMention,
		"Operator":    operator,
		"Symbol":      sbl,
		"Channels":    finalChannels,
	})
	parts.channelMessageSend(msg)
}
