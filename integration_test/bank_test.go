package integration

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	sdktx "github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/bank/types"
)

func TestMsgSend(t *testing.T) {
	client, err := newClient()
	assert.NoError(t, err)

	msg := types.MsgSend{
		FromAddress: testAcc1.addr,
		ToAddress:   "teleport199l57ddd3jepsu3rjen5snyd5x58y2qv9ydpja",
		Amount:      sdk.NewCoins(sdk.NewCoin("atele", sdk.NewInt(10000000))),
	}

	res1, err := client.Send(msg, func(txf sdktx.Factory) sdktx.Factory {
		return txf.WithFees("100atele")
	})
	assert.NoError(t, err)
	assert.EqualValues(t, 0, res1.TxResponse.Code)
	fmt.Println(res1.String())

	res2, err := client.Send(msg, func(txf sdktx.Factory) sdktx.Factory {
		return txf.WithFees("100atele")
	})
	assert.NoError(t, err)
	assert.EqualValues(t, 0, res2.TxResponse.Code)
	fmt.Println(res2.String())

	time.Sleep(5 * time.Second)
	txRes, err := client.GetTx(res1.TxResponse.TxHash)
	assert.NoError(t, err)
	fmt.Println(txRes.TxResponse.String())

	txRes, err = client.GetTx(res2.TxResponse.TxHash)
	assert.NoError(t, err)
	fmt.Println(txRes.TxResponse.String())
}
