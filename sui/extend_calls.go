package sui

import (
	"context"
	"fmt"
	"math/big"
	"strings"

	"github.com/fardream/go-bcs/bcs"
	"github.com/howjmay/sui-go/models"
	"github.com/howjmay/sui-go/sui_signer"
	"github.com/howjmay/sui-go/sui_types"
	"github.com/howjmay/sui-go/sui_types/serialization"
	"github.com/howjmay/sui-go/utils"
)

func (s *ImplSuiAPI) GetCoinObjsForTargetAmount(
	ctx context.Context,
	address *sui_types.SuiAddress,
	targetAmount uint64,
) (models.Coins, error) {
	coinType := models.SuiCoinType
	coins, err := s.GetCoins(ctx, address, &coinType, nil, 200)
	if err != nil {
		return nil, fmt.Errorf("failed to call GetCoins(): %w", err)
	}
	pickedCoins, err := models.PickupCoins(coins, new(big.Int).SetUint64(targetAmount), 0, 0, 0)
	if err != nil {
		return nil, err
	}
	return pickedCoins.Coins, nil
}

func (s *ImplSuiAPI) SignAndExecuteTransaction(
	ctx context.Context,
	signer *sui_signer.Signer,
	txBytes sui_types.Base64Data,
	options *models.SuiTransactionBlockResponseOptions,
) (*models.SuiTransactionBlockResponse, error) {
	// FIXME we need to support other intent
	signature, err := signer.SignTransactionBlock(txBytes, sui_signer.DefaultIntent())
	if err != nil {
		return nil, fmt.Errorf("failed to sign transaction block: %w", err)
	}
	resp, err := s.ExecuteTransactionBlock(
		ctx,
		txBytes,
		[]*sui_signer.Signature{&signature},
		options,
		models.TxnRequestTypeWaitForLocalExecution,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to execute transaction: %w", err)
	}
	if options.ShowEffects && !resp.Effects.Data.IsSuccess() {
		return resp, fmt.Errorf("failed to execute transaction: %v", resp.Effects.Data.V1.Status)
	}
	return resp, nil
}

func (s *ImplSuiAPI) BuildAndPublishContract(
	ctx context.Context,
	signer *sui_signer.Signer,
	contractPath string,
	gasBudget uint64,
	options *models.SuiTransactionBlockResponseOptions,
) (*models.SuiTransactionBlockResponse, *sui_types.PackageID, error) {
	modules, err := utils.MoveBuild(contractPath)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to build move contract: %w", err)
	}

	txnBytes, err := s.Publish(
		context.Background(),
		signer.Address,
		modules.Modules,
		modules.Dependencies,
		nil,
		models.NewSafeSuiBigInt(gasBudget),
	)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to publish move contract: %w", err)
	}
	txnResponse, err := s.SignAndExecuteTransaction(context.Background(), signer, txnBytes.TxBytes, options)
	if err != nil || !txnResponse.Effects.Data.IsSuccess() {
		return nil, nil, fmt.Errorf("failed to sign move contract tx: %w", err)
	}

	packageID, err := txnResponse.GetPublishedPackageID()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get move contract package ID: %w", err)
	}
	return txnResponse, packageID, nil
}

func (s *ImplSuiAPI) PublishContract(
	ctx context.Context,
	signer *sui_signer.Signer,
	modules []*sui_types.Base64Data,
	dependencies []*sui_types.SuiAddress,
	gasBudget uint64,
	options *models.SuiTransactionBlockResponseOptions,
) (*models.SuiTransactionBlockResponse, *sui_types.PackageID, error) {
	txnBytes, err := s.Publish(
		context.Background(),
		signer.Address,
		modules,
		dependencies,
		nil,
		models.NewSafeSuiBigInt(gasBudget),
	)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to publish move contract: %w", err)
	}
	txnResponse, err := s.SignAndExecuteTransaction(context.Background(), signer, txnBytes.TxBytes, options)
	if err != nil || !txnResponse.Effects.Data.IsSuccess() {
		return nil, nil, fmt.Errorf("failed to sign move contract tx: %w", err)
	}

	packageID, err := txnResponse.GetPublishedPackageID()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get move contract package ID: %w", err)
	}
	return txnResponse, packageID, nil
}

func (s *ImplSuiAPI) MintToken(
	ctx context.Context,
	signer *sui_signer.Signer,
	packageID *sui_types.PackageID,
	tokenName string,
	treasuryCap *sui_types.ObjectID,
	mintAmount uint64,
	options *models.SuiTransactionBlockResponseOptions,
) (*models.SuiTransactionBlockResponse, error) {
	txnBytes, err := s.MoveCall(
		ctx,
		signer.Address,
		packageID,
		tokenName,
		"mint",
		[]string{},
		[]any{treasuryCap.String(), fmt.Sprintf("%d", mintAmount), signer.Address.String()},
		nil,
		models.NewSafeSuiBigInt(DefaultGasBudget),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to call mint() move call: %w", err)
	}

	txnResponse, err := s.SignAndExecuteTransaction(ctx, signer, txnBytes.TxBytes, options)
	if err != nil {
		return nil, fmt.Errorf("can't execute the transaction: %w", err)
	}

	return txnResponse, nil
}

