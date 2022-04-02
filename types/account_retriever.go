package types

import (
	"context"
	"sync"

	"github.com/cosmos/cosmos-sdk/client"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"

	"github.com/teleport-network/teleport-sdk-go/common"
	grpcclient "github.com/teleport-network/teleport-sdk-go/grpc"
)

type AccountRetriever struct {
	QueryClient grpcclient.GClient
	Cache       common.Cache
	mu          sync.Mutex
}

func (ar *AccountRetriever) RemoveCache(addr sdk.AccAddress) {
	ar.Cache.Remove(addr.String())
}

func (ar *AccountRetriever) IncreaseSequence(addr sdk.AccAddress) {
	ar.mu.Lock()
	defer ar.mu.Unlock()
	acc, err := ar.Cache.Get(addr.String())
	if err == nil && acc != nil {
		acc, ok := acc.(authtypes.AccountI)
		if ok {
			_ = acc.SetSequence(acc.GetSequence() + 1)
		}
	}
}

func (ar *AccountRetriever) getFromCache(addr sdk.AccAddress) authtypes.AccountI {
	if ar.Cache != nil {
		v, err := ar.Cache.Get(addr.String())
		if err != nil || v == nil {
			return nil
		}
		if v, ok := v.(authtypes.AccountI); ok {
			return v
		}
	}

	return nil
}

func (ar *AccountRetriever) GetAccount(clientCtx client.Context, addr sdk.AccAddress) (authtypes.AccountI, error) {
	if acc := ar.getFromCache(addr); acc != nil {
		return acc, nil
	}

	res, err := ar.QueryClient.AuthQuery.Account(context.Background(), &authtypes.QueryAccountRequest{Address: addr.String()})
	if err != nil {
		return nil, err
	}

	var acc authtypes.AccountI
	if err := clientCtx.InterfaceRegistry.UnpackAny(res.Account, &acc); err != nil {
		return nil, err
	}

	_ = ar.Cache.Set(addr.String(), acc)

	return acc, err
}

func (ar *AccountRetriever) RefreshSequence(clientCtx client.Context) (authtypes.AccountI, error) {
	from := clientCtx.GetFromAddress()

	res, err := ar.QueryClient.AuthQuery.Account(context.Background(), &authtypes.QueryAccountRequest{Address: from.String()})
	if err != nil {
		return nil, err
	}

	var acc authtypes.AccountI
	if err := clientCtx.InterfaceRegistry.UnpackAny(res.Account, &acc); err != nil {
		return nil, err
	}
	err = acc.SetSequence(acc.GetSequence() + 1)
	if err != nil {
		return nil, err
	}
	_ = ar.Cache.Set(from.String(), acc)

	return acc, err
}

// EnsureExists returns an error if no account exists for the given address else nil.
func (ar *AccountRetriever) EnsureExists(clientCtx client.Context, addr sdk.AccAddress) error {
	if _, err := ar.GetAccount(clientCtx, addr); err != nil {
		return err
	}

	return nil
}

// GetAccountNumberSequence returns sequence and account number for the given address.
// It returns an error if the account couldn't be retrieved from the state.
func (ar *AccountRetriever) GetAccountNumberSequence(clientCtx client.Context, addr sdk.AccAddress) (uint64, uint64, error) {
	acc, err := ar.GetAccount(clientCtx, addr)
	if err != nil {
		return 0, 0, err
	}

	return acc.GetAccountNumber(), acc.GetSequence(), nil
}
