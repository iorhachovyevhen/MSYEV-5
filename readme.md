Instruction for creating new mosaic and transfer it from one newly account to another.

Before creating new account create:
- new confing:
```go
conf, err := sdk.NewConfig(context.Background(), []string{baseURL})
if err != nil {
    log.Fatalf("sdk.NewConfig finished with error: %s\n", err)
}
```
where `baseURl` is your server
- new client:
```go 
client := sdk.NewClient(nil, conf)
```
where the func `NewClient` takes two params:
- http.Client
- Config

If http.Client is nil func `NewClient` returned `DefaultClient`.

After config and client have been created you should create the first new account:
```go
account1, err := client.NewAccountFromPrivateKey(privateKey)
if err != nil {
    log.Fatalf("NewAccountFromPrivateKey finished with err: %s", err)
}
```
where func `NewAccountFromPrivateKey takes` takes param `privateKey` - private key of future account

Before creating new mosaic you need publish it definition to the newtwork using:
```go
mosaicDefinitionTransaction, err := client.NewMosaicDefinitionTransaction(
    sdk.NewDeadline(time.Hour*1),     // deadline
    generateNonce(),                  // nonce
    account1.PublicAccount.PublicKey, // ownerPublicKey
    sdk.NewMosaicProperties(          // mosaicprops
        true,                         // supplyMutable
        true,                         // transferable
        3,                            // divisibility
        sdk.Duration(10000),          // duration
    ),
)
if err != nil {
    log.Fatalf("NewMosaicDefinitionTransaction finished with err : %s", err)
}
```
where `NewMosaicDefinitionTransaction` takes next params:
- deadline (Timestamp) - maximum time for including the transaction in blockchain;
- nonce (uint32) - TODO;
- ownerPublicKey - public key of public account, that create mosaic definition transaction;
- mosaicProps - creating new mosaic properties using func `NewMosaicProperties` with next params:
    - supplyMutable (bool) - determines whether the creator is allowed in future to decrease the supply within the limits of mosaics owned or it a immutable supply.
    - transferable (bool) - determines whether the mosaic can be transferred in other accounts or only the creator can be recepient after first the first transfer.
    - divisibility (uint8) - determines which decimal number to divide the mosaic into. For example, if value=3 than a mosaic can be divided into smallest parts of 0.001 mosaics. Range of 0 and 6.
    - duration (baseInt64) - the number of confirmed blocks we would like to rent our namespace for. Should be inferior or equal to namespace duration.

And you should sign it using created account:
```go
func signTransaction(account *sdk.Account, transaction sdk.Transaction) (*sdk.SignedTransaction, error) {
    signedTransaction, err := account.Sign(transaction)
    if err != nil {
        err = fmt.Errorf("Sign finished with err: %s", err)
        return nil, err
    }
    return signedTransaction, nil
}
```
where func `signTransaction` takes next params:
- account (Account struct) - previously created account;
- transaction (Transaction interface) - transaction to be signed

After you should announce this mosaic using func `announceTransaction`:
```go
func announceTransaction(client *sdk.Client, transaction *sdk.SignedTransaction) error {
    _, err := client.Transaction.Announce(context.Background(), transaction)
    if err != nil {
        return fmt.Errorf("Transaction.Announce finished with err: %s", err)
    }
    return nil
}
```
where func `announceTransaction` takes next params:
- client (Client struct) - created early client
- transaction (SignedTransaction struct) - signed early transaction

Creating new mosaic in the network takes 15s.

When mosaic has been created in the network you should create new local mosaic:
```go
mosaic, err := sdk.NewMosaic(mosaicDefinitionTransaction.MosaicId, 1)
if err != nil {
    log.Fatalf("NewMosaic finished with err: %s", err)
}
```
where func `NewMosaic` takes next params:
- assetId (AssetId struct) - blockchain identifier
- amount (Amount struct) - amount of blocks

And create new transfer transaction:
```go
transferTransaction, err := client.NewTransferTransaction(
    sdk.NewDeadline(time.Hour*1), // deadline
    account2.Address,             // recipient
    []*sdk.Mosaic{mosaic},        // mosaics
    sdk.NewPlainMessage("empty"), // message
)
if err != nil {
    log.Fatalf("NewTransferTransaction finished with err: %s", err)
}
```
where `NewTransferTransaction` takes next params:
- deadline (Timestamp) - maximum time for including the transaction in blockchain;
- recipient (*Address) - address of recipient
- mosaics ([]*Mosaic) - array of mosaics to transfer
- message (Message interface) - attached message to the transaction (max size 1024 characters)

and where `account2` is:
```go
account2, err := client.NewAccount()
    if err != nil {
        log.Fatalf("NewAccount finished with err: %s", err)
}
```

As before you should sign this transaction and after announce.