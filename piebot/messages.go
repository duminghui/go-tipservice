// Package main provides ...
package main

import (
	"math"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/duminghui/go-tipservice/amount"
)

type cmdHandler func(*guildConfigPresenter, *msgParts)

// type cmdHandler func(*discordgo.Session, *discordgo.MessageCreate, *msgParts)

type msgParts struct {
	s         *discordgo.Session
	m         *discordgo.MessageCreate
	isManager bool
	channel   *discordgo.Channel
	guild     *discordgo.Guild
	prefix    prefixWrap
	symbol    symbolWrap
	contents  []string
}

func (p *msgParts) channelMessageSend(msg string) {
	p.s.ChannelMessageSend(p.channel.ID, msg)
}

type cmdInfo struct {
	name         string
	managerCmd   bool
	channelLimit bool
	handler      cmdHandler
}

var cmdInfoMap = make(map[string]*cmdInfo)

var cmdChannel = "channel"

func reigsterBotCmdHandler() {
	help := &cmdInfo{
		name:         "help",
		channelLimit: true,
		handler:      (*guildConfigPresenter).cmdPieHelperHandler,
	}
	cmdInfoMap[help.name] = help
	pie := &cmdInfo{
		name:         "pie",
		channelLimit: true,
		handler:      (*guildConfigPresenter).cmdPieHandler,
	}
	cmdInfoMap[pie.name] = pie
	deposit := &cmdInfo{
		name:         "deposit",
		channelLimit: true,
		handler:      (*guildConfigPresenter).cmdDepositHandler,
	}
	cmdInfoMap[deposit.name] = deposit
	bal := &cmdInfo{
		name:         "bal",
		channelLimit: true,
		handler:      (*guildConfigPresenter).cmdBalHandler,
	}
	cmdInfoMap[bal.name] = bal
	withdraw := &cmdInfo{
		name:         "withdraw",
		channelLimit: true,
		handler:      (*guildConfigPresenter).cmdWithdrawHandler,
	}
	cmdInfoMap[withdraw.name] = withdraw

	setChannel := &cmdInfo{
		name:         "channel",
		managerCmd:   true,
		channelLimit: false,
		handler:      (*guildConfigPresenter).cmdChannelHandler,
	}
	cmdInfoMap[setChannel.name] = setChannel

}

func (p *guildConfigPresenter) cmdWithdrawHandler(parts *msgParts) {
	cmdPrefix := parts.prefix
	userID := parts.m.Author.ID
	userMention := parts.m.Author.Mention()
	symbol := parts.symbol
	presenter, ok := coinPresenters[symbol]
	if !ok {
		return
	}
	withdrawMinAmount, _ := amount.FromFloat64(presenter.coin.Withdraw.Min)
	minTxFee, _ := amount.FromFloat64(presenter.coin.Withdraw.TxFee)
	txFeePercent := presenter.coin.Withdraw.TxFeePercent * 100
	withdrawUsageInfo := &cmdWithdrawUserInfo{
		cmdUsageInfo: cmdUsageInfo{
			tmplName:        "withdrawUsage",
			IsShowUsageHint: true,
			CmdName:         "withdraw",
			UserMention:     userMention,
			Prefix:          string(cmdPrefix),
			Symbol:          string(symbol),
		},
		WithdrawMin:  withdrawMinAmount,
		TxFeePercent: txFeePercent,
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
			"Symbol":      symbol,
		})
		parts.channelMessageSend(msg)
		return
	}

	address := parts.contents[0]
	validateAddress, err := presenter.rpc.ValidateAddress(address)
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
			"Symbol":      symbol,
		})
		parts.channelMessageSend(msg)
		return
	}
	if validateAddress.IsMine {
		msg := msgFromTmpl("withdrawBotAddrErr", tmplValueMap{
			"UserMention": userMention,
			"Addr":        address,
			"Prefix":      cmdPrefix,
			"Symbol":      symbol,
		})
		parts.channelMessageSend(msg)
		return
	}
	pieer, err := presenter.db.UserByID(nil, userID)
	if err != nil {
		log.Errorf("[CMD]pie UserByID Error:%s", err)
		return
	}
	userAmount := amount.Zero
	userUnconfirmedAmount := amount.Zero
	if pieer != nil {
		userAmount = pieer.Amount
		userUnconfirmedAmount = pieer.UnconfirmedAmount
	}
	if userAmount.CmpFloat(withdrawAmount) == -1 {
		msg := msgFromTmpl("withdrawAmountNotEnoughErr", tmplValueMap{
			"UserMention":       userMention,
			"Amount":            userAmount,
			"UnconfirmedAmount": userUnconfirmedAmount,
			"Symbol":            symbol,
		})
		parts.channelMessageSend(msg)
		return
	}
	txfee, _ := amount.FromFloat64(withdrawAmount * txFeePercent)

	if txfee.Cmp(minTxFee) == -1 {
		txfee = minTxFee
	}
	withdrawAmountProxy, _ := amount.FromFloat64(withdrawAmount)
	finalWithdrawAmount := withdrawAmountProxy.Sub(txfee)

	withdrawTxID, err := presenter.rpc.SendToAddress(address, finalWithdrawAmount.Float64())
	if err != nil {
		msg := msgFromTmpl("walletMaintenance", userMention)
		parts.channelMessageSend(msg)
		return
	}
	msg := msgFromTmpl("withdrawSuccess", tmplValueMap{
		"UserMention": userMention,
		"Amount":      withdrawAmountProxy,
		"Symbol":      symbol,
		"Addr":        address,
		"TxFee":       txfee,
		"TxExpUrl":    presenter.coin.TxExplorerURL,
		"TxID":        withdrawTxID,
	})
	parts.channelMessageSend(msg)
	err = presenter.db.UserAmountUpsert(userID, parts.m.Author.Username, -withdrawAmount)
	if err != nil {
		log.Errorf("[%s] Withdraw Amount Update Error:%s[%s][%s][%s][%.8f]", symbol, err, userID, parts.m.Author.Username, withdrawTxID, withdrawAmount)
	}
	presenter.db.SaveWithdraw(userID, address, withdrawTxID, withdrawAmount)
}

