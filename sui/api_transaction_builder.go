package sui

import (
	"context"

	"github.com/howjmay/sui-go/models"
	"github.com/howjmay/sui-go/sui_types"
)

type BatchTransactionRequest struct {
	Signer    *sui_types.SuiAddress
	TxnParams []models.RPCTransactionRequestParams
	Gas       *sui_types.ObjectID // optional
	GasBudget *models.BigInt
	// txnBuilderMode // optional // FIXME SuiTransactionBlockBuilderMode
}

// TODO: execution_mode : <SuiTransactionBlockBuilderMode>
func (s *ImplSuiAPI) BatchTransaction(
	ctx context.Context,
	req *BatchTransactionRequest,
) (*models.TransactionBytes, error) {
	resp := models.TransactionBytes{}
	return &resp, s.http.CallContext(ctx, &resp, batchTransaction, req.Signer, req.TxnParams, req.Gas, req.GasBudget)
}

type MergeCoinsRequest struct {
	Signer      *sui_types.SuiAddress `json:"signer,omitempty"`
	PrimaryCoin *sui_types.ObjectID   `json:"primaryCoin,omitempty"`
	CoinToMerge *sui_types.ObjectID   `json:"coinToMerge,omitempty"`
	Gas         *sui_types.ObjectID   `json:"gas,omitempty"` // optional
	GasBudget   *models.BigInt        `json:"gasBudget,omitempty"`
}

// MergeCoins Create an unsigned transaction to merge multiple coins into one coin.
func (s *ImplSuiAPI) MergeCoins(
	ctx context.Context,
	req *MergeCoinsRequest,
) (*models.TransactionBytes, error) {
	resp := models.TransactionBytes{}
	return &resp, s.http.CallContext(ctx, &resp, mergeCoins, req.Signer, req.PrimaryCoin, req.CoinToMerge, req.Gas, req.GasBudget)
}

type MoveCallRequest struct {
	Signer    *sui_types.SuiAddress
	PackageID *sui_types.PackageID
	Module    string
	Function  string
	TypeArgs  []string
	Arguments []any
	Gas       *sui_types.ObjectID // optional
	GasBudget *models.BigInt
	// txnBuilderMode // optional // FIXME SuiTransactionBlockBuilderMode
}

