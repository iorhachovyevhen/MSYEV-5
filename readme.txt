Goal: create new mosaic and transfer it from one newly account to another.

Before creating new account create:
    - new confing (use NewConfig(params)):
        conf, err := sdk.NewConfig(context.Background(), []string{baseURL})
        if err != nil {
            log.Fatalf("sdk.NewConfig finished with error: %s\n", err)
        }
    - new client:
        client := sdk.NewClient(nil, conf)
After config and client have been created create first new account:
    account1, err := client.NewAccountFromPrivateKey(privateKey)
	if err != nil {
		log.Fatalf("NewAccountFromPrivateKey finished with err: %s", err)
	}

Creating mosaic

Before creating new mosaic you need publish it definition to the newtwork using:
    mosaicDefinitionTransaction, err := client.NewMosaicDefinitionTransaction(
		sdk.NewDeadline(time.Hour*1),
		generateNonce(),
		account1.PublicAccount.PublicKey,
		sdk.NewMosaicProperties(
			true,
			true,
			3,
			sdk.Duration(10000),
		),
	)
	if err != nil {
		log.Fatalf("NewMosaicDefinitionTransaction finished with err : %s", err)
	}

And you should sign it using created account using func signTransaction:
    func signTransaction(account *sdk.Account, transaction sdk.Transaction) (*sdk.SignedTransaction, error) {
        signedTransaction, err := account.Sign(transaction)
        if err != nil {
            err = fmt.Errorf("Sign finished with err: %s", err)
            return nil, err
        }
        return signedTransaction, nil
    }

After you should announce this mosaic using func announceTransaction:
    func announceTransaction(client *sdk.Client, transaction *sdk.SignedTransaction) error {
        _, err := client.Transaction.Announce(context.Background(), transaction)
        if err != nil {
            return fmt.Errorf("Transaction.Announce finished with err: %s", err)
        }
        return nil
    }

Creating new mosaic in the network takes 15s.

When mosaic has been created in the network you should create new local mosaic:
	mosaic, err := sdk.NewMosaic(mosaicDefinitionTransaction.MosaicId, 1)
	if err != nil {
		log.Fatalf("NewMosaic finished with err: %s", err)
	}

And create new trancfer transaction:
	transferTransaction, err := client.NewTransferTransaction(
		sdk.NewDeadline(time.Hour*1),
		account2.Address,
		[]*sdk.Mosaic{mosaic},
		sdk.NewPlainMessage("empty"),
	)
	if err != nil {
		log.Fatalf("NewTransferTransaction finished with err: %s", err)
	}
where account2 is:
    account2, err := client.NewAccount()
        if err != nil {
            log.Fatalf("NewAccount finished with err: %s", err)
	}


As before you should sign this transaction and after announce.