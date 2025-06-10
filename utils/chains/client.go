package chains

import (
	"context"
	"io"
	"log"
	"math/big"
	"os"
	"strings"
	"sync/atomic"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/ethclient/simulated"
	"github.com/pkg/errors"
	"github.com/spf13/cast"
)

var (
	ENABLE_CHAIN_CLIENT_DEBUG = cast.ToBool(getEnv("ENABLE_CHAIN_CLIENT_DEBUG", false))
	MAX_RETRY_COUNT           = cast.ToInt(getEnv("MAX_RETRY_COUNT", 3))
	INDEX                     = atomic.Int64{}
)

type MockClient interface {
	simulated.Client
}

type Client struct {
	clients []*ethclient.Client
	rpcs    []string
	chainId *big.Int
}

func getEnv(key string, defualt interface{}) interface{} {
	v := os.Getenv(key)
	if v == "" {
		return defualt
	}
	return v
}

func try(f func() error) {
	var tryCount = 0
	for {
		err := f()
		if err == nil {
			return
		}
		tryCount++
		if tryCount > MAX_RETRY_COUNT {
			return
		}
	}
}

func needTry(err error) bool {
	if err == nil {
		return false
	}
	return strings.Contains(err.Error(), "connection refused") ||
		strings.Contains(err.Error(), "no such host") ||
		strings.Contains(err.Error(), "i/o timeout") ||
		errors.Is(err, io.EOF)

}

