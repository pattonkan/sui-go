package sui

import (
	"context"

	"github.com/howjmay/sui-go/models"
)

// TODO: execution_mode : <SuiTransactionBlockBuilderMode>
func (s *ImplSuiAPI) BatchTransaction(
	ctx context.Context,
	req *models.BatchTransactionRequest,
) (*models.TransactionBytes, error) {
	resp := models.TransactionBytes{}
	return &resp, s.http.CallContext(ctx, &resp, batchTransaction, req.Signer, req.TxnParams, req.Gas, req.GasBudget)
}

// MergeCoins Create an unsigned transaction to merge multiple coins into one coin.
func (s *ImplSuiAPI) MergeCoins(
	ctx context.Context,
	req *models.MergeCoinsRequest,
) (*models.TransactionBytes, error) {
	resp := models.TransactionBytes{}
	return &resp, s.http.CallContext(ctx, &resp, mergeCoins, req.Signer, req.PrimaryCoin, req.CoinToMerge, req.Gas, req.GasBudget)
}

// MoveCall Create an unsigned transaction to execute a Move call on the network, by calling the specified function in the module of a given package.
// TODO: execution_mode : <SuiTransactionBlockBuilderMode>
// `arguments: []any` *SuiAddress can be arguments here, it will automatically convert to Address in hex string.
// [][]byte can't be passed. User should encode array of hex string.
func (s *ImplSuiAPI) MoveCall(
	ctx context.Context,
	req *models.MoveCallRequest,
) (*models.TransactionBytes, error) {
	resp := models.TransactionBytes{}
	return &resp, s.http.CallContext(
		ctx,
		&resp,
		moveCall,
		req.Signer,
		req.PackageID,
		req.Module,
		req.Function,
		req.TypeArgs,
		req.Arguments,
		req.Gas,
		req.GasBudget,
	)
}

func (s *ImplSuiAPI) Pay(
	ctx context.Context,
	req *models.PayRequest,
) (*models.TransactionBytes, error) {
	resp := models.TransactionBytes{}
	return &resp, s.http.CallContext(ctx, &resp, pay, req.Signer, req.InputCoins, req.Recipients, req.Amount, req.Gas, req.GasBudget)
}

// PayAllSui Create an unsigned transaction to send all SUI coins to one recipient.
func (s *ImplSuiAPI) PayAllSui(
	ctx context.Context,
	req *models.PayAllSuiRequest,
) (*models.TransactionBytes, error) {
	resp := models.TransactionBytes{}
	return &resp, s.http.CallContext(ctx, &resp, payAllSui, req.Signer, req.InputCoins, req.Recipient, req.GasBudget)
}

// see explanation in https://forums.sui.io/t/how-to-use-the-sui-paysui-method/2282
func (s *ImplSuiAPI) PaySui(
	ctx context.Context,
	req *models.PaySuiRequest,
) (*models.TransactionBytes, error) {
	resp := models.TransactionBytes{}
	return &resp, s.http.CallContext(ctx, &resp, paySui, req.Signer, req.InputCoins, req.Recipients, req.Amount, req.GasBudget)
}

func (s *ImplSuiAPI) Publish(
	ctx context.Context,
	req *models.PublishRequest,
) (*models.TransactionBytes, error) {
	var resp models.TransactionBytes
	return &resp, s.http.CallContext(ctx, &resp, publish, req.Sender, req.CompiledModules, req.Dependencies, req.Gas, req.GasBudget)
}

func (s *ImplSuiAPI) RequestAddStake(
	ctx context.Context,
	req *models.RequestAddStakeRequest,
) (*models.TransactionBytes, error) {
	var resp models.TransactionBytes
	return &resp, s.http.CallContext(ctx, &resp, requestAddStake, req.Signer, req.Coins, req.Amount, req.Validator, req.Gas, req.GasBudget)
}

func (s *ImplSuiAPI) RequestWithdrawStake(
	ctx context.Context,
	req *models.RequestWithdrawStakeRequest,
) (*models.TransactionBytes, error) {
	var resp models.TransactionBytes
	return &resp, s.http.CallContext(ctx, &resp, requestWithdrawStake, req.Signer, req.StakedSuiId, req.Gas, req.GasBudget)
}

// SplitCoin Creates an unsigned transaction to split a coin object into multiple coins.
// better to replace with unsafe_pay API which consumes less gas
func (s *ImplSuiAPI) SplitCoin(
	ctx context.Context,
	req *models.SplitCoinRequest,
) (*models.TransactionBytes, error) {
	resp := models.TransactionBytes{}
	return &resp, s.http.CallContext(ctx, &resp, splitCoin, req.Signer, req.Coin, req.SplitAmounts, req.Gas, req.GasBudget)
}

// SplitCoinEqual Creates an unsigned transaction to split a coin object into multiple equal-size coins.
// better to replace with unsafe_pay API which consumes less gas
func (s *ImplSuiAPI) SplitCoinEqual(
	ctx context.Context,
	req *models.SplitCoinEqualRequest,
) (*models.TransactionBytes, error) {
	resp := models.TransactionBytes{}
	return &resp, s.http.CallContext(ctx, &resp, splitCoinEqual, req.Signer, req.Coin, req.SplitCount, req.Gas, req.GasBudget)
}

// TransferObject Create an unsigned transaction to transfer an object from one address to another. The object's type must allow public transfers
func (s *ImplSuiAPI) TransferObject(
	ctx context.Context,
	req *models.TransferObjectRequest,
) (*models.TransactionBytes, error) {
	resp := models.TransactionBytes{}
	return &resp, s.http.CallContext(ctx, &resp, transferObject, req.Signer, req.ObjectID, req.Gas, req.GasBudget, req.Recipient)
}

// TransferSui Create an unsigned transaction to send SUI coin object to a Sui address. The SUI object is also used as the gas object.
func (s *ImplSuiAPI) TransferSui(
	ctx context.Context,
	req *models.TransferSuiRequest,
) (*models.TransactionBytes, error) {
	resp := models.TransactionBytes{}
	return &resp, s.http.CallContext(ctx, &resp, transferSui, req.Signer, req.ObjectID, req.GasBudget, req.Recipient, req.Amount)
}
