package client

import (
	"github.com/cosmos/cosmos-sdk/types/tx"
	"github.com/cosmos/cosmos-sdk/x/gov/types"
)

func (client *TeleportClient) SubmitProposal(msg types.MsgSubmitProposal, options ...Option) (*tx.BroadcastTxResponse, error) {
	txf, err := Prepare(client, msg.GetSigners()[0], &msg, options...)
	if err != nil {
		return nil, err
	}
	return client.Broadcast(txf, &msg)
}

func (client *TeleportClient) Deposit(msg types.MsgDeposit, options ...Option) (*tx.BroadcastTxResponse, error) {
	txf, err := Prepare(client, msg.GetSigners()[0], &msg, options...)
	if err != nil {
		return nil, err
	}
	return client.Broadcast(txf, &msg)
}

func (client *TeleportClient) Vote(msg types.MsgVote, options ...Option) (*tx.BroadcastTxResponse, error) {
	txf, err := Prepare(client, msg.GetSigners()[0], &msg, options...)
	if err != nil {
		return nil, err
	}
	return client.Broadcast(txf, &msg)
}

func (client *TeleportClient) VoteWeighted(msg types.MsgVoteWeighted, options ...Option) (*tx.BroadcastTxResponse, error) {
	txf, err := Prepare(client, msg.GetSigners()[0], &msg, options...)
	if err != nil {
		return nil, err
	}
	return client.Broadcast(txf, &msg)
}
