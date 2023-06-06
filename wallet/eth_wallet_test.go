package wallet

import (
	"context"
	"math/big"
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

var (
	walletAccountPK = "e800a26714bf16f22f061868b9e39ab2a2e700b13012178b595da4d7b2c8fc5d"
)

func Test_EthWallet(t *testing.T) {
	ctx := context.Background()
	ethWallet, err := NewETHWallet()
	assert.Empty(t, err)

	_, pbAddress, err := ethWallet.CreateWallet(ctx)
	assert.Empty(t, err)
	assert.NotEmpty(t, pbAddress)

	err = ethWallet.SendTransaction(ctx, walletAccountPK, pbAddress, decimal.NewFromBigInt(big.NewInt(15), 18))
	assert.Empty(t, err)

	balance, err := ethWallet.GetAccountBalance(ctx, pbAddress, nil)
	assert.Empty(t, err)
	assert.True(t, balance.Equal(decimal.NewFromFloat(15)))
}
