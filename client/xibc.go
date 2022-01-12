package client

import (
	"github.com/cosmos/cosmos-sdk/types/tx"

	clienttypes "github.com/teleport-network/teleport/x/xibc/core/client/types"
	packettypes "github.com/teleport-network/teleport/x/xibc/core/packet/types"
)

func (client *TeleportClient) UpdateClient(msg clienttypes.MsgUpdateClient, options ...Option) (*tx.BroadcastTxResponse, error) {
	txf, err := Prepare(client, msg.GetSigners()[0], &msg, options...)
	if err != nil {
		return nil, err
	}
	return client.Broadcast(txf, &msg)
}

func (client *TeleportClient) RecvPacket(msg packettypes.MsgRecvPacket, options ...Option) (*tx.BroadcastTxResponse, error) {
	txf, err := Prepare(client, msg.GetSigners()[0], &msg, options...)
	if err != nil {
		return nil, err
	}
	return client.Broadcast(txf, &msg)
}

func (client *TeleportClient) Acknowledgement(msg packettypes.MsgAcknowledgement, options ...Option) (*tx.BroadcastTxResponse, error) {
	txf, err := Prepare(client, msg.GetSigners()[0], &msg, options...)
	if err != nil {
		return nil, err
	}
	return client.Broadcast(txf, &msg)
}
