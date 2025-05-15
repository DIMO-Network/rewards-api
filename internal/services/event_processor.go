package services

import (
	"github.com/ethereum/go-ethereum/common"
)

type TokenConfig struct {
	Tokens []struct {
		ChainID int64          `yaml:"chainId"`
		Address common.Address `yaml:"address"`
		RPCURL  string         `yaml:"rpcUrl"`
	} `yaml:"tokens"`
}
