// Package config provides ...
package config

import (
	"encoding/json"
	"io/ioutil"

	rpcclient "github.com/duminghui/go-rpcclient"
	"github.com/duminghui/go-util/ulog"
	"github.com/duminghui/go-util/umgo"
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

type CoinInfo struct {
	Name                     string                `json:"name"`
	Symbol                   string                `json:"symbol"`
	Database                 string                `json:"database"`
	TxExplorerURL            string                `json:"txexplorer"`
	MinConfirmations4Deposit int64                 `json:"minConfirmations4Deposit"`
	ScanTxDir                string                `json:"scanTxDir"`
	RPC                      *rpcclient.ConnConfig `json:"rpc"`
	Withdraw                 *Withdraw             `json:"withdraw"`
	Pie                      *Pie                  `json:"pie"`
}

type Discord struct {
	Token          string `json:"token"`
	SuperManagerID string `json:"supermanagerid"`
}

type Config struct {
	Discord *Discord             `json:"discord"`
	Mongodb *umgo.ConnConfig     `json:"mongodb"`
	Log     *ulog.Config         `json:"log"`
	Infos   map[string]*CoinInfo `json:"infos"`
}

func New(file string) (*Config, error) {
	configBytes, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}
	var config Config
	err = json.Unmarshal(configBytes, &config)
	if err != nil {
		return nil, err
	}
	return &config, nil
}
