package integration

import (
	"github.com/cosmos/cosmos-sdk/crypto/keyring"

	"github.com/teleport-network/teleport-sdk-go/client"
)

// editable settings for test
const (
	GrpcUrl = "3.0.202.217:9090"
	ChainId = "teleport_9000-1"
)

var (
	testAcc1 = testAcc{
		name:     "acc1",
		addr:     "teleport1n855w0jwp2n04dtxvuu3zedxdlxz5c5u85mxv3",
		mnemonic: "donate broccoli change around paper fetch rifle matrix guide pioneer catalog blur okay absorb rude much chef assist virtual turn exhaust wing output scene",
	}

	// testAcc2 = testAcc{
	// 	name:     "acc2",
	// 	addr:     "teleport199l57ddd3jepsu3rjen5snyd5x58y2qv9ydpja",
	// 	mnemonic: "roast advice exit enrich august super talk dash expire flag attract glare release moral perfect depth mountain keep lake sudden wing have electric wild",
	// }
)

type testAcc struct {
	name     string
	addr     string
	mnemonic string
}

func newClient() (*client.TeleportClient, error) {
	c, err := client.NewClient(GrpcUrl, ChainId)
	if err != nil {
		return nil, err
	}
	if err = c.WithKeyring(
		keyring.NewInMemory(c.GetCtx().KeyringOptions...),
	).ImportMnemonic(
		testAcc1.name, testAcc1.mnemonic,
	); err != nil {
		return nil, err
	}
	return c, nil
}
