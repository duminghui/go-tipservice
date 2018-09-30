// Package main provides ...
package main

import (
	"bytes"
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/duminghui/go-tipservice/amount"
)

type cmdHandler func(*discordgo.Session, *discordgo.MessageCreate, *msgParts)

var cmdMap = make(map[string]cmdHandler)

type msgParts struct {
	prefix prefixWrap
	cmd    string
	parts  []string
}

var cmdSetChannel = "setChannel"

func reigsterBotCmdHandler() {
	cmdMap["piehelp"] = cmdPieHelperHandler
	cmdMap["pie"] = cmdPieHandler
	cmdMap["deposit"] = cmdDepositHandler
	cmdMap["bal"] = cmdBalHandler
	cmdMap["withdraw"] = cmdWithdrawHandler
	cmdMap[cmdSetChannel] = cmdSetChannelHandler
}

func cmdWithdrawHandler(s *discordgo.Session, m *discordgo.MessageCreate, parts *msgParts) {
	cmdPrefix := parts.prefix
	userID := m.Author.ID
	userMention := m.Author.Mention()
	channelID := m.ChannelID
	channel, err := channel(s, channelID)
	if err != nil {
		log.Error("cmdWithdraw Error:", err)
		return
	}
	symbol, err := guildConfigManagers.symbolByPrefix(channel.GuildID, cmdPrefix)
	if err != nil {
		log.Error("cmdWithdraw Error:", err)
		return
	}
	presenter, ok := coinPresenters[symbol]
	if !ok {
		return
	}
	withdrawMinAmount := presenter.coin.Withdraw.Min
	minTxFee := presenter.coin.Withdraw.TxFee
	txFeePercent := presenter.coin.Withdraw.TxFeePercent
	cmdPartErrMsg := fmt.Sprintf("%s withdraw command usage:\n```withdraw command format:\n  %swithdraw <address> <amount>\n--minimum amount:%.8f %s\n--txfee:%g%% or %.8f %s```", userMention, cmdPrefix, withdrawMinAmount, symbol, txFeePercent*100, minTxFee, symbol)
	if len(parts.parts) != 2 {
		s.ChannelMessageSend(channelID, cmdPartErrMsg)
		return
	}

	withdrawAmount, err := strconv.ParseFloat(parts.parts[1], 64)
	if err != nil {
		s.ChannelMessageSend(channelID, cmdPartErrMsg)
		return
	}

	if withdrawAmount < withdrawMinAmount {
		msg := fmt.Sprintf("%s withdraw minimum amount is `%.8f %s`", userMention, withdrawMinAmount, symbol)
		s.ChannelMessageSend(channelID, msg)
		return
	}

	address := parts.parts[0]
	validateAddress, err := presenter.rpc.ValidateAddress(address)
	if err != nil {
		log.Error("[CMD]withdraw ValidateAddress Error:", err)
		return
	}
	if !validateAddress.IsValid {
		msg := fmt.Sprintf("%s `%s` is not %s address", userMention, address, symbol)
		s.ChannelMessageSend(channelID, msg)
		return
	}
	if validateAddress.IsMine {
		msg := fmt.Sprintf("%s `%s` is in bot's wallet, you can use `%spie` command to give some one %s", userMention, address, cmdPrefix, symbol)
		s.ChannelMessageSend(channelID, msg)
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
	if userAmount.CmpFloat(withdrawAmount) == -1 {
		msg := fmt.Sprintf("%s your don't have enough amount to winthdraw\n ```Balance amount:%s %s\nUnconfirmed amount:%s %s```", userMention, userAmount, symbol, userImmatureAmount, symbol)
		s.ChannelMessageSend(channelID, msg)
		return
	}
	txfee := withdrawAmount * txFeePercent
	if txfee < minTxFee {
		txfee = minTxFee
	}
	withdrawAmountProxy, _ := amount.FromFloat64(withdrawAmount)
	txfeeProxy, _ := amount.FromFloat64(txfee)
	finalWithdrawAmount := withdrawAmountProxy.Sub(txfeeProxy)

	withdrawTxID, err := presenter.rpc.SendToAddress(address, finalWithdrawAmount.Float64())
	if err != nil {
		msg := fmt.Sprintf("%s Wallet maintenance", userMention)
		s.ChannelMessageSend(channelID, msg)
		return
	}
	msg := fmt.Sprintf("%s you withdraw %.8f %s to `%s`\ntxfee: %.8f %s\n%s%s", userMention, withdrawAmount, symbol, address, txfee, symbol, presenter.coin.TxExplorerURL, withdrawTxID)
	s.ChannelMessageSend(channelID, msg)
	err = presenter.db.UserAmountUpsert(userID, m.Author.Username, -withdrawAmount)
	if err != nil {
		log.Errorf("[%s] Withdraw Amount Update Error:%s[%s][%s][%s][%.8f]", symbol, err, userID, m.Author.Username, withdrawTxID, withdrawAmount)
	}
	presenter.db.SaveWithdraw(userID, address, withdrawTxID, withdrawAmount)
}

func cmdBalHandler(s *discordgo.Session, m *discordgo.MessageCreate, parts *msgParts) {
	cmdPrefix := parts.prefix
	userID := m.Author.ID
	channelID := m.ChannelID
	channel, err := channel(s, channelID)
	if err != nil {
		log.Error("cmdBal Error:", err)
		return
	}
	symbol, err := guildConfigManagers.symbolByPrefix(channel.GuildID, cmdPrefix)
	if err != nil {
		log.Error("cmdBal Error:", err)
		return
	}
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
	msgFormat := "%s Your balance:\n```Confirmed: %s %s\nUnconfirmed: %s %s```"
	msg := fmt.Sprintf(msgFormat, m.Author.Mention(), confirmed, symbol, unconfirmed, symbol)
	s.ChannelMessageSend(m.ChannelID, msg)

}

func cmdDepositHandler(s *discordgo.Session, m *discordgo.MessageCreate, parts *msgParts) {
	cmdPrefix := parts.prefix
	userID := m.Author.ID
	userMention := m.Author.Mention()
	channelID := m.ChannelID
	channel, err := channel(s, channelID)
	if err != nil {
		log.Error("cmdBal Error:", err)
		return
	}
	symbol, err := guildConfigManagers.symbolByPrefix(channel.GuildID, cmdPrefix)
	if err != nil {
		log.Error("cmdBal Error:", err)
		return
	}
	presenter := coinPresenters[symbol]
	user, err := presenter.db.UserByID(nil, userID)
	if err != nil {
		log.Errorf("[%s] Deposit UserByID Error:%s", symbol, err)
		return
	}
	if user != nil && user.Address != "" {
		msg := fmt.Sprintf("%s You deposit Address is:\n`%s`", userMention, user.Address)
		s.ChannelMessageSend(m.ChannelID, msg)
		return
	}
	address, err := presenter.rpc.GetNewAddress(userID)
	if err != nil {
		msg := fmt.Sprintf("%s Wallet maintenance", userMention)
		s.ChannelMessageSend(m.ChannelID, msg)
		return
	}
	err = presenter.db.UserAddressUpsert(userID, m.Author.Username, address, user == nil)
	if err != nil {
		log.Errorf("[%s] Deposit UserAddressUpsert Error:%s", symbol, err)
		return
	}
	msg := fmt.Sprintf("%s You deposit Address is:\n`%s`", userMention, address)
	s.ChannelMessageSend(m.ChannelID, msg)
}

func cmdPieHelperHandler(s *discordgo.Session, m *discordgo.MessageCreate, parts *msgParts) {
	s.ChannelMessageSend(m.ChannelID, "this is help command")
}

func pieReceivers(s *discordgo.Session, channelID, pieUserID string, isEveryone bool, isNeedOnline bool, roles []string, users []*discordgo.User) ([]*discordgo.User, error) {
	channel, err := channel(s, channelID)
	if err != nil {
		return nil, err
	}
	guild, err := guild(s, channel.GuildID)
	if err != nil {
		return nil, err
	}
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

func cmdPieHandler(s *discordgo.Session, m *discordgo.MessageCreate, parts *msgParts) {
	userID := m.Author.ID
	userMention := m.Author.Mention()
	cmdPrefix := parts.prefix
	channelID := m.ChannelID
	channel, err := channel(s, channelID)
	if err != nil {
		log.Error("cmdBal Error:", err)
		return
	}
	symbol, err := guildConfigManagers.symbolByPrefix(channel.GuildID, cmdPrefix)
	if err != nil {
		log.Error("cmdBal Error:", err)
		return
	}
	presenter := coinPresenters[symbol]
	coinConfig := presenter.coin
	pieMinAmount, err := amount.FromFloat64(coinConfig.Pie.Min)
	msg := fmt.Sprintf("%s pie command usage:\n```pie command format:\n  %spie [@receiver...] <amount>\n--minimum amount:%s %s```", userMention, cmdPrefix, pieMinAmount, symbol)
	partLen := len(parts.parts)
	if partLen == 0 {
		s.ChannelMessageSend(channelID, msg)
		return
	}
	sendAmount, err := amount.FromNumString(parts.parts[partLen-1])
	if err != nil {
		s.ChannelMessageSend(channelID, msg)
		return
	}

	if sendAmount < pieMinAmount {
		msg := fmt.Sprintf("%s Minimum amount `%s %s` allowed to be distribute", userMention, pieMinAmount, symbol)
		s.ChannelMessageSend(channelID, msg)
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
		msg := fmt.Sprintf("%s your don't have enough amount to distribute\n please use command `%sdeposit` to get deposit address\n```Balance amount:%s %s\nUnconfirmed amount:%s %s```", userMention, cmdPrefix, userAmount, symbol, userImmatureAmount, symbol)
		s.ChannelMessageSend(channelID, msg)
		return
	}

	isEveryone := false
	isNeedOnline := true
	if partLen == 1 {
		isEveryone = true
	} else if m.MentionEveryone {
		isEveryone = true
		for _, part := range parts.parts {
			if strings.Contains(part, "@everyone") {
				isNeedOnline = false
				break
			}
		}
	}

	receivers, err := pieReceivers(s, channelID, userID, isEveryone, isNeedOnline, m.MentionRoles, m.Mentions)
	if err != nil {
		log.Errorf("Pie get receivers error:%s", err)
		return
	}

	receiversLen := len(receivers)
	if receiversLen == 0 {
		msg := fmt.Sprintf("%s No people to be distribute pie, Try again when people are online", userMention)
		s.ChannelMessageSend(channelID, msg)
		return
	}

	amountEach := sendAmount.DivFloat64(float64(receiversLen))

	if amountEach.Cmp(amount.Zero) == 0 {
		msg := fmt.Sprintf("%s %s is not enough to distribute %d peoples", sendAmount, symbol, receiversLen)
		s.ChannelMessageSend(channelID, msg)
		return
	}

	err = presenter.db.UserAmountUpsert(userID, m.Author.Username, -sendAmount.Float64())
	if err != nil {
		log.Errorf("Pie modify sender amount error:%s", err)
		return
	}

	// Max msg count
	sendMsgCount := int(math.Ceil(float64(receiversLen) / eachMsgReceiverNum))
	sendMsgs := make([]*bytes.Buffer, sendMsgCount)
	for i := 0; i < sendMsgCount; i++ {
		sendMsgs[i] = new(bytes.Buffer)
		if i == 0 {
			msg := fmt.Sprintf(":lollipop: ~ ~ ~ ~ ~ ~ ~ ~ %s pie ~ ~ ~ ~ ~ ~ ~ ~:candy:\n%s %s to", coinConfig.Name, amountEach, symbol)
			sendMsgs[i].WriteString(msg)
		}
	}

	for i, receiver := range receivers {
		//msg index
		index := int(math.Floor(float64(i) / eachMsgReceiverNum))
		sendMsgs[index].WriteString(" ")
		sendMsgs[index].WriteString(receiver.Mention())
		err = presenter.db.UserAmountUpsert(receiver.ID, receiver.Username, amountEach.Float64())
		if err != nil {
			log.Errorf("Pie modify receiver amount error:%s", err)
		}
	}

	for _, sendMsg := range sendMsgs {
		s.ChannelMessageSend(channelID, sendMsg.String())
	}

}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.Bot || m.Author.ID == s.State.User.ID {
		return
	}
	cntParts := strings.Fields(m.Content)
	if cntParts[0] == "?pieSet" {
		cmdPieSet(s, m, cntParts[1:])
		return
	}
	channel, err := channel(s, m.ChannelID)
	if err != nil {
		log.Error("messageCreateHandler Error:", err)
		return
	}
	prefixList := guildConfigManagers.prefixList(channel.GuildID)
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
		log.Error("can't find match prefix for:", channel.GuildID)
		return
	}
	if strings.Compare(string(prefix), cntParts[0]) == 0 {
		return
	}
	gcm, ok := guildConfigManagers[channel.GuildID]
	if !ok {
		log.Error("No Guild Config for", channel.GuildID)
		return
	}
	symbol, err := guildConfigManagers.symbolByPrefix(channel.GuildID, prefix)
	cmd := strings.Replace(cntParts[0], string(prefix), "", 1)
	msgParts := &msgParts{
		prefix: prefix,
		cmd:    cmd,
		parts:  cntParts[1:],
	}
	if cmd == cmdSetChannel {
		cmdMap[cmdSetChannel](s, m, msgParts)
		return
	}
	isInChannel := gcm.guildCoinConfig[symbol].InChannels(m.ChannelID)
	if !isInChannel {
		return
	}
	if cmdFun, ok := cmdMap[cmd]; ok {
		cmdFun(s, m, msgParts)
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