package lib

import (
	"context"
	"log"

	"github.com/howjmay/go-sui-sdk/sui"
	"github.com/howjmay/go-sui-sdk/sui_types"
	"github.com/howjmay/go-sui-sdk/types"
)

type Publisher struct {
	client  *sui.ImplSuiAPI
	account *sui_types.Account
}

func NewPublisher(client *sui.ImplSuiAPI, account *sui_types.Account) *Publisher {
	return &Publisher{
		client:  client,
		account: account,
	}
}

func (p *Publisher) PublishEvents(ctx context.Context, packageID *sui_types.PackageID) {
	txnBytes, err := p.client.MoveCall(
		ctx,
		p.account.AccountAddress(),
		packageID,
		"eventpub",
		"emit_clock",
		[]string{},
		[]any{},
		nil,
		types.NewSafeSuiBigInt(uint64(100000)),
	)
	if err != nil {
		log.Panic(err)
	}

	signature, err := p.account.SignSecureWithoutEncode(txnBytes.TxBytes.Data(), sui_types.DefaultIntent())
	if err != nil {
		log.Panic(err)
	}

	txnResponse, err := p.client.ExecuteTransactionBlock(ctx, txnBytes.TxBytes.Data(), []any{signature}, &types.SuiTransactionBlockResponseOptions{
		ShowInput:          true,
		ShowEffects:        true,
		ShowEvents:         true,
		ShowObjectChanges:  true,
		ShowBalanceChanges: true,
	}, types.TxnRequestTypeWaitForLocalExecution)
	if err != nil {
		log.Panic(err)
	}

	log.Println(txnResponse)
}