func NewTryClient(ctx context.Context, endpoints []string) (*Client, error) {
	t := &Client{}
	for _, v := range endpoints {
		client, err := ethclient.DialContext(ctx, v)
		if err != nil {
			log.Printf("dial %s error: %s\n", v, err)
			continue
		}
		t.clients = append(t.clients, client)
		t.rpcs = append(t.rpcs, v)
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
	index := INDEX.Add(1) % int64(len(t.clients))
	if ENABLE_CHAIN_CLIENT_DEBUG {
		log.Printf("use client: %s\n", t.rpcs[index])
	}
	return t.clients[index]
}

func (t *Client) BlockNumber(ctx context.Context) (blockNumber uint64, err error) {
	try(func() error {
		blockNumber, err = t.Clients().BlockNumber(ctx)
		if needTry(err) {
			return err
		}
		return nil
	})
	return
}

func (t *Client) BlockByHash(ctx context.Context, hash common.Hash) (block *types.Block, err error) {
	try(func() error {
		block, err = t.Clients().BlockByHash(ctx, hash)
		if needTry(err) {
			return err
		}
		return nil
	})
	return
}

func (t *Client) BlockByNumber(ctx context.Context, number *big.Int) (block *types.Block, err error) {
	try(func() error {
		block, err = t.Clients().BlockByNumber(ctx, number)
		if needTry(err) {
			return err
		}
		return nil
	})
	return
}

func (t *Client) HeaderByHash(ctx context.Context, hash common.Hash) (header *types.Header, err error) {
	try(func() error {
		header, err = t.Clients().HeaderByHash(ctx, hash)
		if needTry(err) {
			return err
		}
		return nil
	})
	return
}

func (t *Client) HeaderByNumber(ctx context.Context, number *big.Int) (header *types.Header, err error) {
	try(func() error {
		header, err = t.Clients().HeaderByNumber(ctx, number)
		if needTry(err) {
			return err
		}
		return nil
	})
	return
}

func (t *Client) TransactionCount(ctx context.Context, blockHash common.Hash) (transactionCount uint, err error) {
	try(func() error {
		transactionCount, err = t.Clients().TransactionCount(ctx, blockHash)
		if needTry(err) {
			return err
		}
		return nil
	})
	return
}

func (t *Client) TransactionInBlock(ctx context.Context, blockHash common.Hash, index uint) (transaction *types.Transaction, err error) {
	try(func() error {
		transaction, err = t.Clients().TransactionInBlock(ctx, blockHash, index)
		if needTry(err) {
			return err
		}
		return nil
	})
	return
}

func (t *Client) SubscribeNewHead(ctx context.Context, ch chan<- *types.Header) (ethereum.Subscription, error) {
	return t.Clients().SubscribeNewHead(ctx, ch)
}

func (t *Client) BalanceAt(ctx context.Context, account common.Address, blockNumber *big.Int) (balance *big.Int, err error) {
	try(func() error {
		balance, err = t.Clients().BalanceAt(ctx, account, blockNumber)
		if needTry(err) {
			return err
		}
		return nil
	})
	return
}

func (t *Client) StorageAt(ctx context.Context, account common.Address, key common.Hash, blockNumber *big.Int) (storage []byte, err error) {
	try(func() error {
		storage, err = t.Clients().StorageAt(ctx, account, key, blockNumber)
		if needTry(err) {
			return err
		}
		return nil
	})
	return
}

func (t *Client) CodeAt(ctx context.Context, account common.Address, blockNumber *big.Int) (code []byte, err error) {
	try(func() error {
		code, err = t.Clients().CodeAt(ctx, account, blockNumber)
		if needTry(err) {
			return err
		}
		return nil
	})
	return
}

func (t *Client) NonceAt(ctx context.Context, account common.Address, blockNumber *big.Int) (nonce uint64, err error) {
	try(func() error {
		nonce, err = t.Clients().NonceAt(ctx, account, blockNumber)
		if needTry(err) {
			return err
		}
		return nil
	})
	return
}

func (t *Client) CallContract(ctx context.Context, call ethereum.CallMsg, blockNumber *big.Int) (result []byte, err error) {
	try(func() error {
		result, err = t.Clients().CallContract(ctx, call, blockNumber)
		if needTry(err) {
			return err
		}
		return nil
	})
	return
}

func (t *Client) EstimateGas(ctx context.Context, call ethereum.CallMsg) (gas uint64, err error) {
	try(func() error {
		gas, err = t.Clients().EstimateGas(ctx, call)
		if needTry(err) {
			return err
		}
		return nil
	})
	return
}

func (t *Client) SuggestGasPrice(ctx context.Context) (gasPrice *big.Int, err error) {
	try(func() error {
		gasPrice, err = t.Clients().SuggestGasPrice(ctx)
		if needTry(err) {
			return err
		}
		return nil
	})
	return
}

func (t *Client) SuggestGasTipCap(ctx context.Context) (gasTipCap *big.Int, err error) {
	try(func() error {
		gasTipCap, err = t.Clients().SuggestGasTipCap(ctx)
		if needTry(err) {
			return err
		}
		return nil
	})
	return
}

func (t *Client) FeeHistory(ctx context.Context, blockCount uint64, lastBlock *big.Int, rewardPercentiles []float64) (feeHistory *ethereum.FeeHistory, err error) {
	try(func() error {
		feeHistory, err = t.Clients().FeeHistory(ctx, blockCount, lastBlock, rewardPercentiles)
		if needTry(err) {
			return err
		}
		return nil
	})
	return
}

func (t *Client) FilterLogs(ctx context.Context, q ethereum.FilterQuery) (logs []types.Log, err error) {
	try(func() error {
		logs, err = t.Clients().FilterLogs(ctx, q)
		if needTry(err) {
			return err
		}
		return nil
	})
	return
}

func (t *Client) SubscribeFilterLogs(ctx context.Context, q ethereum.FilterQuery, ch chan<- types.Log) (ethereum.Subscription, error) {
	return t.Clients().SubscribeFilterLogs(ctx, q, ch)
}

func (t *Client) PendingBalanceAt(ctx context.Context, account common.Address) (balance *big.Int, err error) {
	try(func() error {
		balance, err = t.Clients().PendingBalanceAt(ctx, account)
		if needTry(err) {
			return err
		}
		return nil
	})
	return
}

func (t *Client) PendingStorageAt(ctx context.Context, account common.Address, key common.Hash) (result []byte, err error) {
	try(func() error {
		result, err = t.Clients().PendingStorageAt(ctx, account, key)
		if needTry(err) {
			return err
		}
		return nil
	})
	return
}

func (t *Client) PendingCodeAt(ctx context.Context, account common.Address) ([]byte, error) {
	return t.Clients().PendingCodeAt(ctx, account)
}

func (t *Client) PendingNonceAt(ctx context.Context, account common.Address) (nonce uint64, err error) {
	try(func() error {
		nonce, err = t.Clients().PendingNonceAt(ctx, account)
		if needTry(err) {
			return err
		}
		return nil
	})
	return
}

func (t *Client) PendingTransactionCount(ctx context.Context) (count uint, err error) {
	try(func() error {
		count, err = t.Clients().PendingTransactionCount(ctx)
		if needTry(err) {
			return err
		}
		return nil
	})
	return
}

func (t *Client) PendingCallContract(ctx context.Context, call ethereum.CallMsg) (result []byte, err error) {
	try(func() error {
		result, err = t.Clients().PendingCallContract(ctx, call)
		if needTry(err) {
			return err
		}
		return nil
	})
	return
}

func (t *Client) TransactionByHash(ctx context.Context, txHash common.Hash) (tx *types.Transaction, isPending bool, err error) {
	try(func() error {
		tx, isPending, err = t.Clients().TransactionByHash(ctx, txHash)
		if needTry(err) {
			return err
		}
		return nil
	})
	return
}

func (t *Client) TransactionReceipt(ctx context.Context, txHash common.Hash) (receipt *types.Receipt, err error) {
	try(func() error {
		receipt, err = t.Clients().TransactionReceipt(ctx, txHash)
		if needTry(err) {
			return err
		}
		return nil
	})
	return
}

func (t *Client) SendTransaction(ctx context.Context, tx *types.Transaction) (err error) {
	try(func() error {
		err = t.Clients().SendTransaction(ctx, tx)
		if needTry(err) {
			return err
		}
		return nil
	})
	return
}

func (t *Client) ChainID(ctx context.Context) (chainId *big.Int, err error) {
	if t.chainId != nil {
		return t.chainId, nil
	}
	try(func() error {
		chainId, err = t.Clients().ChainID(ctx)
		if needTry(err) {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	t.chainId = chainId
	return chainId, nil
}
