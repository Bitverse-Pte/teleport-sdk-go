# Teleport Go SDK

The Teleport GO SDK provides a wrapper for clients to access grpc endpoints exposed by Teleport node.

**WARNING**: This branch is under active development, all subject to potential future change without notification and not ready for production use.

**Note**: Requires Go 1.17+

## Usage

### Init Client

```go
import sdk "github.com/teleport-network/teleport-sdk-go/client"

client, err := sdk.NewClient(GrpcUrl).WithChainId(ChainId)
```

The above is the basic way to initialize a client with no keyring imported. We can use it to call read-only endpoints, like the below

```go
res, err := client.BankQuery.Balance(context.Background(), &types.QueryBalanceRequest{Address: "teleport1qz4xxmn73s8tkttqkw396vklcanl5nzkappyzy", Denom: "atele"})
```

However, mostly we need to broadcast transactions to the Teleport, then we have to register the keyring, and import the sender account, like the below

```go
import sdk "github.com/teleport-network/teleport-sdk-go/client"

client, err := sdk.NewClient(GrpcUrl).WithChainId(ChainId)

// new a keyring 
kr :=  keyring.NewInMemory(c.GetCtx().KeyringOptions...)
// import the account with mnemonic
err = c.WithKeyring(kr).ImportMnemonic(testAcc1.name, testAcc1.mnemonic)
```

With the keyring installed, we can use it to send a transaction to the grpc server.

```go
// build the send message
msg := types.MsgSend{
    FromAddress: testAcc1.addr,
    ToAddress:   "teleport199l57ddd3jepsu3rjen5snyd5x58y2qv9ydpja",
    Amount:      sdk.NewCoins(sdk.NewCoin("atele", sdk.NewInt(10000000))),
}

res1, err := client.Send(msg, func(txf sdktx.Factory) sdktx.Factory {
    return txf.WithFees("100atele") // either fee or gasPrice has to be provided
})
```

### Keyring Management

The example above describes the way to create a keyring based on memory. We also can create a particular instance of keyring

```go
import ( 
    "github.com/cosmos/cosmos-sdk/crypto/keyring"
    "github.com/tharsis/ethermint/crypto/hd"
)

......
kb, err := keyring.New(
    "teleport",         // it is a generic service name that is used by backends that support the concept
    keyringBackend,  // keyring backend, available backends are "os", "file", "kwallet", "memory", "pass", "test".
    rootDir,         // is the directory that keyring files are stored in
    inBuf,           // it is a stdin reader to read the password user needs to input(only os/file backend needed)
    hd.EthSecp256k1Option(), // it defines a function keys options for the ethereum Secp256k1 curve. It is a MUST to our Teleport chain
)
```

We support several keyring backends, as the below

```go
// Backend options for Keyring
const (
    BackendFile    = "file"
    BackendOS      = "os"
    BackendKWallet = "kwallet"
    BackendPass    = "pass"
    BackendTest    = "test"
    BackendMemory  = "memory"
)
```

### Query Endpoints

The Teleport Go SDK client imports grpc query clients from several modules of `teleport`, `cosmos-sdk` to support access the chain data.
It supports for querying from

- auth
- bank
- gov
- xibc
- tmservice

### Broadcast Endpoint

The Teleport Go SDK imports tx service to broadcast transactions. Besides, it wraps various transaction types for clients to submit the transactions. It includes the transaction messages of

- bank
- xibc
- gov

The details please refer to `client` package

## Advanced Usage

### Tx Factory Configuration

If we want to submit the transaction with particular settings, then the `Option` function is what we need. The SDK provides the `Option` as one of the arguments in each transaction sending method to enrich the tx with specific parameters.

```go
type Option func(txf sdktx.Factory) sdktx.Factory
```

For instance, we can send a transfer message like the following:

```go
res1, err := client.Send(msg, func(txf sdktx.Factory) sdktx.Factory {
    return txf.WithFees("100atele").                       // either fee or gasPrice has to be provided
        WithGas(100000).                                // gas limit
        WithGasAdjustment(1.1).                         // default 1
        WithMemo("testMemo").                           // memo 
        WithKeybase(keyring).                           // keyring
        WithSignMode(signing.SignMode_SIGN_MODE_DIRECT) // side mode
})
```

### Account Cache

Once the `client` is initialized, the account cache is enabled by default, which means each time when we build the tx, the sequence is acquired from the account cache. And the sequence will be increased automatically when the tx is successfully submitted.
However, if you want to acquire the sequence from the node each time, you can disable the cache by:

```go
client.DisableCache()
```

### Broadcast Mode

The grpc server defines 3 broadcast modes

```go
// BroadcastBlock defines a tx broadcasting mode where the client waits for
// the tx to be committed in a block.
BroadcastBlock = "block"
// BroadcastSync defines a tx broadcasting mode where the client waits for
// a CheckTx execution response only.
BroadcastSync = "sync"
// BroadcastAsync defines a tx broadcasting mode where the client returns
// immediately.
BroadcastAsync = "async"
```

The client is initialized with `BroadcastSync` mode by default, and can be changed by:

```go
client.WithBroadcastMode("sync")
```
