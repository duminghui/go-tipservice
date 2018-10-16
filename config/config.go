// Package config provides ...
package config

import (
	"encoding/json"
	"io/ioutil"

	rpcclient "github.com/duminghui/go-rpcclient"
)

type Withdraw struct {
	Min          float64 `json:"min"`
	TxFee        float64 `json:"txfee"`
	TxFeePercent float64 `json:"txfeepercent"`
}

type Pie struct {
	Min            float64 `json:"min"`
	EachMinReceive float64 `json:"eachminreceive"`
}

//VipGuildID:for control vip set command
type CoinInfo struct {
	Name                     string                `json:"name"`
	Symbol                   string                `json:"symbol"`
	Website                  string                `json:"website"`
	IconURL                  string                `json:"iconUrl"`
	VipGuildID               string                `json:"vipguildid"`
	Database                 string                `json:"database"`
	TxExplorerURL            string                `json:"txexplorer"`
	MinConfirmations4Deposit int64                 `json:"minConfirmations4Deposit"`
	ScanTxDir                string                `json:"scanTxDir"`
	RPC                      *rpcclient.ConnConfig `json:"rpc"`
	Withdraw                 *Withdraw             `json:"withdraw"`
	Pie                      *Pie                  `json:"pie"`
}

func FromFile(file string, v interface{}) (interface{}, error) {
	configBytes, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(configBytes, v)
	if err != nil {
		return nil, err
	}
	return v, nil

}