func (p *guildConfigPresenter) cmdBalHandler(parts *msgParts) {
	userID := parts.m.Author.ID
	symbol := parts.symbol
	presenter := coinPresenters[symbol]
	user, err := presenter.db.UserByID(nil, userID)
	if err != nil {
		log.Errorf("[%s] Deposit UserByID Error:%s", symbol, err)
		return
	}
	confirmed := amount.Zero
	unconfirmed := amount.Zero
	if user != nil {
		confirmed = user.Amount
		unconfirmed = user.UnconfirmedAmount
	}
	msg := msgFromTmpl("balAmount", tmplValueMap{
		"UserMention":       parts.m.Author.Mention(),
		"Amount":            confirmed,
		"UnconfirmedAmount": unconfirmed,
		"Symbol":            symbol,
	})
	parts.channelMessageSend(msg)

}

func (p *guildConfigPresenter) cmdDepositHandler(parts *msgParts) {
	userID := parts.m.Author.ID
	userMention := parts.m.Author.Mention()
	symbol := parts.symbol
	presenter := coinPresenters[symbol]
	user, err := presenter.db.UserByID(nil, userID)
	if err != nil {
		log.Errorf("[%s] Deposit UserByID Error:%s", symbol, err)
		return
	}
	if user != nil && user.Address != "" {
		msg := msgFromTmpl("depositInfo", tmplValueMap{
			"UserMention": userMention,
			"Symbol":      symbol,
			"Addr":        user.Address,
		})
		parts.channelMessageSend(msg)
		return
	}
	address, err := presenter.rpc.GetNewAddress(userID)
	if err != nil {
		msg := msgFromTmpl("walletMaintenance", userMention)
		parts.channelMessageSend(msg)
		return
	}
	err = presenter.db.UserAddressUpsert(userID, parts.m.Author.Username, address, user == nil)
	if err != nil {
		log.Errorf("[%s] Deposit UserAddressUpsert Error:%s", symbol, err)
		return
	}
	msg := msgFromTmpl("depositInfo", tmplValueMap{
		"UserMention": userMention,
		"Symbol":      symbol,
		"Addr":        address,
	})
	parts.channelMessageSend(msg)
}

