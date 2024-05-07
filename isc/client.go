package isc

import (
	"context"

	"github.com/howjmay/sui-go/models"
	"github.com/howjmay/sui-go/sui"
	"github.com/howjmay/sui-go/sui_signer"
	"github.com/howjmay/sui-go/sui_types"
)

type Client struct {
	API *sui.ImplSuiAPI
}

func NewClient(api *sui.ImplSuiAPI) *Client {
	return &Client{
		API: api,
	}
}

func (c *Client) StartNewChain(ctx context.Context, signer *sui_signer.Signer, packageID *sui_types.PackageID, anchorCap *sui_types.ObjectID) (*models.SuiTransactionBlockResponse, error) {
	txnBytes, err := c.API.MoveCall(ctx,
		signer.Address,
		packageID,
		"isc",
		"start_new_chain",
		[]string{},
		[]any{anchorCap.String()},
		nil,
		models.NewSafeSuiBigInt(uint64(10000000)),
	)
	if err != nil {
		panic(err)
	}

	signature, err := signer.SignTransactionBlock(txnBytes.TxBytes.Data(), sui_signer.DefaultIntent())
	if err != nil {
		panic(err)
	}
	txnResponse, err := c.API.ExecuteTransactionBlock(context.TODO(), txnBytes.TxBytes.Data(), []any{signature}, &models.SuiTransactionBlockResponseOptions{
		ShowInput:          true,
		ShowEffects:        true,
		ShowEvents:         true,
		ShowObjectChanges:  true,
		ShowBalanceChanges: true,
	}, models.TxnRequestTypeWaitForLocalExecution)
	if err != nil {
		panic(err)
	}

	return txnResponse, nil
}

func (c *Client) RegisterIscToken(ctx context.Context, signer *sui_signer.Signer, packageID *sui_types.PackageID) (*models.SuiTransactionBlockResponse, error) {
	txnBytes, err := c.API.MoveCall(ctx,
		signer.Address,
		packageID,
		"iscanchor",
		"register_isc_token",
		[]string{"0x28e450c30438509db86248253e07c8cdfd3e47e070ec69776f869fafbe6752c4"},
		[]any{},
		nil,
		models.NewSafeSuiBigInt(uint64(10000000)),
	)
	if err != nil {
		panic(err)
	}

	signature, err := signer.SignTransactionBlock(txnBytes.TxBytes.Data(), sui_signer.DefaultIntent())
	if err != nil {
		panic(err)
	}
	txnResponse, err := c.API.ExecuteTransactionBlock(context.TODO(), txnBytes.TxBytes.Data(), []any{signature}, &models.SuiTransactionBlockResponseOptions{
		ShowInput:          true,
		ShowEffects:        true,
		ShowEvents:         true,
		ShowObjectChanges:  true,
		ShowBalanceChanges: true,
	}, models.TxnRequestTypeWaitForLocalExecution)
	if err != nil {
		panic(err)
	}

	return txnResponse, nil
}

func (c *Client) MintToken(ctx context.Context, signer *sui_signer.Signer, packageID *sui_types.PackageID) (*models.SuiTransactionBlockResponse, error) {
	txnBytes, err := c.API.MoveCall(
		ctx,
		signer.Address,
		packageID,
		"testcoin",
		"mint",
		[]string{},
		[]any{},
		nil,
		models.NewSafeSuiBigInt(uint64(10000000)),
	)
	if err != nil {
		panic(err)
	}

	signature, err := signer.SignTransactionBlock(txnBytes.TxBytes.Data(), sui_signer.DefaultIntent())
	if err != nil {
		panic(err)
	}
	txnResponse, err := c.API.ExecuteTransactionBlock(context.TODO(), txnBytes.TxBytes.Data(), []any{signature}, &models.SuiTransactionBlockResponseOptions{
		ShowInput:          true,
		ShowEffects:        true,
		ShowEvents:         true,
		ShowObjectChanges:  true,
		ShowBalanceChanges: true,
	}, models.TxnRequestTypeWaitForLocalExecution)
	if err != nil {
		panic(err)
	}

	return txnResponse, nil
}
