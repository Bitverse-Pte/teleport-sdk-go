package grpc

import (
	"context"
	"fmt"
	"testing"

	"github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/stretchr/testify/require"
)

const GrpcUrl = "grpc0.testnet.teleport.network:443"

func TestQueryBalance(t *testing.T) {
	c, err := NewGRPCClientWithTLSDefault(GrpcUrl)
	require.NoError(t, err)
	res, err := c.BankQuery.Balance(context.Background(), &types.QueryBalanceRequest{Address: "teleport1r60jksyacp3cstz3q5l83suyhtfmm3cautjs68", Denom: "atele"})
	require.NoError(t, err)
	fmt.Println(res.String())
}