func (p *guildConfigPresenter) cmdPieHelperHandler(parts *msgParts) {
	cmdPrefix := parts.prefix
	symbol := parts.symbol

	presenter := coinPresenters[symbol]
	coinConfig := presenter.coin
	isManager := parts.isManager
	pieMinAmount, _ := amount.FromFloat64(coinConfig.Pie.Min)
	withdrawMinAmount, _ := amount.FromFloat64(coinConfig.Withdraw.Min)
	minTxFee, _ := amount.FromFloat64(coinConfig.Withdraw.TxFee)
	txFeePercent := coinConfig.Withdraw.TxFeePercent * 100
	cmdMsg := &cmdHelpUsageInfo{
		cmdUsageInfo: cmdUsageInfo{
			tmplName:        "helpUsage",
			IsShowUsageHint: false,
			CmdName:         "help",
			UserMention:     parts.m.Author.Mention(),
			Prefix:          string(cmdPrefix),
			Symbol:          string(symbol),
		},
		cmdPieUsageInfo: cmdPieUsageInfo{
			PieMin: pieMinAmount,
		},
		cmdWithdrawUserInfo: cmdWithdrawUserInfo{
			WithdrawMin:  withdrawMinAmount,
			TxFeePercent: txFeePercent,
			TxFeeMin:     minTxFee,
		},
		IsManager: isManager,
	}
	parts.channelMessageSend(cmdMsg.String())
}

func (p *guildConfigPresenter) pieReceivers(s *discordgo.Session, guild *discordgo.Guild, channelID, pieUserID string, isEveryone bool, isNeedOnline bool, roles []string, users []*discordgo.User) ([]*discordgo.User, error) {
	receivers := []*discordgo.User{}
	for _, member := range guild.Members {
		userID := member.User.ID
		switch {
		case member.User.Bot:
			fallthrough
		case member.User.ID == pieUserID:
			continue
		}

		userPermission, err := s.State.UserChannelPermissions(userID, channelID)
		if err != nil {
			log.Errorf("PieReceivers get permission Error:%s", err)
			continue
		}
		if (userPermission & discordgo.PermissionReadMessages) != discordgo.PermissionReadMessages {
			continue
		}
		isOnline := false
		presence, err := s.State.Presence(guild.ID, userID)
		if err == nil && presence.Status == discordgo.StatusOnline {
			isOnline = true
		}
		isAdd := false
		if isEveryone && !isNeedOnline {
			isAdd = true
		} else if isEveryone && isOnline {
			isAdd = true
		} else if len(roles) > 0 {
			rolesStr := strings.Join(roles, "|")
			for _, role := range member.Roles {
				if strings.Contains(rolesStr, role) {
					isAdd = true
					break
				}
			}
		} else if len(users) > 0 {
			for _, user := range users {
				if member.User.ID == user.ID {
					isAdd = true
					break
				}
			}
		}
		if isAdd {
			receivers = append(receivers, member.User)
		}
	}
	return receivers, nil
}

const eachMsgReceiverNum = 30