// MoveCall Create an unsigned transaction to execute a Move call on the network, by calling the specified function in the module of a given package.
// TODO: execution_mode : <SuiTransactionBlockBuilderMode>
// `arguments: []any` *SuiAddress can be arguments here, it will automatically convert to Address in hex string.
// [][]byte can't be passed. User should encode array of hex string.
func (s *ImplSuiAPI) MoveCall(
	ctx context.Context,
	req *MoveCallRequest,
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

type PayRequest struct {
	Signer     *sui_types.SuiAddress
	InputCoins []*sui_types.ObjectID
	Recipients []*sui_types.SuiAddress
	Amount     []*models.BigInt
	Gas        *sui_types.ObjectID // optional
	GasBudget  *models.BigInt
}

func (s *ImplSuiAPI) Pay(
	ctx context.Context,
	req *PayRequest,
) (*models.TransactionBytes, error) {
	resp := models.TransactionBytes{}
	return &resp, s.http.CallContext(ctx, &resp, pay, req.Signer, req.InputCoins, req.Recipients, req.Amount, req.Gas, req.GasBudget)
}

type PayAllSuiRequest struct {
	Signer     *sui_types.SuiAddress
	Recipient  *sui_types.SuiAddress
	InputCoins []*sui_types.ObjectID
	GasBudget  *models.BigInt
}

// PayAllSui Create an unsigned transaction to send all SUI coins to one recipient.
func (s *ImplSuiAPI) PayAllSui(
	ctx context.Context,
	req *PayAllSuiRequest,
) (*models.TransactionBytes, error) {
	resp := models.TransactionBytes{}
	return &resp, s.http.CallContext(ctx, &resp, payAllSui, req.Signer, req.InputCoins, req.Recipient, req.GasBudget)
}

type PaySuiRequest struct {
	Signer     *sui_types.SuiAddress
	InputCoins []*sui_types.ObjectID
	Recipients []*sui_types.SuiAddress
	Amount     []*models.BigInt
	GasBudget  *models.BigInt
}

// see explanation in https://forums.sui.io/t/how-to-use-the-sui-paysui-method/2282
func (s *ImplSuiAPI) PaySui(
	ctx context.Context,
	req *PaySuiRequest,
) (*models.TransactionBytes, error) {
	resp := models.TransactionBytes{}
	return &resp, s.http.CallContext(ctx, &resp, paySui, req.Signer, req.InputCoins, req.Recipients, req.Amount, req.GasBudget)
}

type PublishRequest struct {
	Sender          *sui_types.SuiAddress
	CompiledModules []*sui_types.Base64Data
	Dependencies    []*sui_types.ObjectID
	Gas             *sui_types.ObjectID // optional
	GasBudget       *models.BigInt
}

func (s *ImplSuiAPI) Publish(
	ctx context.Context,
	req *PublishRequest,
) (*models.TransactionBytes, error) {
	var resp models.TransactionBytes
	return &resp, s.http.CallContext(ctx, &resp, publish, req.Sender, req.CompiledModules, req.Dependencies, req.Gas, req.GasBudget)
}

type RequestAddStakeRequest struct {
	Signer    *sui_types.SuiAddress
	Coins     []*sui_types.ObjectID
	Amount    *models.BigInt // optional
	Validator *sui_types.SuiAddress
	Gas       *sui_types.ObjectID // optional
	GasBudget *models.BigInt
}

func (s *ImplSuiAPI) RequestAddStake(
	ctx context.Context,
	req *RequestAddStakeRequest,
) (*models.TransactionBytes, error) {
	var resp models.TransactionBytes
	return &resp, s.http.CallContext(ctx, &resp, requestAddStake, req.Signer, req.Coins, req.Amount, req.Validator, req.Gas, req.GasBudget)
}

type RequestWithdrawStakeRequest struct {
	Signer      *sui_types.SuiAddress
	StakedSuiId *sui_types.ObjectID
	Gas         *sui_types.ObjectID // optional
	GasBudget   *models.BigInt
}

func (s *ImplSuiAPI) RequestWithdrawStake(
	ctx context.Context,
	req *RequestWithdrawStakeRequest,
) (*models.TransactionBytes, error) {
	var resp models.TransactionBytes
	return &resp, s.http.CallContext(ctx, &resp, requestWithdrawStake, req.Signer, req.StakedSuiId, req.Gas, req.GasBudget)
}

type SplitCoinRequest struct {
	Signer       *sui_types.SuiAddress
	Coin         *sui_types.ObjectID
	SplitAmounts []*models.BigInt
	Gas          *sui_types.ObjectID // optional
	GasBudget    *models.BigInt
}

// SplitCoin Creates an unsigned transaction to split a coin object into multiple coins.
// better to replace with unsafe_pay API which consumes less gas
func (s *ImplSuiAPI) SplitCoin(
	ctx context.Context,
	req *SplitCoinRequest,
) (*models.TransactionBytes, error) {
	resp := models.TransactionBytes{}
	return &resp, s.http.CallContext(ctx, &resp, splitCoin, req.Signer, req.Coin, req.SplitAmounts, req.Gas, req.GasBudget)
}

type SplitCoinEqualRequest struct {
	Signer     *sui_types.SuiAddress
	Coin       *sui_types.ObjectID
	SplitCount *models.BigInt
	Gas        *sui_types.ObjectID // optional
	GasBudget  *models.BigInt
}

// SplitCoinEqual Creates an unsigned transaction to split a coin object into multiple equal-size coins.
// better to replace with unsafe_pay API which consumes less gas
func (s *ImplSuiAPI) SplitCoinEqual(
	ctx context.Context,
	req *SplitCoinEqualRequest,
) (*models.TransactionBytes, error) {
	resp := models.TransactionBytes{}
	return &resp, s.http.CallContext(ctx, &resp, splitCoinEqual, req.Signer, req.Coin, req.SplitCount, req.Gas, req.GasBudget)
}

type TransferObjectRequest struct {
	Signer    *sui_types.SuiAddress
	ObjectID  *sui_types.ObjectID
	Gas       *sui_types.ObjectID // optional
	GasBudget *models.BigInt
	Recipient *sui_types.SuiAddress
}

// TransferObject Create an unsigned transaction to transfer an object from one address to another. The object's type must allow public transfers
func (s *ImplSuiAPI) TransferObject(
	ctx context.Context,
	req *TransferObjectRequest,
) (*models.TransactionBytes, error) {
	resp := models.TransactionBytes{}
	return &resp, s.http.CallContext(ctx, &resp, transferObject, req.Signer, req.ObjectID, req.Gas, req.GasBudget, req.Recipient)
}

type TransferSuiRequest struct {
	Signer    *sui_types.SuiAddress
	ObjectID  *sui_types.ObjectID
	GasBudget *models.BigInt
	Recipient *sui_types.SuiAddress
	Amount    *models.BigInt // optional
}

// TransferSui Create an unsigned transaction to send SUI coin object to a Sui address. The SUI object is also used as the gas object.
func (s *ImplSuiAPI) TransferSui(
	ctx context.Context,
	req *TransferSuiRequest,
) (*models.TransactionBytes, error) {
	resp := models.TransactionBytes{}
	return &resp, s.http.CallContext(ctx, &resp, transferSui, req.Signer, req.ObjectID, req.GasBudget, req.Recipient, req.Amount)
}
