package client

import (
	"errors"

	sdkclient "github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/tharsis/ethermint/crypto/hd"
	"github.com/tharsis/ethermint/encoding"

	"github.com/teleport-network/teleport-sdk-go/common"
	"github.com/teleport-network/teleport-sdk-go/grpc"
	"github.com/teleport-network/teleport-sdk-go/types"
	"github.com/teleport-network/teleport/app"
)

func init() {
	sdk.GetConfig().SetBech32PrefixForAccount("teleport", "teleportpub")
	sdk.GetConfig().SetBech32PrefixForValidator("teleportvaloper", "teleportvaloperpub")
	sdk.GetConfig().SetBech32PrefixForConsensusNode("teleportvalcons", "teleportvalconspub")
	sdk.GetConfig().SetCoinType(60)
	sdk.GetConfig().SetPurpose(44)
	sdk.GetConfig().Seal()
}

type TeleportClient struct {
	grpc.GClient
	ctx sdkclient.Context

	accountRetriever *types.AccountRetriever
}

func NewClient(url string, chainId string) (*TeleportClient, error) {
	if len(url) == 0 {
		return nil, errors.New("url can not be empty")
	}
	if len(chainId) == 0 {
		return nil, errors.New("chainId can not be empty")
	}
	encodingConfig := encoding.MakeConfig(app.ModuleBasics)
	grpcClient, err := grpc.NewGRPCClient(url)
	if err != nil {
		return nil, err
	}

	accountCache := common.NewCache(1000, true)
	ctx := sdkclient.Context{}.
		WithCodec(encodingConfig.Marshaler).
		WithInterfaceRegistry(encodingConfig.InterfaceRegistry).
		WithTxConfig(encodingConfig.TxConfig).
		WithLegacyAmino(encodingConfig.Amino).
		WithBroadcastMode(flags.BroadcastSync).
		WithKeyringOptions(hd.EthSecp256k1Option()).
		WithChainID(chainId)

	return &TeleportClient{
		ctx:              ctx,
		GClient:          grpcClient,
		accountRetriever: &types.AccountRetriever{QueryClient: grpcClient, Cache: accountCache},
	}, nil
}

func (client *TeleportClient) WithAccountRetrieverCache(cache common.Cache) *TeleportClient {
	client.accountRetriever.Cache = cache
	return client
}

func (client *TeleportClient) GetAccountRetriever() *types.AccountRetriever {
	return client.accountRetriever
}

func (client *TeleportClient) WithChainId(chainId string) *TeleportClient {
	client.ctx = client.ctx.WithChainID(chainId)
	return client
}

func (client *TeleportClient) WithKeyring(k keyring.Keyring) *TeleportClient {
	client.ctx = client.ctx.WithKeyring(k)
	return client
}

func (client *TeleportClient) WithBroadcastMode(mode string) *TeleportClient {
	client.ctx = client.ctx.WithBroadcastMode(mode)
	return client
}

func (client *TeleportClient) DisableCache() {
	client.accountRetriever.Cache.Disable()
}

func (client *TeleportClient) EnableCache() {
	client.accountRetriever.Cache.Enable()
}

func (client *TeleportClient) Key(name string) (string, error) {
	if client.ctx.Keyring == nil {
		return "", errors.New("no keyring found, please add keyring first")
	}
	info, err := client.ctx.Keyring.Key(name)
	if err != nil {
		return "", err
	}
	return info.GetAddress().String(), nil
}

func (client *TeleportClient) ImportKey(name, armor, passphrase string) error {
	if client.ctx.Keyring == nil {
		return errors.New("no keyring found, please add keyring first")
	}
	return client.ctx.Keyring.ImportPrivKey(name, armor, passphrase)
}

func (client *TeleportClient) ImportMnemonic(name, mnemonic string) error {
	if client.ctx.Keyring == nil {
		return errors.New("no keyring found, please add keyring first")
	}

	_, err := client.ctx.Keyring.NewAccount(name, mnemonic, "", sdk.GetConfig().GetFullBIP44Path(), hd.EthSecp256k1)
	return err
}

func (client *TeleportClient) GetCtx() sdkclient.Context {
	return client.ctx
}
