package main

import (
	"context"
	"log"
	"math/rand"
	"time"

	"github.com/proximax-storage/go-xpx-chain-sdk/sdk"
)

const (
	baseURL    = "http://localhost:8000"
	privateKey = "3B9670B5CB19C893694FC49B461CE489BF9588BE16DBE8DC29CF06338133DEE6"
)

func main() {
	// client := createNewClient()
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

	signedMosaicDefinitionTransaction, err := account1.Sign(mosaicDefinitionTransaction)
	if err != nil {
		log.Fatalf("Sign finished with err: %s", err)
	}
	_, err = client.Transaction.Announce(context.Background(), signedMosaicDefinitionTransaction)
	if err != nil {
		log.Fatalf("Transaction.Announce finished with err: %s", err)
	}

	account2, err := client.NewAccount()
	if err != nil {
		log.Fatalf("NewAccount finished with err: %s", err)
	}

	transferTransaction, err := client.NewTransferTransaction(
		sdk.NewDeadline(time.Hour*1),
		sdk.NewAddress(account2.Address.Address, client.NetworkType()),
		[]*sdk.Mosaic{sdk.Xpx(10000000)},
		sdk.NewPlainMessage("empty"),
	)
	if err != nil {
		log.Printf("NewTransferTransaction finished with err: %s", err)
	}

	signedTransferTransaction, err := account2.Sign(transferTransaction)
	if err != nil {
		log.Fatalf("Sign finished with err: %s", err)
	}
	_, err = client.Transaction.Announce(context.Background(), signedTransferTransaction)
	if err != nil {
		log.Fatalf("Transaction.Announce finished with err: %s", err)
	}
}

func generateNonce() (nonce uint32) {
	random := rand.New(rand.NewSource(time.Now().UTC().UnixNano()))
	nonce = random.Uint32()
	return
}
