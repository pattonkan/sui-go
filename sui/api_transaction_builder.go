package sui

import (
	"context"

	"github.com/howjmay/sui-go/models"
	"github.com/howjmay/sui-go/sui_types"
)

// TODO: execution_mode : <SuiTransactionBlockBuilderMode>
func (s *ImplSuiAPI) BatchTransaction(
	ctx context.Context,
	signer *sui_types.SuiAddress,
	txnParams []map[string]interface{},
	gas *sui_types.ObjectID,
	gasBudget uint64,
) (*models.TransactionBytes, error) {
	resp := models.TransactionBytes{}
	return &resp, s.http.CallContext(ctx, &resp, batchTransaction, signer, txnParams, gas, gasBudget)
}

// MergeCoins Create an unsigned transaction to merge multiple coins into one coin.
func (s *ImplSuiAPI) MergeCoins(
	ctx context.Context,
	signer *sui_types.SuiAddress,
	primaryCoin, coinToMerge *sui_types.ObjectID,
	gas *sui_types.ObjectID,
	gasBudget *models.BigInt,
) (*models.TransactionBytes, error) {
	resp := models.TransactionBytes{}
	return &resp, s.http.CallContext(ctx, &resp, mergeCoins, signer, primaryCoin, coinToMerge, gas, gasBudget)
}

// MoveCall Create an unsigned transaction to execute a Move call on the network, by calling the specified function in the module of a given package.
// TODO: execution_mode : <SuiTransactionBlockBuilderMode>
// `arguments: []any` *SuiAddress can be arguments here, it will automatically convert to Address in hex string.
// [][]byte can't be passed. User should encode array of hex string.
func (s *ImplSuiAPI) MoveCall(
	ctx context.Context,
	signer *sui_types.SuiAddress,
	packageId *sui_types.ObjectID,
	module string,
	function string,
	typeArgs []string,
	arguments []any,
	gas *sui_types.ObjectID,
	gasBudget *models.BigInt,
) (*models.TransactionBytes, error) {
	resp := models.TransactionBytes{}
	return &resp, s.http.CallContext(
		ctx,
		&resp,
		moveCall,
		signer,
		packageId,
		module,
		function,
		typeArgs,
		arguments,
		gas,
		gasBudget,
	)
}

func (s *ImplSuiAPI) Pay(
	ctx context.Context,
	signer *sui_types.SuiAddress,
	inputCoins []*sui_types.ObjectID,
	recipients []*sui_types.SuiAddress,
	amount []*models.BigInt,
	gas *sui_types.ObjectID,
	gasBudget *models.BigInt,
) (*models.TransactionBytes, error) {
	resp := models.TransactionBytes{}
	return &resp, s.http.CallContext(ctx, &resp, pay, signer, inputCoins, recipients, amount, gas, gasBudget)
}

// PayAllSui Create an unsigned transaction to send all SUI coins to one recipient.
func (s *ImplSuiAPI) PayAllSui(
	ctx context.Context,
	signer *sui_types.SuiAddress,
	recipient *sui_types.SuiAddress,
	inputCoins []*sui_types.ObjectID,
	gasBudget *models.BigInt,
) (*models.TransactionBytes, error) {
	resp := models.TransactionBytes{}
	return &resp, s.http.CallContext(ctx, &resp, payAllSui, signer, inputCoins, recipient, gasBudget)
}

// see explanation in https://forums.sui.io/t/how-to-use-the-sui-paysui-method/2282
func (s *ImplSuiAPI) PaySui(
	ctx context.Context,
	signer *sui_types.SuiAddress,
	inputCoins []*sui_types.ObjectID,
	recipients []*sui_types.SuiAddress,
	amount []*models.BigInt,
	gasBudget *models.BigInt,
) (*models.TransactionBytes, error) {
	resp := models.TransactionBytes{}
	return &resp, s.http.CallContext(ctx, &resp, paySui, signer, inputCoins, recipients, amount, gasBudget)
}

