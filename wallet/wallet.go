package wallet

import (
	"context"
	"crypto/ecdsa"
	"math/big"

	"github.com/shopspring/decimal"
)

type Wallet interface {
	CreateWallet(ctx context.Context) (*ecdsa.PrivateKey, string, error)
	SendTransaction(ctx context.Context, pk, targetAddr string, val decimal.Decimal) error
	GetAccountBalance(ctx context.Context, accountAddr string, blockNumber *big.Int) (decimal.Decimal, error)
	IsValidAddress(ctx context.Context, accountAddr string) bool
}
