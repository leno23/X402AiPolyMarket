// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package config

import "github.com/zeromicro/go-zero/rest"

type Config struct {
	rest.RestConf
	X402 X402Config
}

// X402Config X402 支付配置
type X402Config struct {
	Enabled   bool    `json:"Enabled" yaml:"Enabled"`
	Amount    float64 `json:"Amount" yaml:"Amount"`       // 服务费用（SOL）
	Recipient string  `json:"Recipient" yaml:"Recipient"` // Solana 收款地址
	RPCURL    string  `json:"RPCURL" yaml:"RPCURL"`       // Solana RPC 端点
}
