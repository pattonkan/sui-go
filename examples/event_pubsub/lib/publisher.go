package serialization

import (
	"context"
	"log"

	"github.com/howjmay/sui-go/models"
	"github.com/howjmay/sui-go/sui"
	"github.com/howjmay/sui-go/sui_signer"
	"github.com/howjmay/sui-go/sui_types"
)

type Publisher struct {
	client *sui.ImplSuiAPI
	signer *sui_signer.Signer
}

func NewPublisher(client *sui.ImplSuiAPI, signer *sui_signer.Signer) *Publisher {
	return &Publisher{
		client: client,
		signer: signer,
	}
}

func (p *Publisher) PublishEvents(ctx context.Context, packageID *sui_types.PackageID) {
	txnBytes, err := p.client.MoveCall(
		ctx,
		&models.MoveCallRequest{
			Signer:    p.signer.Address,
			PackageID: packageID,
			Module:    "eventpub",
			Function:  "emit_clock",
			TypeArgs:  []string{},
			Arguments: []any{},
			GasBudget: models.NewBigInt(100000),
		},
	)
	if err != nil {
		log.Panic(err)
	}

	signature, err := p.signer.SignTransactionBlock(txnBytes.TxBytes.Data(), sui_signer.DefaultIntent())
	if err != nil {
		log.Panic(err)
	}

	txnResponse, err := p.client.ExecuteTransactionBlock(ctx, &models.ExecuteTransactionBlockRequest{
		TxDataBytes: txnBytes.TxBytes.Data(),
		Signatures:  []*sui_signer.Signature{&signature},
		Options: &models.SuiTransactionBlockResponseOptions{
			ShowInput:          true,
			ShowEffects:        true,
			ShowEvents:         true,
			ShowObjectChanges:  true,
			ShowBalanceChanges: true,
		},
		RequestType: models.TxnRequestTypeWaitForLocalExecution,
	})
	if err != nil {
		log.Panic(err)
	}

	log.Println(txnResponse)
}