func (p *guildConfigPresenter) cmdPieHandler(parts *msgParts) {
	userID := parts.m.Author.ID
	userMention := parts.m.Author.Mention()
	cmdPrefix := parts.prefix
	symbol := parts.symbol
	presenter := coinPresenters[symbol]
	coinConfig := presenter.coin
	pieMinAmount, err := amount.FromFloat64(coinConfig.Pie.Min)
	pieUsageInfo := &cmdPieUsageInfo{
		cmdUsageInfo: cmdUsageInfo{
			tmplName:        "pieUsage",
			IsShowUsageHint: true,
			CmdName:         "pie",
			UserMention:     userMention,
			Prefix:          string(cmdPrefix),
			Symbol:          string(symbol),
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
			"Symbol":      symbol,
		})
		parts.channelMessageSend(msg)
		return
	}

	pieer, err := presenter.db.UserByID(nil, userID)
	if err != nil {
		log.Errorf("[CMD]pie UserByID Error:%s", err)
		return
	}
	userAmount := amount.Zero
	userImmatureAmount := amount.Zero
	if pieer != nil {
		userAmount = pieer.Amount
		userImmatureAmount = pieer.UnconfirmedAmount
	}
	if userAmount.Cmp(sendAmount) == -1 {
		msg := msgFromTmpl("pieNotEnoughAmountErr", tmplValueMap{
			"UserMention":       userMention,
			"Prefix":            cmdPrefix,
			"Amount":            userAmount,
			"UnconfirmedAmount": userImmatureAmount,
			"Symbol":            symbol,
		})
		parts.channelMessageSend(msg)
		return
	}

	isEveryone := false
	isNeedOnline := true
	if partLen == 1 {
		isEveryone = true
	} else if parts.m.MentionEveryone {
		isEveryone = true
		for _, part := range parts.contents {
			if strings.Contains(part, "@everyone") {
				isNeedOnline = false
				break
			}
		}
	}

	receivers, err := p.pieReceivers(parts.s, parts.guild, parts.m.ChannelID, userID, isEveryone, isNeedOnline, parts.m.MentionRoles, parts.m.Mentions)
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
			"Symbol":        symbol,
			"ReceiverCount": receiversLen,
		})
		parts.channelMessageSend(msg)
		return
	}

	err = presenter.db.UserAmountUpsert(userID, parts.m.Author.Username, -sendAmount.Float64())
	if err != nil {
		log.Errorf("Pie modify sender amount error:%s", err)
		return
	}

	receiversMap := make(map[int][]string)
	for i, receiver := range receivers {
		//msg index
		index := int(math.Floor(float64(i) / eachMsgReceiverNum))
		receiversMap[index] = append(receiversMap[index], receiver.Mention())
		err = presenter.db.UserAmountUpsert(receiver.ID, receiver.Username, amountEach.Float64())
		if err != nil {
			log.Errorf("Pie modify receiver amount error:%s", err)
		}
	}

	for _, receivers := range receiversMap {
		msg := msgFromTmpl("pieSuccess", tmplValueMap{
			"CoinName":   coinConfig.Name,
			"AmountEach": amountEach,
			"Symbol":     symbol,
			"Receivers":  receivers,
		})
		parts.channelMessageSend(msg)
	}

}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
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
	guild, err := guild(s, channel.GuildID)
	if err != nil {
		log.Error("messageCreateHandler guild Error:", err)
		return
	}
	gcp, ok := guildConfigPresenters[channel.GuildID]
	if !ok {
		gcp = guildConfigPresenters.initGuildConfigPresenter(channel.GuildID)
	}
	msgParts := &msgParts{
		s:        s,
		m:        m,
		channel:  channel,
		guild:    guild,
		contents: cntParts[1:],
	}
	if cntParts[0] == "?pie" {
		gcp.cmdMainPie(msgParts)
		return
	}
	prefixList := gcp.prefixList()
	if len(prefixList) == 0 {
		log.Error("Prefix List is Empty:", channel.GuildID)
		return
	}
	var prefix prefixWrap
	for _, pfx := range prefixList {
		if strings.HasPrefix(m.Content, string(pfx)) {
			prefix = pfx
			break
		}
	}
	if prefix == "" {
		// log.Error("can't find match prefix for:", channel.GuildID)
		return
	}
	// just only prefix
	if strings.Compare(string(prefix), cntParts[0]) == 0 {
		return
	}

	msgParts.prefix = prefix

	symbol, err := gcp.symbolByPrefix(prefix)
	msgParts.symbol = symbol

	cmd := strings.Replace(cntParts[0], string(prefix), "", 1)
	isInChannel := gcp.gccMap[symbol].inChannels(m.ChannelID)
	if cmdInfo, ok := cmdInfoMap[cmd]; ok {
		if cmdInfo.channelLimit && !isInChannel {
			return
		}
		isManager := gcp.isBotManager(s, guild, m.Author.ID)
		msgParts.isManager = isManager
		if cmdInfo.managerCmd && !isManager {
			return
		}
		cmdInfo.handler(gcp, msgParts)
	}
}

func channel(s *discordgo.Session, channelID string) (*discordgo.Channel, error) {
	channel, err := s.State.Channel(channelID)
	if err != nil {
		channel, err = s.Channel(channelID)
		if err != nil {
			log.Errorf("Get Channel Error:%s", err)
			return nil, err
		}
	}
	return channel, err
}

func guild(s *discordgo.Session, guildID string) (*discordgo.Guild, error) {
	guild, err := s.State.Guild(guildID)
	if err != nil {
		guild, err = s.Guild(guildID)
		if err != nil {
			log.Errorf("Get guild Error:%s", err)
			return nil, err
		}
	}
	return guild, err
}
