package sequencesender

import (
	"encoding/json"
	jTypes "github.com/0xPolygon/cdk-data-availability/rpc"
	daTypes "github.com/0xPolygon/cdk-data-availability/types"
	"github.com/0xPolygonHermez/zkevm-node/etherman/types"
	"github.com/0xPolygonHermez/zkevm-node/log"
	kzg_sdk "github.com/domicon-labs/kzg-sdk"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"math/big"
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

func (s *SequenceSender) GenerateDomiconDA(sequences []types.Sequence) (*DomiconDAItem, error) {
	sequence := daTypes.Sequence{
		Batches: []daTypes.Batch{},
	}
	for _, seq := range sequences {
		sequence.Batches = append(sequence.Batches, daTypes.Batch{
			Number:         jTypes.ArgUint64(seq.BatchNumber),
			GlobalExitRoot: seq.GlobalExitRoot,
			Timestamp:      jTypes.ArgUint64(seq.Timestamp),
			Coinbase:       s.cfg.L2Coinbase,
			L2Data:         seq.BatchL2Data,
		})
	}
	signedSequence, err := sequence.Sign(s.privKey)
	if err != nil {
		log.Errorw("sequence Sign", "error", err)
		return nil, err
	}

	rawData, err := json.Marshal(signedSequence)
	if err != nil {
		log.Errorw("sequence data json Marshal", "error", err)
		return nil, err
	}
	dataCM, err := s.domiconKzg.GenerateDataCommit(rawData)
	if err != nil {
		log.Errorw("domicon GenerateDataCommit", "error", err)
		return nil, err
	}
	commit := dataCM.Bytes()

	singer := kzg_sdk.NewEIP155FdSigner(big.NewInt(31337)) // todo
	var length uint64 = uint64(len(rawData))
	var index uint64 = 1
	_, sigData, err := kzg_sdk.SignFd(s.cfg.SenderAddress, common.HexToAddress(""), 0, index, length, commit[:], singer, s.privKey)
	if err != nil {
		return nil, err
	}
	return &DomiconDAItem{
		index:       uint64(index),
		legth:       length,
		broadcaster: s.cfg.SenderAddress,
		user:        common.HexToAddress(""), //todo
		commitment:  hexutil.Bytes(commit[:]),
		sign:        hexutil.Bytes(sigData),
		data:        hexutil.Bytes(rawData)}, nil
}
