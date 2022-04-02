package client

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/avast/retry-go"

	"github.com/cosmos/cosmos-sdk/client"
	sdktx "github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/tx"

	"github.com/teleport-network/teleport-sdk-go/types"
)

type Option func(txf sdktx.Factory) sdktx.Factory

func Prepare(client *TeleportClient, signer sdk.AccAddress, msg sdk.Msg, options ...Option) (sdktx.Factory, error) {
	if err := msg.ValidateBasic(); err != nil {
		return sdktx.Factory{}, err
	}
	if client.ctx.Keyring == nil {
		return sdktx.Factory{}, errors.New("keyring must be imported")
	}
	info, err := client.ctx.Keyring.KeyByAddress(signer)
	if err != nil {
		return sdktx.Factory{}, err
	}
	client.ctx.FromAddress = signer
	client.ctx.FromName = info.GetName()

	txf := sdktx.Factory{}.
		WithChainID(client.ctx.ChainID).
		WithTxConfig(client.ctx.TxConfig).
		WithKeybase(client.ctx.Keyring).
		WithGasAdjustment(1) // default 1, can be changed by 'Option'
	for _, op := range options {
		txf = op(txf)
	}
	if txf.Gas() == 0 {
		txf = txf.WithSimulateAndExecute(true)
	}
	return txf, nil
}

// Broadcast Sign and broadcast to node. It is retryable.
func (client *TeleportClient) Broadcast(txf sdktx.Factory, msgs ...sdk.Msg) (res *tx.BroadcastTxResponse, err error) {
	retryableFunc := func() error {
		txf, err := SetupAccNumberSequence(client.ctx, client.accountRetriever, txf)
		if err != nil {
			return err
		}
		res, err = client.broadcast(txf, msgs...)
		if err == nil && res.TxResponse.Code == 0 {
			client.accountRetriever.IncreaseSequence(client.ctx.FromAddress)
		}
		return err
	}

	retryIfFunc := func(err error) bool {
		return strings.Contains(err.Error(), "account sequence mismatch")
	}

	onRetryFunc := func(n uint, err error) {
		client.accountRetriever.RemoveCache(client.ctx.FromAddress)
	}

	err = retry.Do(
		retryableFunc,
		retry.Attempts(3),
		retry.RetryIf(retryIfFunc),
		retry.OnRetry(onRetryFunc),
	)

	return
}

func (client *TeleportClient) broadcast(txf sdktx.Factory, msgs ...sdk.Msg) (*tx.BroadcastTxResponse, error) {
	if txf.SimulateAndExecute() {
		_, adjusted, err := client.calculateGas(txf, msgs...)
		if err != nil {
			return nil, err
		}

		txf = txf.WithGas(adjusted)
		_, _ = fmt.Fprintf(os.Stderr, "%s\n", sdktx.GasEstimateResponse{GasEstimate: txf.Gas()})
	}
	txBuilder, err := sdktx.BuildUnsignedTx(txf, msgs...)
	if err != nil {
		return nil, err
	}

	err = sdktx.Sign(txf, client.ctx.GetFromName(), txBuilder, true)
	if err != nil {
		return nil, err
	}

	txBytes, err := client.ctx.TxConfig.TxEncoder()(txBuilder.GetTx())
	if err != nil {
		return nil, err
	}

	res, err := client.BroadcastTx(txBytes)

	return res, err
}

// CalculateGas simulation response obtained by the query and the adjusted gas amount
// it is retryable
func (client *TeleportClient) CalculateGas(txf sdktx.Factory, msgs ...sdk.Msg) (res *tx.SimulateResponse, gas uint64, err error) {
	retryableFunc := func() error {
		txf, err = SetupAccNumberSequence(client.ctx, client.accountRetriever, txf)
		if err != nil {
			return err
		}
		res, gas, err = client.calculateGas(txf, msgs...)
		return err
	}

	retryIfFunc := func(err error) bool {
		return strings.Contains(err.Error(), "account sequence mismatch")
	}

	onRetryFunc := func(n uint, err error) {
		client.accountRetriever.RemoveCache(client.ctx.FromAddress)
	}

	err = retry.Do(
		retryableFunc,
		retry.Attempts(3),
		retry.RetryIf(retryIfFunc),
		retry.OnRetry(onRetryFunc),
	)

	return
}

func (client *TeleportClient) calculateGas(txf sdktx.Factory, msgs ...sdk.Msg) (*tx.SimulateResponse, uint64, error) {
	txBytes, err := sdktx.BuildSimTx(txf, msgs...)
	if err != nil {
		return nil, 0, err
	}

	simRes, err := client.Simulate(txBytes)
	if err != nil {
		return nil, 0, err
	}

	return simRes, uint64(txf.GasAdjustment() * float64(simRes.GasInfo.GasUsed)), nil
}

func (client *TeleportClient) Simulate(txBytes []byte) (*tx.SimulateResponse, error) {
	return client.TxClient.Simulate(
		context.Background(),
		&tx.SimulateRequest{TxBytes: txBytes},
	)
}

func (client *TeleportClient) BroadcastTx(txBytes []byte) (*tx.BroadcastTxResponse, error) {
	return client.TxClient.BroadcastTx(
		context.Background(),
		&tx.BroadcastTxRequest{
			TxBytes: txBytes,
			Mode:    convertBroadcastMode(client.ctx.BroadcastMode),
		},
	)
}

func (client *TeleportClient) GetTx(hash string) (*tx.GetTxResponse, error) {
	return client.TxClient.GetTx(context.Background(), &tx.GetTxRequest{Hash: hash})
}

func (client *TeleportClient) GetTxsEvent(req *tx.GetTxsEventRequest) (*tx.GetTxsEventResponse, error) {
	return client.TxClient.GetTxsEvent(context.Background(), req)
}

// to be passed into the clientCtx.
func convertBroadcastMode(mode string) tx.BroadcastMode {
	switch mode {
	case "async":
		return tx.BroadcastMode_BROADCAST_MODE_ASYNC
	case "block":
		return tx.BroadcastMode_BROADCAST_MODE_BLOCK
	case "sync":
		return tx.BroadcastMode_BROADCAST_MODE_SYNC
	default:
		return tx.BroadcastMode_BROADCAST_MODE_UNSPECIFIED
	}
}

// SetupAccNumberSequence ensures the account defined by ctx.GetFromAddress() exists and
// if the account number and/or the account sequence number are zero (not set),
// they will be queried for and set on the provided Factory. A new Factory with
// the updated fields will be returned.
func SetupAccNumberSequence(clientCtx client.Context, accountRetriever *types.AccountRetriever, txf sdktx.Factory) (sdktx.Factory, error) {
	from := clientCtx.GetFromAddress()

	if err := accountRetriever.EnsureExists(clientCtx, from); err != nil {
		return txf, err
	}

	num, seq, err := accountRetriever.GetAccountNumberSequence(clientCtx, from)
	if err != nil {
		return txf, err
	}

	txf = txf.WithAccountNumber(num).WithSequence(seq)

	return txf, nil
}
