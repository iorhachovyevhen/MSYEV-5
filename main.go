package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/proximax-storage/go-xpx-chain-sdk/sdk"
)

const (
	baseURL    = "http://localhost:3000"
	privateKey = "EFFEE8FB10A0141C38643566C86C0401FC38F3A768D409272F5F5368BB1A2B25"
)

func main() {
	conf, err := sdk.NewConfig(context.Background(), []string{baseURL})
	if err != nil {
		log.Fatalf("sdk.NewConfig finished with error: %s\n", err)
	}

	client := sdk.NewClient(nil, conf)

	account1, err := client.NewAccountFromPrivateKey(privateKey)
	if err != nil {
		log.Fatalf("NewAccountFromPrivateKey finished with err: %s", err)
	}

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

	signedMosaicDefinitionTransaction, err := signTransaction(account1, mosaicDefinitionTransaction)
	if err != nil {
		log.Fatal(err)
	}

	err = announceTransaction(client, signedMosaicDefinitionTransaction)
	if err != nil {
		log.Fatal(err)
	}

	time.Sleep(time.Second * 20)

	account2, err := client.NewAccount()
	if err != nil {
		log.Fatalf("NewAccount finished with err: %s", err)
	}

	mosaic, err := sdk.NewMosaic(mosaicDefinitionTransaction.MosaicId, 1)
	if err != nil {
		log.Fatalf("NewMosaic finished with err: %s", err)
	}

	transferTransaction, err := client.NewTransferTransaction(
		sdk.NewDeadline(time.Hour*1),
		account2.Address,
		[]*sdk.Mosaic{mosaic},
		sdk.NewPlainMessage("empty"),
	)
	if err != nil {
		log.Fatalf("NewTransferTransaction finished with err: %s", err)
	}

	signedTransferTransaction, err := signTransaction(account1, transferTransaction)
	if err != nil {
		log.Fatal(err)
	}

	err = announceTransaction(client, signedTransferTransaction)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Mosaic: %s\n", mosaic.String())
	fmt.Printf("Mosaic defenition: %s\n", signedMosaicDefinitionTransaction.String())
	fmt.Printf("Transfer: %s\n", signedTransferTransaction.String())
	fmt.Printf("Account1: %s\n", account1.String())
	fmt.Printf("Account2: %s\n", account2.String())

	fmt.Println(client.Account.GetAccountInfo(context.Background(), account2.Address))
}

func signTransaction(account *sdk.Account, transaction sdk.Transaction) (*sdk.SignedTransaction, error) {
	signedTransaction, err := account.Sign(transaction)
	if err != nil {
		err = fmt.Errorf("Sign finished with err: %s", err)
		return nil, err
	}
	return signedTransaction, nil
}

func announceTransaction(client *sdk.Client, transaction *sdk.SignedTransaction) error {
	_, err := client.Transaction.Announce(context.Background(), transaction)
	if err != nil {
		return fmt.Errorf("Transaction.Announce finished with err: %s", err)
	}
	return nil
}

func generateNonce() (nonce uint32) {
	random := rand.New(rand.NewSource(time.Now().UTC().UnixNano()))
	nonce = random.Uint32()
	return
}
