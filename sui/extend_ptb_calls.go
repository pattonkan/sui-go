package sui

import (
	"context"
	"strings"

	"github.com/fardream/go-bcs/bcs"
	"github.com/howjmay/go-sui-sdk/lib"
	"github.com/howjmay/go-sui-sdk/sui_types"
	"github.com/howjmay/go-sui-sdk/sui_types/sui_system_state"
	"github.com/howjmay/go-sui-sdk/types"
)

// NOTE: This a copy the query limit from our Rust JSON RPC backend, this needs to be kept in sync!
const QUERY_MAX_RESULT_LIMIT = 50

type suiBase64Data = lib.Base64Data

// GetSuiCoinsOwnedByAddress This function will retrieve a maximum of 200 coins.
func (s *ImplSuiAPI) GetSuiCoinsOwnedByAddress(ctx context.Context, address *sui_types.SuiAddress) (types.Coins, error) {
	coinType := types.SuiCoinType
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
	options *types.SuiObjectDataOptions,
	filterType string,
) ([]types.SuiObjectResponse, error) {
	filterType = strings.TrimSpace(filterType)
	return s.BatchGetFilteredObjectsOwnedByAddress(
		ctx, address, options, func(sod *types.SuiObjectData) bool {
			return filterType == "" || filterType == *sod.Type
		},
	)
}

func (s *ImplSuiAPI) BatchGetFilteredObjectsOwnedByAddress(
	ctx context.Context,
	address *sui_types.SuiAddress,
	options *types.SuiObjectDataOptions,
	filter func(*types.SuiObjectData) bool,
) ([]types.SuiObjectResponse, error) {
	query := types.SuiObjectResponseQuery{
		Options: &types.SuiObjectDataOptions{
			ShowType: true,
		},
	}
	filteringObjs, err := s.GetOwnedObjects(ctx, address, &query, nil, nil)
	if err != nil {
		return nil, err
	}
	objIds := make([]sui_types.ObjectID, 0)
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

// PTB impl

func BCS_RequestAddStake(
	signer *sui_types.SuiAddress,
	coins []*sui_types.ObjectRef,
	amount types.SafeSuiBigInt[uint64],
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
			SplitCoins: &struct {
				Argument  sui_types.Argument
				Arguments []sui_types.Argument
			}{
				Argument:  sui_types.Argument{GasCoin: &lib.EmptyEnum{}},
				Arguments: []sui_types.Argument{amtArg},
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
				Package:  sui_types.SuiSystemAddress,
				Module:   sui_system_state.SuiSystemModuleName,
				Function: sui_types.AddStakeFunName,
				Arguments: []sui_types.Argument{
					arg0, arg1, arg2,
				},
			},
		},
	)
	pt := ptb.Finish()
	tx := sui_types.NewProgrammable(
		*signer, coins, pt, gasBudget, gasPrice,
	)
	return bcs.Marshal(tx)
}

func BCS_RequestWithdrawStake(signer *sui_types.SuiAddress, stakedSuiRef sui_types.ObjectRef, gas []*sui_types.ObjectRef, gasBudget, gasPrice uint64) ([]byte, error) {
	// build with BCS
	ptb := sui_types.NewProgrammableTransactionBuilder()
	arg0, err := ptb.Obj(sui_types.SuiSystemMutObj)
	if err != nil {
		return nil, err
	}
	arg1, err := ptb.Obj(sui_types.ObjectArg{
		ImmOrOwnedObject: &stakedSuiRef,
	})
	if err != nil {
		return nil, err
	}
	ptb.Command(sui_types.Command{
		MoveCall: &sui_types.ProgrammableMoveCall{
			Package:  sui_types.SuiSystemAddress,
			Module:   sui_system_state.SuiSystemModuleName,
			Function: sui_types.WithdrawStakeFunName,
			Arguments: []sui_types.Argument{
				arg0, arg1,
			},
		},
	})
	pt := ptb.Finish()
	tx := sui_types.NewProgrammable(
		*signer, gas, pt, gasBudget, gasPrice,
	)
	return bcs.Marshal(tx)
}