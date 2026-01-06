// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package config

import "github.com/zeromicro/go-zero/rest"

type Config struct {
	rest.RestConf

	MySQL      MySQLConfig
	Redis      RedisConfig
	Blockchain BlockchainConfig
	Auth       AuthConfig
	Business   BusinessConfig
	X402       X402Config
}

type MySQLConfig struct {
	Host            string
	Port            int
	User            string
	Password        string
	Database        string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime int
}

type RedisConfig struct {
	Host     string
	Password string
	DB       int
	PoolSize int
}

type BlockchainConfig struct {
	RpcUrl     string
	ChainId    int64
	PrivateKey string
}

type AuthConfig struct {
	AccessSecret  string
	AccessExpire  int64
	RefreshSecret string
	RefreshExpire int64
}

type BusinessConfig struct {
	PlatformFeeRate   float64
	MinMarketDuration int64
	MaxMarketDuration int64
	MinLiquidity      float64
}

// X402Config X402 支付配置
type X402Config struct {
	Enabled   bool    `json:"Enabled" yaml:"Enabled"`
	Amount    float64 `json:"Amount" yaml:"Amount"`       // 服务费用（SOL）
	Recipient string  `json:"Recipient" yaml:"Recipient"` // Solana 收款地址
	RPCURL    string  `json:"RPCURL" yaml:"RPCURL"`       // Solana RPC 端点
}
