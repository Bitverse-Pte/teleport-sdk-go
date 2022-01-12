package grpc

import (
	grpc1 "github.com/gogo/protobuf/grpc"
	"google.golang.org/grpc"

	"github.com/cosmos/cosmos-sdk/client/grpc/tmservice"
	"github.com/cosmos/cosmos-sdk/types/tx"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"

	abcitypes "github.com/teleport-network/teleport/grpc_abci/types"
	xibcclitypes "github.com/teleport-network/teleport/x/xibc/core/client/types"
	xibcpkttypes "github.com/teleport-network/teleport/x/xibc/core/packet/types"
)

type GClient struct {
	clientConn      grpc1.ClientConn
	AuthQuery       authtypes.QueryClient
	BankQuery       banktypes.QueryClient
	GovQuery        govtypes.QueryClient
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
	clientConn, err := grpc.Dial(url, dialOpts...)
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
		TMServiceQuery:  tmservice.NewServiceClient(clientConn),
		TxClient:        tx.NewServiceClient(clientConn),
	}, nil
}
