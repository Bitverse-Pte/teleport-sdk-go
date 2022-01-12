package client

import (
	"github.com/cosmos/cosmos-sdk/types/tx"
	"github.com/cosmos/cosmos-sdk/x/bank/types"
)

func (client *TeleportClient) Send(msg types.MsgSend, options ...Option) (*tx.BroadcastTxResponse, error) {
	txf, err := Prepare(client, msg.GetSigners()[0], &msg, options...)
	if err != nil {
		return nil, err
	}
	return client.Broadcast(txf, &msg)
}