func (s *ImplSuiAPI) Publish(
	ctx context.Context,
	sender *sui_types.SuiAddress,
	compiledModules []*sui_types.Base64Data,
	dependencies []*sui_types.ObjectID,
	gas *sui_types.ObjectID, // optional
	gasBudget *models.BigInt,
) (*models.TransactionBytes, error) {
	var resp models.TransactionBytes
	return &resp, s.http.CallContext(ctx, &resp, publish, sender, compiledModules, dependencies, gas, gasBudget)
}

func (s *ImplSuiAPI) RequestAddStake(
	ctx context.Context,
	signer *sui_types.SuiAddress,
	coins []*sui_types.ObjectID,
	amount *models.BigInt,
	validator *sui_types.SuiAddress,
	gas *sui_types.ObjectID,
	gasBudget *models.BigInt,
) (*models.TransactionBytes, error) {
	var resp models.TransactionBytes
	return &resp, s.http.CallContext(ctx, &resp, requestAddStake, signer, coins, amount, validator, gas, gasBudget)
}

func (s *ImplSuiAPI) RequestWithdrawStake(
	ctx context.Context,
	signer *sui_types.SuiAddress,
	stakedSuiId *sui_types.ObjectID,
	gas *sui_types.ObjectID,
	gasBudget *models.BigInt,
) (*models.TransactionBytes, error) {
	var resp models.TransactionBytes
	return &resp, s.http.CallContext(ctx, &resp, requestWithdrawStake, signer, stakedSuiId, gas, gasBudget)
}

// SplitCoin Creates an unsigned transaction to split a coin object into multiple coins.
// better to replace with unsafe_pay API which consumes less gas
func (s *ImplSuiAPI) SplitCoin(
	ctx context.Context,
	signer *sui_types.SuiAddress,
	coin *sui_types.ObjectID,
	splitAmounts []*models.BigInt,
	gas *sui_types.ObjectID, // optional
	gasBudget *models.BigInt,
) (*models.TransactionBytes, error) {
	resp := models.TransactionBytes{}
	return &resp, s.http.CallContext(ctx, &resp, splitCoin, signer, coin, splitAmounts, gas, gasBudget)
}

// SplitCoinEqual Creates an unsigned transaction to split a coin object into multiple equal-size coins.
// better to replace with unsafe_pay API which consumes less gas
func (s *ImplSuiAPI) SplitCoinEqual(
	ctx context.Context,
	signer *sui_types.SuiAddress,
	coin *sui_types.ObjectID,
	splitCount *models.BigInt,
	gas *sui_types.ObjectID,
	gasBudget *models.BigInt,
) (*models.TransactionBytes, error) {
	resp := models.TransactionBytes{}
	return &resp, s.http.CallContext(ctx, &resp, splitCoinEqual, signer, coin, splitCount, gas, gasBudget)
}

// TransferObject Create an unsigned transaction to transfer an object from one address to another. The object's type must allow public transfers
func (s *ImplSuiAPI) TransferObject(
	ctx context.Context,
	signer *sui_types.SuiAddress,
	recipient *sui_types.SuiAddress,
	objID *sui_types.ObjectID,
	gas *sui_types.ObjectID,
	gasBudget *models.BigInt,
) (*models.TransactionBytes, error) {
	resp := models.TransactionBytes{}
	return &resp, s.http.CallContext(ctx, &resp, transferObject, signer, objID, gas, gasBudget, recipient)
}

// TransferSui Create an unsigned transaction to send SUI coin object to a Sui address. The SUI object is also used as the gas object.
func (s *ImplSuiAPI) TransferSui(
	ctx context.Context,
	signer *sui_types.SuiAddress,
	recipient *sui_types.SuiAddress,
	objID *sui_types.ObjectID,
	amount *models.BigInt,
	gasBudget *models.BigInt,
) (*models.TransactionBytes, error) {
	resp := models.TransactionBytes{}
	return &resp, s.http.CallContext(ctx, &resp, transferSui, signer, objID, gasBudget, recipient, amount)
}
