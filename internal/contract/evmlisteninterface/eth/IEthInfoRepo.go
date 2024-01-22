package eth

import (
	"time"
)

type EthDefaultInfo struct {
	CreatedAt        time.Time
	UpdatedAt        time.Time
	ChainID          string `db:"eth_chain_id"`
	LastScannedBlock int64  `db:"eth_last_scan_block"`
}

func (*EthDefaultInfo) TableName() string {
	return "eth_block_info"

}

type IEthInfoRepo interface {
	Get(ChainID string) (*EthDefaultInfo, error)
	Create(ChainID string, info *EthDefaultInfo) error
	Update(ChainID string, info *EthDefaultInfo) error
}
