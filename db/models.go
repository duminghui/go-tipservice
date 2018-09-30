// Package db provides ...
package db

import "github.com/duminghui/go-tipservice/amount"

const colDeposit = "deposit"

type Deposit struct {
	UserID      string  `bson:"user_id,omitempty"`
	Amount      float64 `bson:"amount,omitempty"`
	TxID        string  `bson:"txid,omitempty"`
	Address     string  `bson:"addresses,omitempty"`
	IsConfirmed bool    `bson:"isConfirmed,omitempty"`
}

const colWithdraw = "withdraw"

type Withdraw struct {
	UserID  string  `bson:"user_id"`
	Amount  float64 `bson:"amount"`
	Address string  `bson:"address"`
	TxID    string  `bson:"txid"`
}

const colNoOwnerDeposit = "no_owner_deposit"

type NoOwnerDeposit struct {
	TxID    string  `bson:"txid"`
	Address string  `bson:"address"`
	Amount  float64 `bson:"amount"`
}

const colUser = "user"

type User struct {
	UserID            string        `bson:"user_id,omitempty"`
	UserName          string        `bson:"user_name,omitempty"`
	Address           string        `bson:"address,omitempty"`
	Amount            amount.Amount `bson:"amount,omitempty"`
	UnconfirmedAmount amount.Amount `bson:"unconfirmed_amount,omitempty"`
}
