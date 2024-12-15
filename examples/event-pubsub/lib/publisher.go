package lib

import (
	"context"
	"log"

	"github.com/howjmay/sui-go/sui"
	"github.com/howjmay/sui-go/suiclient"
	"github.com/howjmay/sui-go/suisigner"
)

type Publisher struct {
	client *suiclient.ClientImpl
	signer *suisigner.Signer
}

func NewPublisher(client *suiclient.ClientImpl, signer *suisigner.Signer) *Publisher {
	return &Publisher{
		client: client,
		signer: signer,
	}
}

func (p *Publisher) PublishEvents(ctx context.Context, packageId *sui.PackageId) {
	txnBytes, err := p.client.MoveCall(
		ctx,
		&suiclient.MoveCallRequest{
			Signer:    p.signer.Address,
			PackageId: packageId,
			Module:    "eventpub",
			Function:  "emit_clock",
			TypeArgs:  []string{},
			Arguments: []any{},
			GasBudget: sui.NewBigInt(100000),
		},
	)
	if err != nil {
		log.Panic(err)
	}

	signature, err := p.signer.SignTransactionBlock(txnBytes.TxBytes.Data(), suisigner.DefaultIntent())
	if err != nil {
		log.Panic(err)
	}

	txnResponse, err := p.client.ExecuteTransactionBlock(ctx, &suiclient.ExecuteTransactionBlockRequest{
		TxDataBytes: txnBytes.TxBytes.Data(),
		Signatures:  []*suisigner.Signature{&signature},
		Options: &suiclient.SuiTransactionBlockResponseOptions{
			ShowInput:          true,
			ShowEffects:        true,
			ShowEvents:         true,
			ShowObjectChanges:  true,
			ShowBalanceChanges: true,
		},
		RequestType: suiclient.TxnRequestTypeWaitForLocalExecution,
	})
	if err != nil {
		log.Panic(err)
	}

	log.Println(txnResponse)
}
