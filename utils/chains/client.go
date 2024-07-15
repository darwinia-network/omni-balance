package chains

import (
	"context"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/ethclient/simulated"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"math/big"
	"sync/atomic"
)

type MockClient interface {
	simulated.Client
}

type Client struct {
	index   atomic.Int64
	clients []*ethclient.Client
	chainId *big.Int
}

func NewTryClient(ctx context.Context, endpoints []string) (*Client, error) {
	t := &Client{}

	for _, v := range endpoints {
		client, err := ethclient.DialContext(ctx, v)
		if err != nil {
			logrus.Warnf("dial %s error: %s", v, err)
			continue
		}
		t.clients = append(t.clients, client)
	}
	if len(t.clients) == 0 {
		return nil, errors.Errorf("no available endpoint, endpoints: %+v", endpoints)
	}
	return t, nil
}

func (t *Client) Close() {
	for _, v := range t.clients {
		v.Close()
	}
}

func (t *Client) Clients() simulated.Client {
	return t.clients[t.index.Load()%int64(len(t.clients))]
}

func (t *Client) BlockNumber(ctx context.Context) (uint64, error) {
	return t.Clients().BlockNumber(ctx)
}

func (t *Client) BlockByHash(ctx context.Context, hash common.Hash) (*types.Block, error) {
	return t.Clients().BlockByHash(ctx, hash)
}

func (t *Client) BlockByNumber(ctx context.Context, number *big.Int) (*types.Block, error) {
	return t.Clients().BlockByNumber(ctx, number)
}

func (t *Client) HeaderByHash(ctx context.Context, hash common.Hash) (*types.Header, error) {
	return t.Clients().HeaderByHash(ctx, hash)
}

func (t *Client) HeaderByNumber(ctx context.Context, number *big.Int) (*types.Header, error) {
	return t.Clients().HeaderByNumber(ctx, number)
}

func (t *Client) TransactionCount(ctx context.Context, blockHash common.Hash) (uint, error) {
	return t.Clients().TransactionCount(ctx, blockHash)
}

func (t *Client) TransactionInBlock(ctx context.Context, blockHash common.Hash, index uint) (*types.Transaction, error) {
	return t.Clients().TransactionInBlock(ctx, blockHash, index)
}

func (t *Client) SubscribeNewHead(ctx context.Context, ch chan<- *types.Header) (ethereum.Subscription, error) {
	return t.Clients().SubscribeNewHead(ctx, ch)
}

func (t *Client) BalanceAt(ctx context.Context, account common.Address, blockNumber *big.Int) (*big.Int, error) {
	return t.Clients().BalanceAt(ctx, account, blockNumber)
}

func (t *Client) StorageAt(ctx context.Context, account common.Address, key common.Hash, blockNumber *big.Int) ([]byte, error) {
	return t.Clients().StorageAt(ctx, account, key, blockNumber)
}

func (t *Client) CodeAt(ctx context.Context, account common.Address, blockNumber *big.Int) ([]byte, error) {
	return t.Clients().CodeAt(ctx, account, blockNumber)
}

func (t *Client) NonceAt(ctx context.Context, account common.Address, blockNumber *big.Int) (uint64, error) {
	return t.Clients().NonceAt(ctx, account, blockNumber)
}

func (t *Client) CallContract(ctx context.Context, call ethereum.CallMsg, blockNumber *big.Int) ([]byte, error) {
	return t.Clients().CallContract(ctx, call, blockNumber)
}

func (t *Client) EstimateGas(ctx context.Context, call ethereum.CallMsg) (uint64, error) {
	return t.Clients().EstimateGas(ctx, call)
}

func (t *Client) SuggestGasPrice(ctx context.Context) (*big.Int, error) {
	return t.Clients().SuggestGasPrice(ctx)
}

func (t *Client) SuggestGasTipCap(ctx context.Context) (*big.Int, error) {
	return t.Clients().SuggestGasTipCap(ctx)
}

func (t *Client) FeeHistory(ctx context.Context, blockCount uint64, lastBlock *big.Int, rewardPercentiles []float64) (*ethereum.FeeHistory, error) {
	return t.Clients().FeeHistory(ctx, blockCount, lastBlock, rewardPercentiles)
}

func (t *Client) FilterLogs(ctx context.Context, q ethereum.FilterQuery) ([]types.Log, error) {
	return t.Clients().FilterLogs(ctx, q)
}

func (t *Client) SubscribeFilterLogs(ctx context.Context, q ethereum.FilterQuery, ch chan<- types.Log) (ethereum.Subscription, error) {
	return t.Clients().SubscribeFilterLogs(ctx, q, ch)
}

func (t *Client) PendingBalanceAt(ctx context.Context, account common.Address) (*big.Int, error) {
	return t.Clients().PendingBalanceAt(ctx, account)
}

func (t *Client) PendingStorageAt(ctx context.Context, account common.Address, key common.Hash) ([]byte, error) {
	return t.Clients().PendingStorageAt(ctx, account, key)
}

func (t *Client) PendingCodeAt(ctx context.Context, account common.Address) ([]byte, error) {
	return t.Clients().PendingCodeAt(ctx, account)
}

func (t *Client) PendingNonceAt(ctx context.Context, account common.Address) (uint64, error) {
	return t.Clients().PendingNonceAt(ctx, account)
}

func (t *Client) PendingTransactionCount(ctx context.Context) (uint, error) {
	return t.Clients().PendingTransactionCount(ctx)
}

func (t *Client) PendingCallContract(ctx context.Context, call ethereum.CallMsg) ([]byte, error) {
	return t.Clients().PendingCallContract(ctx, call)
}

func (t *Client) TransactionByHash(ctx context.Context, txHash common.Hash) (tx *types.Transaction, isPending bool, err error) {
	return t.Clients().TransactionByHash(ctx, txHash)
}

func (t *Client) TransactionReceipt(ctx context.Context, txHash common.Hash) (*types.Receipt, error) {
	return t.Clients().TransactionReceipt(ctx, txHash)
}

func (t *Client) SendTransaction(ctx context.Context, tx *types.Transaction) error {
	return t.Clients().SendTransaction(ctx, tx)
}

func (t *Client) ChainID(ctx context.Context) (*big.Int, error) {
	if t.chainId != nil {
		return t.chainId, nil
	}
	chainId, err := t.Clients().ChainID(ctx)
	if err != nil {
		return nil, err
	}
	t.chainId = chainId
	return chainId, nil
}