// NOTE: This a copy the query limit from our Rust JSON RPC backend, this needs to be kept in sync!
const QUERY_MAX_RESULT_LIMIT = 50

// GetSuiCoinsOwnedByAddress This function will retrieve a maximum of 200 coins.
func (s *ImplSuiAPI) GetSuiCoinsOwnedByAddress(ctx context.Context, address *sui_types.SuiAddress) (
	models.Coins,
	error,
) {
	coinType := models.SuiCoinType
	page, err := s.GetCoins(ctx, address, &coinType, nil, 200)
	if err != nil {
		return nil, err
	}
	return page.Data, nil
}

// BatchGetObjectsOwnedByAddress @param filterType You can specify filtering out the specified resources, this will fetch all resources if it is not empty ""
func (s *ImplSuiAPI) BatchGetObjectsOwnedByAddress(
	ctx context.Context,
	address *sui_types.SuiAddress,
	options *models.SuiObjectDataOptions,
	filterType string,
) ([]models.SuiObjectResponse, error) {
	filterType = strings.TrimSpace(filterType)
	return s.BatchGetFilteredObjectsOwnedByAddress(
		ctx, address, options, func(sod *models.SuiObjectData) bool {
			return filterType == "" || filterType == *sod.Type
		},
	)
}

func (s *ImplSuiAPI) BatchGetFilteredObjectsOwnedByAddress(
	ctx context.Context,
	address *sui_types.SuiAddress,
	options *models.SuiObjectDataOptions,
	filter func(*models.SuiObjectData) bool,
) ([]models.SuiObjectResponse, error) {
	query := models.SuiObjectResponseQuery{
		Options: &models.SuiObjectDataOptions{
			ShowType: true,
		},
	}
	filteringObjs, err := s.GetOwnedObjects(ctx, address, &query, nil, nil)
	if err != nil {
		return nil, err
	}
	objIds := make([]*sui_types.ObjectID, 0)
	for _, obj := range filteringObjs.Data {
		if obj.Data == nil {
			continue // error obj
		}
		if filter != nil && !filter(obj.Data) {
			continue // ignore objects if non-specified type
		}
		objIds = append(objIds, obj.Data.ObjectID)
	}

	return s.MultiGetObjects(ctx, objIds, options)
}

////// PTB impl

func BCS_RequestAddStake(
	signer *sui_types.SuiAddress,
	coins []*sui_types.ObjectRef,
	amount models.SafeSuiBigInt[uint64],
	validator *sui_types.SuiAddress,
	gasBudget, gasPrice uint64,
) ([]byte, error) {
	// build with BCS
	ptb := sui_types.NewProgrammableTransactionBuilder()
	amtArg, err := ptb.Pure(amount.Uint64())
	if err != nil {
		return nil, err
	}
	arg0, err := ptb.Obj(sui_types.SuiSystemMutObj)
	if err != nil {
		return nil, err
	}
	arg1 := ptb.Command(
		sui_types.Command{
			SplitCoins: &sui_types.ProgrammableSplitCoins{
				Coin:    sui_types.Argument{GasCoin: &serialization.EmptyEnum{}},
				Amounts: []sui_types.Argument{amtArg},
			},
		},
	) // the coin is split result argument
	arg2, err := ptb.Pure(validator)
	if err != nil {
		return nil, err
	}

	ptb.Command(
		sui_types.Command{
			MoveCall: &sui_types.ProgrammableMoveCall{
				Package:  sui_types.SuiPackageIdSuiSystem,
				Module:   sui_types.SuiSystemModuleName,
				Function: sui_types.AddStakeFunName,
				Arguments: []sui_types.Argument{
					arg0, arg1, arg2,
				},
			},
		},
	)
	pt := ptb.Finish()
	tx := sui_types.NewProgrammable(
		signer, pt, coins, gasBudget, gasPrice,
	)
	return bcs.Marshal(tx)
}

func BCS_RequestWithdrawStake(
	signer *sui_types.SuiAddress,
	stakedSuiRef sui_types.ObjectRef,
	gas []*sui_types.ObjectRef,
	gasBudget, gasPrice uint64,
) ([]byte, error) {
	// build with BCS
	ptb := sui_types.NewProgrammableTransactionBuilder()
	arg0, err := ptb.Obj(sui_types.SuiSystemMutObj)
	if err != nil {
		return nil, err
	}
	arg1, err := ptb.Obj(
		sui_types.ObjectArg{
			ImmOrOwnedObject: &stakedSuiRef,
		},
	)
	if err != nil {
		return nil, err
	}
	ptb.Command(
		sui_types.Command{
			MoveCall: &sui_types.ProgrammableMoveCall{
				Package:  sui_types.SuiPackageIdSuiSystem,
				Module:   sui_types.SuiSystemModuleName,
				Function: sui_types.WithdrawStakeFunName,
				Arguments: []sui_types.Argument{
					arg0, arg1,
				},
			},
		},
	)
	pt := ptb.Finish()
	tx := sui_types.NewProgrammable(
		signer, pt, gas, gasBudget, gasPrice,
	)
	return bcs.Marshal(tx)
}
