// Package main provides ...
package main

import (
	"errors"

	"github.com/bwmarrin/discordgo"
	"github.com/duminghui/go-tipservice/amount"
	"github.com/duminghui/go-tipservice/db"
)

type pieReceiverGenerator interface {
	Receivers() ([]*discordgo.User, error)
}

type pie struct {
	symbol            symbol
	userID            string
	userName          string
	amount            amount.Amount
	receiverGenerator pieReceiverGenerator
}

type pieReport struct {
	pieer         *db.User
	receivers     []*discordgo.User
	receiverCount int
	eachAmount    amount.Amount
}

var (
	errPieUnknown             = errors.New("Unknown")
	errPieNoReceiverGenerator = errors.New("PieNoReceiverGeneratorError")
	errPieNoSymbol            = errors.New("PieNoSymbol")
	errPieAmountMin           = errors.New("PieAmountMinError")
	errPieUserNotExists       = errors.New("PieUserNotExistsError")
	errPieNotEnoughAmount     = errors.New("PieNotEnoughAmountError")
	errPieNoReceiver          = errors.New("PieNoReceiverError")
	errPieNotEnoughEachAmount = errors.New("PieNotEnoughEachAmountError")
)

func (p *pie) pie() (*pieReport, error) {
	if p.receiverGenerator == nil {
		return nil, errPieNoReceiverGenerator
	}
	cp, ok := coinPresenters[p.symbol]
	if !ok {
		return nil, errPieNoSymbol
	}
	pieMinAmount, err := amount.FromFloat64(cp.coinInfo.Pie.Min)
	if err != nil {
		return nil, err
	}
	sendAmount := p.amount
	if sendAmount.Cmp(pieMinAmount) == -1 {
		return nil, errPieAmountMin
	}
	userID := p.userID
	pieer, err := cp.dbSymbol.UserByID(nil, userID)
	if err != nil {
		return nil, err
	}
	if pieer == nil {
		return nil, errPieUserNotExists
	}
	pieReport := new(pieReport)
	pieReport.pieer = pieer
	userAmount := pieer.Amount
	if userAmount.Cmp(sendAmount) == -1 {
		return pieReport, errPieNotEnoughAmount
	}
	receivers, err := p.receiverGenerator.Receivers()
	if err != nil {
		return pieReport, err
	}
	receiverCount := len(receivers)
	pieReport.receiverCount = receiverCount
	if receiverCount == 0 {
		return pieReport, errPieNoReceiver
	}
	amountEach := sendAmount.DivFloat64(float64(receiverCount))
	pieReport.eachAmount = amountEach
	if amountEach.Cmp(amount.Zero) == 0 {
		return pieReport, errPieNotEnoughEachAmount
	}
	userName := p.userName
	err = cp.dbSymbol.UserAmountSub(nil, userID, userName, sendAmount)
	if err != nil {
		return pieReport, err
	}
	receiversSuccess := make([]*discordgo.User, 0)
	for _, receiver := range receivers {
		err = cp.dbSymbol.UserAmountAddUpsert(nil, receiver.ID, receiver.Username, amountEach)
		if err == nil {
			receiversSuccess = append(receiversSuccess, receiver)
		}
	}
	pieReport.receiverCount = len(receiversSuccess)
	pieReport.receivers = receiversSuccess
	return pieReport, nil
}
