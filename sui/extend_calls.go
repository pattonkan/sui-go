package sui

import (
	"context"
	"fmt"
	"math/big"

	"github.com/howjmay/sui-go/models"
	"github.com/howjmay/sui-go/sui_signer"
	"github.com/howjmay/sui-go/sui_types"
)

func (s *ImplSuiAPI) GetCoinObjectForGasFee(
	ctx context.Context,
	address *sui_types.SuiAddress,
	targetAmount uint64,
	gasBudget uint64,
) (models.Coins, error) {
	coinType := models.SuiCoinType
	coins, err := s.GetCoins(ctx, address, &coinType, nil, 200)
	if err != nil {
		return nil, fmt.Errorf("failed to call GetCoins(): %w", err)
	}
	pickedCoins, err := models.PickupCoins(coins, new(big.Int).SetUint64(targetAmount), gasBudget, 0, 0)
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
	resp, err := s.ExecuteTransactionBlock(ctx, txBytes, []*sui_signer.Signature{&signature}, options, models.TxnRequestTypeWaitForLocalExecution)
	if err != nil {
		return nil, fmt.Errorf("failed to execute transaction: %w", err)
	}
	if options.ShowEffects && resp.Effects.Data.V1.Status.Status != models.ExecutionStatusSuccess {
		return resp, fmt.Errorf("failed to execute transaction: %v", resp.Effects.Data.V1.Status)
	}
	return resp, nil
}
