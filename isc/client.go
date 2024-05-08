package isc

import (
	"context"
	"fmt"
	"math/big"

	"github.com/fardream/go-bcs/bcs"
	"github.com/howjmay/sui-go/models"
	"github.com/howjmay/sui-go/sui"
	"github.com/howjmay/sui-go/sui_signer"
	"github.com/howjmay/sui-go/sui_types"
)

type Client struct {
	API *sui.ImplSuiAPI
}

func NewIscClient(api *sui.ImplSuiAPI) *Client {
	return &Client{
		API: api,
	}
}

func (c *Client) StartNewChain(ctx context.Context, signer *sui_signer.Signer, packageID *sui_types.PackageID, anchorCap *sui_types.ObjectID) (*models.SuiTransactionBlockResponse, error) {
	ptb := sui_types.NewProgrammableTransactionBuilder()
	argAnchor, err := ptb.Pure(anchorCap.String())
	if err != nil {
		panic(err)
	}
	arg1 := ptb.Command(
		sui_types.Command{
			MoveCall: &sui_types.ProgrammableMoveCall{
				Package:       packageID,
				Module:        "anchor",
				Function:      "start_new_chain",
				TypeArguments: []sui_types.TypeTag{},
				Arguments:     []sui_types.Argument{argAnchor},
			},
		},
	)

	ptb.Command(
		sui_types.Command{
			TransferObjects: &sui_types.ProgrammableTransferObjects{
				Objects: []sui_types.Argument{arg1},
				Address: ptb.MustPure(signer.Address),
			},
		},
	)
	pt := ptb.Finish()

	gasBudget := uint64(1000000)
	coins, err := c.API.GetCoins(context.Background(), signer.Address, nil, nil, 10)
	if err != nil {
		panic(err)
	}
	pickedCoins, err := models.PickupCoins(coins, big.NewInt(100000), gasBudget, 10, 10)
	if err != nil {
		panic(err)
	}

	tx := sui_types.NewProgrammable(
		signer.Address,
		pt,
		[]*sui_types.ObjectRef{
			pickedCoins.CoinRefs()[0],
		},
		gasBudget,
		1000,
	)

	txnBytes, err := bcs.Marshal(tx)
	if err != nil {
		panic(err)
	}

	signature, err := signer.SignTransactionBlock(txnBytes, sui_signer.DefaultIntent())
	if err != nil {
		panic(err)
	}
	txnResponse, err := c.API.ExecuteTransactionBlock(context.TODO(), txnBytes, []any{signature}, &models.SuiTransactionBlockResponseOptions{
		ShowInput:          true,
		ShowEffects:        true,
		ShowEvents:         true,
		ShowObjectChanges:  true,
		ShowBalanceChanges: true,
	}, models.TxnRequestTypeWaitForLocalExecution)
	if err != nil {
		panic(err)
	}
	fmt.Println(txnResponse)

	return txnResponse, nil
}
