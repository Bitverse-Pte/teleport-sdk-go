package grpc

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/cosmos/cosmos-sdk/x/bank/types"
)

const GrpcUrl = "3.0.202.217:9090"

func TestQueryBalance(t *testing.T) {
	c, err := NewGRPCClient(GrpcUrl)
	assert.NoError(t, err)
	res, _ := c.BankQuery.Balance(context.Background(), &types.QueryBalanceRequest{Address: "teleport1qz4xxmn73s8tkttqkw396vklcanl5nzkappyzy", Denom: "atele"})
	fmt.Println(res.String())
}
