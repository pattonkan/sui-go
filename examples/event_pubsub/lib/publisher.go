package lib

import (
	"context"
	"log"

	"github.com/howjmay/sui-go/models"
	"github.com/howjmay/sui-go/sui"
	"github.com/howjmay/sui-go/sui_types"
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
		models.NewSafeSuiBigInt(uint64(100000)),
	)
	if err != nil {
		log.Panic(err)
	}

	signature, err := p.account.SignTransactionBlock(txnBytes.TxBytes.Data(), sui_types.DefaultIntent())
	if err != nil {
		log.Panic(err)
	}

	txnResponse, err := p.client.ExecuteTransactionBlock(ctx, txnBytes.TxBytes.Data(), []any{signature}, &models.SuiTransactionBlockResponseOptions{
		ShowInput:          true,
		ShowEffects:        true,
		ShowEvents:         true,
		ShowObjectChanges:  true,
		ShowBalanceChanges: true,
	}, models.TxnRequestTypeWaitForLocalExecution)
	if err != nil {
		log.Panic(err)
	}

	log.Println(txnResponse)
}
