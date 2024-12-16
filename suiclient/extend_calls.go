package suiclient

import (
	"context"
	"fmt"
	"strings"

	"github.com/fardream/go-bcs/bcs"
	"github.com/pattonkan/sui-go/sui"
	"github.com/pattonkan/sui-go/sui/suiptb"
	"github.com/pattonkan/sui-go/suisigner"
	"github.com/pattonkan/sui-go/utils"
)

func (s *ClientImpl) SignAndExecuteTransaction(
	ctx context.Context,
	signer *suisigner.Signer,
	txBytes sui.Base64Data,
	options *SuiTransactionBlockResponseOptions,
) (*SuiTransactionBlockResponse, error) {
	// FIXME we need to support other intent
	signature, err := signer.SignTransactionBlock(txBytes, suisigner.DefaultIntent())
	if err != nil {
		return nil, fmt.Errorf("failed to sign transaction block: %w", err)
	}
	resp, err := s.ExecuteTransactionBlock(
		ctx,
		&ExecuteTransactionBlockRequest{
			TxDataBytes: txBytes,
			Signatures:  []*suisigner.Signature{&signature},
			Options:     options,
			RequestType: TxnRequestTypeWaitForLocalExecution,
		},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to execute transaction: %w", err)
	}
	if options.ShowEffects && !resp.Effects.Data.IsSuccess() {
		return resp, fmt.Errorf("failed to execute transaction: %v", resp.Effects.Data.V1.Status)
	}
	return resp, nil
}

func (s *ClientImpl) BuildAndPublishContract(
	ctx context.Context,
	signer *suisigner.Signer,
	contractPath string,
	gasBudget uint64,
	options *SuiTransactionBlockResponseOptions,
) (*SuiTransactionBlockResponse, *sui.PackageId, error) {
	modules, err := utils.MoveBuild(contractPath)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to build move contract: %w", err)
	}

	txnBytes, err := s.Publish(
		context.Background(),
		&PublishRequest{
			Sender:          signer.Address,
			CompiledModules: modules.Modules,
			Dependencies:    modules.Dependencies,
			GasBudget:       sui.NewBigInt(gasBudget),
		},
	)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to publish move contract: %w", err)
	}
	txnResponse, err := s.SignAndExecuteTransaction(context.Background(), signer, txnBytes.TxBytes, options)
	if err != nil || !txnResponse.Effects.Data.IsSuccess() {
		return nil, nil, fmt.Errorf("failed to sign move contract tx: %w", err)
	}

	packageId, err := txnResponse.GetPublishedPackageId()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get move contract package Id: %w", err)
	}
	return txnResponse, packageId, nil
}

func (s *ClientImpl) PublishContract(
	ctx context.Context,
	signer *suisigner.Signer,
	modules []*sui.Base64Data,
	dependencies []*sui.Address,
	gasBudget uint64,
	options *SuiTransactionBlockResponseOptions,
) (*SuiTransactionBlockResponse, *sui.PackageId, error) {
	txnBytes, err := s.Publish(
		context.Background(),
		&PublishRequest{
			Sender:          signer.Address,
			CompiledModules: modules,
			Dependencies:    dependencies,
			GasBudget:       sui.NewBigInt(gasBudget),
		},
	)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to publish move contract: %w", err)
	}
	txnResponse, err := s.SignAndExecuteTransaction(context.Background(), signer, txnBytes.TxBytes, options)
	if err != nil || !txnResponse.Effects.Data.IsSuccess() {
		return nil, nil, fmt.Errorf("failed to sign move contract tx: %w", err)
	}

	packageId, err := txnResponse.GetPublishedPackageId()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get move contract package Id: %w", err)
	}
	return txnResponse, packageId, nil
}

