package wallet

import (
	"context"
	"crypto/ecdsa"
	"errors"
	"math/big"
	"regexp"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/shopspring/decimal"
)

var (
	networkURL = "http://127.0.0.1:7545"
)

type ETHWallet struct {
	ethClient *ethclient.Client
}

func NewETHWallet() (Wallet, error) {
	eClient, err := ethclient.Dial(networkURL)
	if err != nil {
		return nil, err
	}

	return &ETHWallet{
		ethClient: eClient,
	}, nil
}

func (e *ETHWallet) CreateWallet(ctx context.Context) (*ecdsa.PrivateKey, string, error) {
	var publicAddress string
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		return privateKey, publicAddress, err
	}

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return privateKey, publicAddress, errors.New("error casting public key to ECDSA")
	}
	publicAddress = crypto.PubkeyToAddress(*publicKeyECDSA).Hex()
	return privateKey, publicAddress, nil
}

func (e *ETHWallet) SendTransaction(ctx context.Context, pk, targetAddr string, val decimal.Decimal) error {
	privateKey, err := crypto.HexToECDSA(pk)
	if err != nil {
		return err
	}

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return errors.New("error casting public key to ECDSA")
	}

	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	nonce, err := e.ethClient.PendingNonceAt(ctx, fromAddress)
	if err != nil {
		return err
	}

	// currently we use suggested gas price by default
	gasLimit := uint64(21000)
	gasPrice, err := e.ethClient.SuggestGasPrice(ctx)
	if err != nil {
		return err
	}

	toAddress := common.HexToAddress(targetAddr)
	tx := types.NewTransaction(nonce, toAddress, val.BigInt(), gasLimit, gasPrice, nil)

	chainID, err := e.ethClient.NetworkID(ctx)
	if err != nil {
		return err
	}

	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
	if err != nil {
		return err
	}

	return e.ethClient.SendTransaction(ctx, signedTx)
}

func (e *ETHWallet) GetAccountBalance(ctx context.Context, accountAddr string, blockNumber *big.Int) (decimal.Decimal, error) {
	balance, err := e.ethClient.BalanceAt(ctx, common.HexToAddress(accountAddr), blockNumber)
	if err != nil {
		return decimal.Zero, err
	}

	return decimal.NewFromBigInt(balance, -18), nil
}

func (e *ETHWallet) IsValidAddress(ctx context.Context, accountAddr string) bool {
	re := regexp.MustCompile("^0x[0-9a-fA-F]{40}$")
	return re.MatchString(accountAddr)
}
