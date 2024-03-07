package rpcapi

import (
	"context"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ten-protocol/go-ten/go/common/gethapi"
	rpc2 "github.com/ten-protocol/go-ten/go/enclave/rpc"
	"github.com/ten-protocol/go-ten/lib/gethfork/rpc"
)

type TransactionAPI struct {
	we *Services
}

func NewTransactionAPI(we *Services) *TransactionAPI {
	return &TransactionAPI{we}
}

func (s *TransactionAPI) GetBlockTransactionCountByNumber(ctx context.Context, blockNr rpc.BlockNumber) *hexutil.Uint {
	count, err := PlaintextTenRPCCall[hexutil.Uint](ctx, s.we, &CacheCfg{cacheTTLCallback: func() time.Duration {
		if blockNr > 0 {
			return longCacheTTL
		}
		return shortCacheTTL
	}}, "eth_getBlockTransactionCountByNumber", blockNr)
	if err != nil {
		return nil
	}
	return count
}

func (s *TransactionAPI) GetBlockTransactionCountByHash(ctx context.Context, blockHash common.Hash) *hexutil.Uint {
	count, err := PlaintextTenRPCCall[hexutil.Uint](ctx, s.we, &CacheCfg{cacheTTL: longCacheTTL}, "eth_getBlockTransactionCountByHash", blockHash)
	if err != nil {
		return nil
	}
	return count
}

func (s *TransactionAPI) GetTransactionByBlockNumberAndIndex(ctx context.Context, blockNr rpc.BlockNumber, index hexutil.Uint) *rpc2.RpcTransaction {
	// not implemented
	return nil
}

func (s *TransactionAPI) GetTransactionByBlockHashAndIndex(ctx context.Context, blockHash common.Hash, index hexutil.Uint) *rpc2.RpcTransaction {
	// not implemented
	return nil
}

func (s *TransactionAPI) GetRawTransactionByBlockNumberAndIndex(ctx context.Context, blockNr rpc.BlockNumber, index hexutil.Uint) hexutil.Bytes {
	// not implemented
	return nil
}

func (s *TransactionAPI) GetRawTransactionByBlockHashAndIndex(ctx context.Context, blockHash common.Hash, index hexutil.Uint) hexutil.Bytes {
	// not implemented
	return nil
}

func (s *TransactionAPI) GetTransactionCount(ctx context.Context, address common.Address, blockNrOrHash rpc.BlockNumberOrHash) (*hexutil.Uint64, error) {
	return ExecAuthRPC[hexutil.Uint64](ctx, s.we, &ExecCfg{account: &address}, "eth_getTransactionCount", address, blockNrOrHash)
}

func (s *TransactionAPI) GetTransactionByHash(ctx context.Context, hash common.Hash) (*rpc2.RpcTransaction, error) {
	return ExecAuthRPC[rpc2.RpcTransaction](ctx, s.we, &ExecCfg{tryAll: true}, "eth_getTransactionByHash", hash)
}

func (s *TransactionAPI) GetRawTransactionByHash(ctx context.Context, hash common.Hash) (hexutil.Bytes, error) {
	tx, err := ExecAuthRPC[hexutil.Bytes](ctx, s.we, &ExecCfg{tryAll: true}, "eth_getRawTransactionByHash", hash)
	if tx != nil {
		return *tx, err
	}
	return nil, err
}

func (s *TransactionAPI) GetTransactionReceipt(ctx context.Context, hash common.Hash) (map[string]interface{}, error) {
	txRec, err := ExecAuthRPC[map[string]interface{}](ctx, s.we, &ExecCfg{tryUntilAuthorised: true}, "eth_getTransactionReceipt", hash)
	if txRec != nil {
		return *txRec, err
	}
	return nil, err
}

func (s *TransactionAPI) SendTransaction(ctx context.Context, args gethapi.TransactionArgs) (common.Hash, error) {
	txRec, err := ExecAuthRPC[common.Hash](ctx, s.we, &ExecCfg{account: args.From}, "eth_sendTransaction", args)
	if txRec != nil {
		return *txRec, err
	}
	return common.Hash{}, err
}

type SignTransactionResult struct {
	Raw hexutil.Bytes      `json:"raw"`
	Tx  *types.Transaction `json:"tx"`
}

func (s *TransactionAPI) FillTransaction(ctx context.Context, args gethapi.TransactionArgs) (*SignTransactionResult, error) {
	// not implemented
	return nil, nil
}

func (s *TransactionAPI) SendRawTransaction(ctx context.Context, input hexutil.Bytes) (common.Hash, error) {
	txRec, err := ExecAuthRPC[common.Hash](ctx, s.we, &ExecCfg{tryAll: true}, "eth_sendRawTransaction", input)
	if txRec != nil {
		return *txRec, err
	}
	return common.Hash{}, err
}

func (s *TransactionAPI) PendingTransactions() ([]*rpc2.RpcTransaction, error) {
	// not implemented
	return nil, nil
}

func (s *TransactionAPI) Resend(ctx context.Context, sendArgs gethapi.TransactionArgs, gasPrice *hexutil.Big, gasLimit *hexutil.Uint64) (common.Hash, error) {
	txRec, err := ExecAuthRPC[common.Hash](ctx, s.we, &ExecCfg{account: sendArgs.From}, "eth_resend", sendArgs, gasPrice, gasLimit)
	if txRec != nil {
		return *txRec, err
	}
	return common.Hash{}, err
}
