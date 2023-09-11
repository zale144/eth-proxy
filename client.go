package main

import (
	"context"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"go.uber.org/zap"
)

// EthClient is a wrapper around the Ethereum client
type EthClient struct {
	geth    *ethclient.Client
	retries int
	delay   time.Duration
	log     *zap.Logger
}

// NewClient creates a new Ethereum client
func NewClient(networkURL string, retries, retryDelay int, log *zap.Logger) (EthClient, error) {
	geth, err := ethclient.Dial(networkURL)
	if err != nil {
		return EthClient{}, fmt.Errorf("failed to connect to Ethereum network: %w", err)
	}
	return EthClient{
		geth:    geth,
		retries: retries,
		delay:   time.Duration(retryDelay) * time.Second,
		log:     log,
	}, nil
}

// GetBalance fetches the balance of an account from the Ethereum network with retries
func (c EthClient) GetBalance(ctx context.Context, account common.Address) (*big.Int, error) {
	f := func() (*big.Int, error) {
		balance, err := c.geth.BalanceAt(ctx, account, nil)
		if err != nil {
			return nil, fmt.Errorf("failed to get balance: %w", err)
		}
		return balance, nil
	}

	balance, err := retry(f, c.retries, c.delay, c.log)
	if err != nil {
		return nil, fmt.Errorf("failed to get balance: %w", err)
	}

	return balance, nil
}

// BlockNumber fetches the latest block number from the Ethereum network
func (c EthClient) BlockNumber(ctx context.Context) (uint64, error) {
	block, err := c.geth.BlockNumber(ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to get block number: %w", err)
	}
	return block, nil
}

// simple retry logic
func retry[T any](f func() (T, error), retries int, delay time.Duration, log *zap.Logger) (result T, err error) {
	for i := 0; i < retries; i++ {
		result, err = f()
		if err == nil {
			return
		}

		log.Error("Retrying...", zap.Error(err))

		time.Sleep(delay)
		delay *= 2
	}
	return
}
