package sequencesender

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

type DomiconDAItem struct {
	index       uint64
	legth       uint64
	broadcaster common.Address
	user        common.Address
	commitment  hexutil.Bytes
	sign        hexutil.Bytes
	data        hexutil.Bytes
}

func (s *SequenceSender) GenerateDomiconDA() *DomiconDAItem {
	//todo
	return &DomiconDAItem{}
}
