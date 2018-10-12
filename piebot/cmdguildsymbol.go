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

func (p *guildSymbolPresenter) cmdWithdrawHandler(parts *msgParts) {
	cmdPrefix := parts.prefix
	userID := parts.m.Author.ID
	username := parts.m.Author.Username
	userMention := parts.m.Author.Mention()
	sbl := p.symbol
	withdrawMinAmount, _ := amount.FromFloat64(p.coinInfo.Withdraw.Min)
	minTxFee, _ := amount.FromFloat64(p.coinInfo.Withdraw.TxFee)
	txFeePercent := p.coinInfo.Withdraw.TxFeePercent
	withdrawUsageInfo := &cmdWithdrawUserInfo{
		cmdUsageInfo: cmdUsageInfo{
			tmplName:        "withdrawUsage",
			IsShowUsageHint: true,
			CmdName:         "withdraw",
			UserMention:     userMention,
			Prefix:          string(cmdPrefix),
			Symbol:          string(sbl),
		},
		WithdrawMin:  withdrawMinAmount,
		TxFeePercent: txFeePercent * 100,
		TxFeeMin:     minTxFee,
	}
	cmdPartErrMsg := withdrawUsageInfo.String()
	if len(parts.contents) != 2 {
		parts.channelMessageSend(cmdPartErrMsg)
		return
	}

	withdrawAmount, err := strconv.ParseFloat(parts.contents[1], 64)
	if err != nil {
		parts.channelMessageSend(cmdPartErrMsg)
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
	if userAmount.CmpFloat(withdrawAmount) == -1 {
		msg := msgFromTmpl("withdrawAmountNotEnoughErr", tmplValueMap{
			"UserMention": userMention,
		})
		balInfoEmbed := balInfoEmbed(pieer, username, sbl)
		send := msgSend(msg, balInfoEmbed)
		parts.channelMessageSendComplex(send)
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
	send := msgSend(msg, balInfoEmbed)
	parts.channelMessageSendComplex(send)
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
	send := msgSend(content, embed)
	parts.channelMessageSendComplex(send)
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
	fields = append(fields, mef("Amount", fmt.Sprintf("%s %s", confirmed, sbl), true))
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
	cmdMsg := &cmdHelpUsageInfo{
		cmdUsageInfo: cmdUsageInfo{
			tmplName:        "helpUsage",
			IsShowUsageHint: false,
			CmdName:         "help",
			UserMention:     parts.m.Author.Mention(),
			Prefix:          string(cmdPrefix),
			Symbol:          string(sbl),
		},
		cmdPieUsageInfo: cmdPieUsageInfo{
			PieMin: pieMinAmount,
		},
		cmdWithdrawUserInfo: cmdWithdrawUserInfo{
			WithdrawMin:  withdrawMinAmount,
			TxFeePercent: txFeePercent * 100,
			TxFeeMin:     minTxFee,
		},
		IsManager: isManager,
	}
	parts.channelMessageSend(cmdMsg.String())
}

// @here @everyone isEveryone is true
func (p *guildSymbolPresenter) pieReceivers(s *discordgo.Session, guild *discordgo.Guild, channelID, pieUserID string, isEveryone bool, roles []string, users []*discordgo.User) ([]*discordgo.User, error) {
	receivers := []*discordgo.User{}
	rolesStr := strings.Join(roles, "|")
	gc := guildConfigs.gc(p.guildID)
	excludeRoles := strings.Join(gc.ExcludeRoles, "|")
	for _, member := range guild.Members {
		userID := member.User.ID
		switch {
		case member.User.Bot:
			fallthrough
		case member.User.ID == pieUserID:
			continue
		}
		isAdd := false
		for _, user := range users {
			if userID == user.ID {
				isAdd = true
				break
			}
		}
		if isAdd {
			receivers = append(receivers, member.User)
			continue
		}

		userPermission, err := s.State.UserChannelPermissions(userID, channelID)
		if err != nil {
			log.Errorf("PieReceivers get permission Error:%s[%s(%s)][%s(%s)]", err, member.User.Username, userID, guild.Name, guild.ID)
			continue
		}
		if (userPermission & discordgo.PermissionReadMessages) != discordgo.PermissionReadMessages {
			continue
		}
		isInExcludeRoles := false
		for _, role := range member.Roles {
			if strings.Contains(excludeRoles, role) {
				isInExcludeRoles = true
				break
			}
		}

		for _, role := range member.Roles {
			if strings.Contains(rolesStr, role) {
				isAdd = !isInExcludeRoles
				break
			}
		}
		if isAdd {
			receivers = append(receivers, member.User)
			continue
		}

		isOnline := false
		presence, err := s.State.Presence(guild.ID, userID)
		if err != nil {
			log.Errorf("PieReceivers get Presence Error:%s[%s(%s)][%s(%s)]", err, member.User.Username, userID, guild.Name, guild.ID)
			continue
		}
		if presence.Status == discordgo.StatusOnline || presence.Status == discordgo.StatusIdle {
			isOnline = true
		}

		if isEveryone && isOnline {
			isAdd = !isInExcludeRoles
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
	pieUsageInfo := &cmdPieUsageInfo{
		cmdUsageInfo: cmdUsageInfo{
			tmplName:        "pieUsage",
			IsShowUsageHint: true,
			CmdName:         "pie",
			UserMention:     userMention,
			Prefix:          string(cmdPrefix),
			Symbol:          string(sbl),
		},
		PieMin: pieMinAmount,
	}
	cmdUsage := pieUsageInfo.String()
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

	if sendAmount < pieMinAmount {
		msg := msgFromTmpl("pieAmountMinErr", tmplValueMap{
			"UserMention": userMention,
			"Min":         pieMinAmount,
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
	if userAmount.Cmp(sendAmount) == -1 {
		msg := msgFromTmpl("pieNotEnoughAmountErr", tmplValueMap{
			"UserMention": userMention,
			"Prefix":      cmdPrefix,
		})
		balInfoEmbed := balInfoEmbed(pieer, parts.m.Author.Username, sbl)
		send := msgSend(msg, balInfoEmbed)
		parts.channelMessageSendComplex(send)
		return
	}

	isEveryone := false
	if partLen == 1 {
		isEveryone = true
	} else if parts.m.MentionEveryone {
		isEveryone = true
	}

	receivers, err := p.pieReceivers(parts.s, parts.guild, parts.m.ChannelID, userID, isEveryone, parts.m.MentionRoles, parts.m.Mentions)
	if err != nil {
		log.Errorf("Pie get receivers error:%s", err)
		return
	}

	receiversLen := len(receivers)
	if receiversLen == 0 {
		msg := msgFromTmpl("pieNoPeopleErr", userMention)
		parts.channelMessageSend(msg)
		return
	}

	amountEach := sendAmount.DivFloat64(float64(receiversLen))

	if amountEach.Cmp(amount.Zero) == 0 {
		msg := msgFromTmpl("pieNotEnoughEachErr", tmplValueMap{
			"UserMention":   userMention,
			"SendAmount":    sendAmount,
			"Symbol":        sbl,
			"ReceiverCount": receiversLen,
		})
		parts.channelMessageSend(msg)
		return
	}

	err = p.dbSymbol.UserAmountSub(nil, userID, parts.m.Author.Username, sendAmount)
	if err != nil {
		log.Errorf("Pie modify sender amount error:%s", err)
		return
	}

	eachMsgReceiverNum := piebotConfig.Discord.EachPieMsgReceiversLimit
	receiversMap := make(map[int][]string)
	for i, receiver := range receivers {
		//msg index
		index := int(math.Floor(float64(i) / float64(eachMsgReceiverNum)))
		receiversMap[index] = append(receiversMap[index], receiver.Mention())
		err = p.dbSymbol.UserAmountAddUpsert(nil, receiver.ID, receiver.Username, amountEach)
		if err != nil {
			log.Errorf("Pie modify receiver amount error:%s", err)
		}
	}

	if receiversLen > eachMsgReceiverNum {
		sendCountMsg := msgFromTmpl("pieSendCountHint", tmplValueMap{
			"UserMention":   userMention,
			"Amount":        sendAmount,
			"Symbol":        sbl,
			"ReceiverCount": receiversLen,
		})
		parts.channelMessageSend(sendCountMsg)
	}

	for _, receivers := range receiversMap {
		msg := msgFromTmpl("pieSuccess", tmplValueMap{
			"CoinName":      coinConfig.Name,
			"AmountEach":    amountEach,
			"Symbol":        sbl,
			"Receivers":     receivers,
			"ReceiverCount": receiversLen,
			"ShowAllPeople": receiversLen > eachMsgReceiverNum,
		})
		parts.channelMessageSend(msg)
	}
	userName := parts.m.Author.Username
	userDmr := parts.m.Author.Discriminator
	log.Infof("[%s:pie]%s#%s(%s) send %s to %d peoples in [%s(%s)]", sbl, userName, userDmr, userMention, sendAmount, receiversLen, parts.guild.Name, parts.guild.ID)

}

func (p *guildSymbolPresenter) cmdChannelHandler(parts *msgParts) {
	cmdPrefix := parts.prefix
	userMention := parts.m.Author.Mention()
	sbl := p.symbol
	cmdUsage := &cmdUsageInfo{
		tmplName:        "channelUsage",
		IsShowUsageHint: true,
		CmdName:         "channel",
		UserMention:     userMention,
		Prefix:          string(cmdPrefix),
		Symbol:          string(sbl),
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
