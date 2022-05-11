package grpc

import (
	"crypto/tls"

	"github.com/cosmos/cosmos-sdk/client/grpc/tmservice"
	"github.com/cosmos/cosmos-sdk/types/tx"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	abcitypes "github.com/teleport-network/teleport/grpc_abci/types"
	xibcclitypes "github.com/teleport-network/teleport/x/xibc/core/client/types"
	xibcpkttypes "github.com/teleport-network/teleport/x/xibc/core/packet/types"

	grpc1 "github.com/gogo/protobuf/grpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type GClient struct {
	clientConn      grpc1.ClientConn
	AuthQuery       authtypes.QueryClient
	BankQuery       banktypes.QueryClient
	GovQuery        govtypes.QueryClient
	StakingQuery    stakingtypes.QueryClient
	XIBCClientQuery xibcclitypes.QueryClient
	XIBCPacketQuery xibcpkttypes.QueryClient
	ABCIQuery       abcitypes.ABCIQueryClient
	TMServiceQuery  tmservice.ServiceClient
	TxClient        tx.ServiceClient
}

func NewGRPCClient(url string) (GClient, error) {
	dialOpts := []grpc.DialOption{
		grpc.WithInsecure(),
	}
	return buildGRPCClient(url, dialOpts...)
}

func NewGRPCClientWithTLSDefault(url string) (GClient, error) {
	dialOpts := []grpc.DialOption{
		grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{})),
	}
	return buildGRPCClient(url, dialOpts...)
}

func NewGRPCClientWithTLS(url string, c *tls.Config) (GClient, error) {
	dialOpts := []grpc.DialOption{
		grpc.WithTransportCredentials(credentials.NewTLS(c)),
	}
	return buildGRPCClient(url, dialOpts...)
}

func buildGRPCClient(url string, opts ...grpc.DialOption) (GClient, error) {
	clientConn, err := grpc.Dial(url, opts...)
	if err != nil {
		return GClient{}, err
	}
	return GClient{
		clientConn:      clientConn,
		XIBCClientQuery: xibcclitypes.NewQueryClient(clientConn),
		XIBCPacketQuery: xibcpkttypes.NewQueryClient(clientConn),
		ABCIQuery:       abcitypes.NewABCIQueryClient(clientConn),
		BankQuery:       banktypes.NewQueryClient(clientConn),
		AuthQuery:       authtypes.NewQueryClient(clientConn),
		StakingQuery:    stakingtypes.NewQueryClient(clientConn),
		TMServiceQuery:  tmservice.NewServiceClient(clientConn),
		TxClient:        tx.NewServiceClient(clientConn),
	}, nil
}
