// Package main provides ...
package main

import (
	"strings"

	"github.com/duminghui/go-tipservice/amount"
)

type cmdUsageInfo struct {
	tmplName        string
	IsShowUsageHint bool
	CmdName         string
	UserMention     string
	Prefix          string
	Symbol          string
}

func (c *cmdUsageInfo) String() string {
	return msgFromTmpl(c.tmplName, c)
}

type cmdPieUsageInfo struct {
	cmdUsageInfo
	PieMin amount.Amount
}

func (c *cmdPieUsageInfo) String() string {
	return msgFromTmpl(c.tmplName, c)
}

type cmdWithdrawUserInfo struct {
	cmdUsageInfo
	WithdrawMin  amount.Amount
	TxFeePercent float64
	TxFeeMin     amount.Amount
}

func (c *cmdWithdrawUserInfo) String() string {
	return msgFromTmpl(c.tmplName, c)
}

type cmdHelpUsageInfo struct {
	cmdUsageInfo
	cmdPieUsageInfo
	cmdWithdrawUserInfo
	IsManager bool
}

func (c *cmdHelpUsageInfo) String() string {
	return msgFromTmpl(c.tmplName, c)
}

func msgFromTmpl(tmplName string, data interface{}) string {
	buf := new(strings.Builder)
	err := msgTmpl.ExecuteTemplate(buf, tmplName, data)
	if err != nil {
		return err.Error()
	}
	return buf.String()
}