func (s *ClientImpl) MintToken(
	ctx context.Context,
	signer *suisigner.Signer,
	packageId *sui.PackageId,
	tokenName string,
	treasuryCap *sui.ObjectId,
	mintAmount uint64,
	options *SuiTransactionBlockResponseOptions,
) (*SuiTransactionBlockResponse, error) {
	txnBytes, err := s.MoveCall(
		ctx,
		&MoveCallRequest{
			Signer:    signer.Address,
			PackageId: packageId,
			Module:    tokenName,
			Function:  "mint",
			TypeArgs:  []string{},
			Arguments: []any{treasuryCap.String(), fmt.Sprintf("%d", mintAmount), signer.Address.String()},
			GasBudget: sui.NewBigInt(DefaultGasBudget),
		},
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
func (s *ClientImpl) GetSuiCoinsOwnedByAddress(ctx context.Context, address *sui.Address) (Coins, error) {
	page, err := s.GetCoins(ctx, &GetCoinsRequest{
		Owner: address,
		Limit: 200,
	})
	if err != nil {
		return nil, err
	}
	return page.Data, nil
}

// BatchGetObjectsOwnedByAddress @param filterType You can specify filtering out the specified resources, this will fetch all resources if it is not empty ""
func (s *ClientImpl) BatchGetObjectsOwnedByAddress(
	ctx context.Context,
	address *sui.Address,
	options *SuiObjectDataOptions,
	filterType string,
) ([]SuiObjectResponse, error) {
	filterType = strings.TrimSpace(filterType)
	return s.BatchGetFilteredObjectsOwnedByAddress(
		ctx, address, options, func(sod *SuiObjectData) bool {
			return filterType == "" || filterType == *sod.Type
		},
	)
}

func (s *ClientImpl) BatchGetFilteredObjectsOwnedByAddress(
	ctx context.Context,
	address *sui.Address,
	options *SuiObjectDataOptions,
	filter func(*SuiObjectData) bool,
) ([]SuiObjectResponse, error) {
	filteringObjs, err := s.GetOwnedObjects(ctx, &GetOwnedObjectsRequest{
		Address: address,
		Query: &SuiObjectResponseQuery{
			Options: &SuiObjectDataOptions{
				ShowType: true,
			},
		},
	})
	if err != nil {
		return nil, err
	}
	objIds := make([]*sui.ObjectId, 0)
	for _, obj := range filteringObjs.Data {
		if obj.Data == nil {
			continue // error obj
		}
		if filter != nil && !filter(obj.Data) {
			continue // ignore objects if non-specified type
		}
		objIds = append(objIds, obj.Data.ObjectId)
	}

	return s.MultiGetObjects(ctx, &MultiGetObjectsRequest{
		ObjectIds: objIds,
		Options:   options,
	})
}

func BCS_RequestAddStake(
	signer *sui.Address,
	coins []*sui.ObjectRef,
	amount *sui.BigInt,
	validator *sui.Address,
	gasBudget, gasPrice uint64,
) ([]byte, error) {
	// build with BCS
	ptb := suiptb.NewTransactionDataTransactionBuilder()
	amtArg, err := ptb.Pure(amount.Uint64())
	if err != nil {
		return nil, err
	}
	arg0, err := ptb.Obj(suiptb.SuiSystemMutObj)
	if err != nil {
		return nil, err
	}
	arg1 := ptb.Command(
		suiptb.Command{
			SplitCoins: &suiptb.ProgrammableSplitCoins{
				Coin:    suiptb.Argument{GasCoin: &sui.EmptyEnum{}},
				Amounts: []suiptb.Argument{amtArg},
			},
		},
	) // the coin is split result argument
	arg2, err := ptb.Pure(validator)
	if err != nil {
		return nil, err
	}

	ptb.Command(
		suiptb.Command{
			MoveCall: &suiptb.ProgrammableMoveCall{
				Package:  sui.SuiPackageIdSuiSystem,
				Module:   sui.SuiSystemModuleName,
				Function: "request_add_stake",
				Arguments: []suiptb.Argument{
					arg0, arg1, arg2,
				},
			},
		},
	)
	pt := ptb.Finish()
	tx := suiptb.NewTransactionData(
		signer, pt, coins, gasBudget, gasPrice,
	)
	return bcs.Marshal(tx)
}

func BCS_RequestWithdrawStake(
	signer *sui.Address,
	stakedSuiRef sui.ObjectRef,
	gas []*sui.ObjectRef,
	gasBudget, gasPrice uint64,
) ([]byte, error) {
	// build with BCS
	ptb := suiptb.NewTransactionDataTransactionBuilder()
	arg0, err := ptb.Obj(suiptb.SuiSystemMutObj)
	if err != nil {
		return nil, err
	}
	arg1, err := ptb.Obj(
		suiptb.ObjectArg{
			ImmOrOwnedObject: &stakedSuiRef,
		},
	)
	if err != nil {
		return nil, err
	}
	ptb.Command(
		suiptb.Command{
			MoveCall: &suiptb.ProgrammableMoveCall{
				Package:  sui.SuiPackageIdSuiSystem,
				Module:   sui.SuiSystemModuleName,
				Function: "request_withdraw_stake",
				Arguments: []suiptb.Argument{
					arg0, arg1,
				},
			},
		},
	)
	pt := ptb.Finish()
	tx := suiptb.NewTransactionData(
		signer, pt, gas, gasBudget, gasPrice,
	)
	return bcs.Marshal(tx)
}
