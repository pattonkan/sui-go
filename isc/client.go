package isc

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/fardream/go-bcs/bcs"
	"github.com/howjmay/sui-go/models"
	"github.com/howjmay/sui-go/sui"
	"github.com/howjmay/sui-go/sui_signer"
	"github.com/howjmay/sui-go/sui_types"
)

type Client struct {
	// API *sui.ImplSuiAPI
	*sui.ImplSuiAPI
}

func NewIscClient(api *sui.ImplSuiAPI) *Client {
	return &Client{
		api,
	}
}

// starts a new chain and transfer the Anchor to the signer
func (c *Client) StartNewChain(
	ctx context.Context,
	signer *sui_signer.Signer,
	packageID *sui_types.PackageID,
	gasBudget uint64,
	execOptions *models.SuiTransactionBlockResponseOptions,
) (*models.SuiTransactionBlockResponse, error) {
	ptb := sui_types.NewProgrammableTransactionBuilder()

	// the return object is an Anchor object
	arg1 := ptb.Command(
		sui_types.Command{
			MoveCall: &sui_types.ProgrammableMoveCall{
				Package:       packageID,
				Module:        "anchor",
				Function:      "start_new_chain",
				TypeArguments: []sui_types.TypeTag{},
				Arguments:     []sui_types.Argument{},
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

	// FIXME set the proper gas price
	coins, err := c.GetCoinObjectForGasFee(ctx, signer.Address, 10000, gasBudget)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch GasPayment object: %w", err)
	}

	tx := sui_types.NewProgrammable(
		signer.Address,
		pt,
		coins.CoinRefs(),
		gasBudget,
		1000, // TODO we may need to pass gas price
	)
	txnBytes, err := bcs.Marshal(tx)
	if err != nil {
		return nil, fmt.Errorf("can't marshal transaction into BCS encoding: %w", err)
	}

	txnResponse, err := c.SignAndExecuteTransaction(ctx, signer, txnBytes, execOptions)
	if err != nil {
		return nil, fmt.Errorf("can't execute the transaction: %w", err)
	}

	return txnResponse, nil
}

func (c *Client) SendCoin(
	ctx context.Context,
	signer *sui_signer.Signer,
	anchorPackageID *sui_types.PackageID,
	anchorAddress *sui_types.ObjectID,
	coinType string,
	coinObject *sui_types.ObjectID,
	gasBudget uint64,
	execOptions *models.SuiTransactionBlockResponseOptions,
) (*models.SuiTransactionBlockResponse, error) {
	txnBytes, err := c.MoveCall(
		ctx,
		signer.Address,
		anchorPackageID,
		"anchor",
		"send_coin",
		[]string{coinType},
		[]any{anchorAddress.String(), coinObject.String()},
		nil,
		models.NewSafeSuiBigInt(gasBudget),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to call send_coin() move call: %w", err)
	}

	txnResponse, err := c.SignAndExecuteTransaction(ctx, signer, txnBytes.TxBytes, execOptions)
	if err != nil {
		return nil, fmt.Errorf("can't execute the transaction: %w", err)
	}

	return txnResponse, nil
}

func (c *Client) ReceiveCoin(
	ctx context.Context,
	signer *sui_signer.Signer,
	anchorPackageID *sui_types.PackageID,
	anchorAddress *sui_types.ObjectID,
	coinType string,
	coinObject *sui_types.ObjectID,
	gasBudget uint64,
	execOptions *models.SuiTransactionBlockResponseOptions,
) (*models.SuiTransactionBlockResponse, error) {
	txnBytes, err := c.MoveCall(
		ctx,
		signer.Address,
		anchorPackageID,
		"anchor",
		"receive_coin",
		[]string{coinType},
		[]any{anchorAddress.String(), coinObject.String()},
		nil,
		models.NewSafeSuiBigInt(gasBudget),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to call receive_coin() move call: %w", err)
	}

	txnResponse, err := c.SignAndExecuteTransaction(ctx, signer, txnBytes.TxBytes, execOptions)
	if err != nil {
		return nil, fmt.Errorf("can't execute the transaction: %w", err)
	}

	return txnResponse, nil
}

func (c *Client) GetAssetsFromAnchor(
	ctx context.Context,
	anchorPackageID *sui_types.PackageID,
	anchorAddress *sui_types.ObjectID,
) (*NormalizedAssets, error) {
	res, err := c.GetObject(
		context.Background(),
		anchorAddress,
		&models.SuiObjectDataOptions{
			ShowContent: true,
		})
	if err != nil {
		return nil, fmt.Errorf("failed to call GetObject(): %w", err)
	}

	b, err := json.Marshal(res.Data.Content.Data.MoveObject.Fields.(map[string]interface{})["assets"])
	if err != nil {
		return nil, fmt.Errorf("failed to access 'assets' fields: %w", err)
	}
	var assets NormalizedAssets
	err = json.Unmarshal(b, &assets)
	if err != nil {
		return nil, fmt.Errorf("failed to cast to 'NormalizedAssets' type: %w", err)
	}
	return &assets, nil
}
