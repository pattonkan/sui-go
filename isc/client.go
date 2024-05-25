package isc

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

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
		new(models.BigInt).SetUint64(gasBudget),
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
		new(models.BigInt).SetUint64(gasBudget),
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

// object 'Assets' is owned by the Anchor object, and an 'Assets' object doesn't have ID, because it is a dynamic-field of Anchor object.
func (c *Client) GetAssets(
	ctx context.Context,
	anchorPackageID *sui_types.PackageID,
	anchorAddress *sui_types.ObjectID,
) (*Assets, error) {
	resGetObject, err := c.GetObject(
		context.Background(),
		anchorAddress,
		&models.SuiObjectDataOptions{
			ShowContent: true,
		})
	if err != nil {
		return nil, fmt.Errorf("failed to call GetObject(): %w", err)
	}

	b, err := json.Marshal(resGetObject.Data.Content.Data.MoveObject.Fields.(map[string]interface{})["assets"])
	if err != nil {
		return nil, fmt.Errorf("failed to access 'assets' fields: %w", err)
	}
	var normalizedAssets NormalizedAssets
	err = json.Unmarshal(b, &normalizedAssets)
	if err != nil {
		return nil, fmt.Errorf("failed to cast to 'NormalizedAssets' type: %w", err)
	}

	CoinsID := normalizedAssets.Fields.Coins.Fields.ID.ID
	resDynamicFields, err := c.GetDynamicFields(
		context.Background(),
		sui_types.MustObjectIDFromHex(CoinsID),
		nil,
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to call GetDynamicFields(): %w", err)
	}

	var assets Assets
	assets.Coins = make([]*models.Coin, len(resDynamicFields.Data))
	for i, coin := range resDynamicFields.Data {
		assets.Coins[i] = &models.Coin{
			CoinType:     coin.Name.Value.(string),
			CoinObjectID: &coin.ObjectID,
			Digest:       &coin.Digest,
		}
	}

	for _, coin := range assets.Coins {
		res, err := c.GetObject(
			context.Background(),
			coin.CoinObjectID,
			&models.SuiObjectDataOptions{
				ShowContent: true,
			})
		if err != nil {
			return nil, fmt.Errorf("failed to call GetObject(): %w", err)
		}
		fieldsMap := res.Data.Content.Data.MoveObject.Fields.((map[string]interface{}))
		bal, _ := strconv.ParseUint(fieldsMap["value"].(string), 10, 64)
		coin.Balance = new(models.BigInt).SetUint64(bal)
	}
	return &assets, nil
}

func (c *Client) CreateRequest(
	ctx context.Context,
	signer *sui_signer.Signer,
	packageID *sui_types.PackageID,
	anchorAddress *sui_types.ObjectID,
	iscContractName string,
	iscFunctionName string,
	args [][]byte,
	gasBudget uint64,
	execOptions *models.SuiTransactionBlockResponseOptions,
) (*models.SuiTransactionBlockResponse, error) {
	ptb := sui_types.NewProgrammableTransactionBuilder()

	// the return object is an Anchor object
	arg1 := ptb.Command(
		sui_types.Command{
			MoveCall: &sui_types.ProgrammableMoveCall{
				Package:       packageID,
				Module:        "request",
				Function:      "create_request",
				TypeArguments: []sui_types.TypeTag{},
				Arguments:     []sui_types.Argument{ptb.MustPure(iscContractName), ptb.MustPure(iscFunctionName), ptb.MustPure(args)},
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

func (c *Client) SendRequest(
	ctx context.Context,
	signer *sui_signer.Signer,
	packageID *sui_types.PackageID,
	anchorAddress *sui_types.ObjectID,
	reqObjID *sui_types.ObjectID,
	gasBudget uint64,
	execOptions *models.SuiTransactionBlockResponseOptions,
) (*models.SuiTransactionBlockResponse, error) {
	txnBytes, err := c.MoveCall(
		ctx,
		signer.Address,
		packageID,
		"anchor",
		"send_request",
		[]string{},
		[]any{anchorAddress.String(), reqObjID.String()},
		nil,
		new(models.BigInt).SetUint64(gasBudget),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to call send_request() move call: %w", err)
	}

	txnResponse, err := c.SignAndExecuteTransaction(ctx, signer, txnBytes.TxBytes, execOptions)
	if err != nil {
		return nil, fmt.Errorf("can't execute the transaction: %w", err)
	}

	return txnResponse, nil
}

func (c *Client) ReceiveRequest(
	ctx context.Context,
	signer *sui_signer.Signer,
	packageID *sui_types.PackageID,
	anchorAddress *sui_types.ObjectID,
	reqObjID *sui_types.ObjectID,
	gasBudget uint64,
	execOptions *models.SuiTransactionBlockResponseOptions,
) (*models.SuiTransactionBlockResponse, error) {
	txnBytes, err := c.MoveCall(
		ctx,
		signer.Address,
		packageID,
		"anchor",
		"receive_request",
		[]string{},
		[]any{anchorAddress.String(), reqObjID.String()},
		nil,
		new(models.BigInt).SetUint64(gasBudget),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to call receive_request() move call: %w", err)
	}

	txnResponse, err := c.SignAndExecuteTransaction(ctx, signer, txnBytes.TxBytes, execOptions)
	if err != nil {
		return nil, fmt.Errorf("can't execute the transaction: %w", err)
	}

	return txnResponse, nil
}
