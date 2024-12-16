package suiclient

import (
	"context"

	"github.com/pattonkan/sui-go/sui"
)

// the returned encoded TransactionData by RPC JSON APIs
type TransactionBytes struct {
	// the gas object to be used
	Gas []*sui.ObjectRef `json:"gas"`
	// objects to be used in this transaction
	InputObjects []*InputObjectKind `json:"inputObjects"`
	// transaction data bytes
	TxBytes sui.Base64Data `json:"txBytes"`
}

type InputObjectKind map[string]interface{}

type BatchTransactionRequest struct {
	Signer    *sui.Address
	TxnParams []RPCTransactionRequestParams
	Gas       *sui.ObjectId // optional
	GasBudget *sui.BigInt
	// txnBuilderMode // optional // FIXME SuiTransactionBlockBuilderMode
}

// TODO: execution_mode : <SuiTransactionBlockBuilderMode>
func (s *ClientImpl) BatchTransaction(
	ctx context.Context,
	req *BatchTransactionRequest,
) (*TransactionBytes, error) {
	resp := TransactionBytes{}
	return &resp, s.http.CallContext(ctx, &resp, batchTransaction, req.Signer, req.TxnParams, req.Gas, req.GasBudget)
}

type MergeCoinsRequest struct {
	Signer      *sui.Address  `json:"signer,omitempty"`
	PrimaryCoin *sui.ObjectId `json:"primaryCoin,omitempty"`
	CoinToMerge *sui.ObjectId `json:"coinToMerge,omitempty"`
	Gas         *sui.ObjectId `json:"gas,omitempty"` // optional
	GasBudget   *sui.BigInt   `json:"gasBudget,omitempty"`
}

// MergeCoins Create an unsigned transaction to merge multiple coins into one coin.
func (s *ClientImpl) MergeCoins(
	ctx context.Context,
	req *MergeCoinsRequest,
) (*TransactionBytes, error) {
	resp := TransactionBytes{}
	return &resp, s.http.CallContext(ctx, &resp, mergeCoins, req.Signer, req.PrimaryCoin, req.CoinToMerge, req.Gas, req.GasBudget)
}

type MoveCallRequest struct {
	Signer    *sui.Address
	PackageId *sui.PackageId
	Module    string
	Function  string
	TypeArgs  []string
	Arguments []any
	Gas       *sui.ObjectId // optional
	GasBudget *sui.BigInt
	// txnBuilderMode // optional // FIXME SuiTransactionBlockBuilderMode
}

// MoveCall Create an unsigned transaction to execute a Move call on the network, by calling the specified function in the module of a given package.
// TODO: execution_mode : <SuiTransactionBlockBuilderMode>
// `arguments: []any` *Address can be arguments here, it will automatically convert to Address in hex string.
// [][]byte can't be passed. User should encode array of hex string.
func (s *ClientImpl) MoveCall(
	ctx context.Context,
	req *MoveCallRequest,
) (*TransactionBytes, error) {
	resp := TransactionBytes{}
	return &resp, s.http.CallContext(
		ctx,
		&resp,
		moveCall,
		req.Signer,
		req.PackageId,
		req.Module,
		req.Function,
		req.TypeArgs,
		req.Arguments,
		req.Gas,
		req.GasBudget,
	)
}

type PayRequest struct {
	Signer     *sui.Address
	InputCoins []*sui.ObjectId
	Recipients []*sui.Address
	Amount     []*sui.BigInt
	Gas        *sui.ObjectId // optional
	GasBudget  *sui.BigInt
}

func (s *ClientImpl) Pay(
	ctx context.Context,
	req *PayRequest,
) (*TransactionBytes, error) {
	resp := TransactionBytes{}
	return &resp, s.http.CallContext(ctx, &resp, pay, req.Signer, req.InputCoins, req.Recipients, req.Amount, req.Gas, req.GasBudget)
}

type PayAllSuiRequest struct {
	Signer     *sui.Address
	Recipient  *sui.Address
	InputCoins []*sui.ObjectId
	GasBudget  *sui.BigInt
}

// PayAllSui Create an unsigned transaction to send all SUI coins to one recipient.
func (s *ClientImpl) PayAllSui(
	ctx context.Context,
	req *PayAllSuiRequest,
) (*TransactionBytes, error) {
	resp := TransactionBytes{}
	return &resp, s.http.CallContext(ctx, &resp, payAllSui, req.Signer, req.InputCoins, req.Recipient, req.GasBudget)
}

type PaySuiRequest struct {
	Signer     *sui.Address
	InputCoins []*sui.ObjectId
	Recipients []*sui.Address
	Amount     []*sui.BigInt
	GasBudget  *sui.BigInt
}

// see explanation in https://forums.sui.io/t/how-to-use-the-sui-paysui-method/2282
func (s *ClientImpl) PaySui(
	ctx context.Context,
	req *PaySuiRequest,
) (*TransactionBytes, error) {
	resp := TransactionBytes{}
	return &resp, s.http.CallContext(ctx, &resp, paySui, req.Signer, req.InputCoins, req.Recipients, req.Amount, req.GasBudget)
}

type PublishRequest struct {
	Sender          *sui.Address
	CompiledModules []*sui.Base64Data
	Dependencies    []*sui.ObjectId
	Gas             *sui.ObjectId // optional
	GasBudget       *sui.BigInt
}

func (s *ClientImpl) Publish(
	ctx context.Context,
	req *PublishRequest,
) (*TransactionBytes, error) {
	var resp TransactionBytes
	return &resp, s.http.CallContext(ctx, &resp, publish, req.Sender, req.CompiledModules, req.Dependencies, req.Gas, req.GasBudget)
}

type RequestAddStakeRequest struct {
	Signer    *sui.Address
	Coins     []*sui.ObjectId
	Amount    *sui.BigInt // optional
	Validator *sui.Address
	Gas       *sui.ObjectId // optional
	GasBudget *sui.BigInt
}

func (s *ClientImpl) RequestAddStake(
	ctx context.Context,
	req *RequestAddStakeRequest,
) (*TransactionBytes, error) {
	var resp TransactionBytes
	return &resp, s.http.CallContext(ctx, &resp, requestAddStake, req.Signer, req.Coins, req.Amount, req.Validator, req.Gas, req.GasBudget)
}

type RequestWithdrawStakeRequest struct {
	Signer      *sui.Address
	StakedSuiId *sui.ObjectId
	Gas         *sui.ObjectId // optional
	GasBudget   *sui.BigInt
}

func (s *ClientImpl) RequestWithdrawStake(
	ctx context.Context,
	req *RequestWithdrawStakeRequest,
) (*TransactionBytes, error) {
	var resp TransactionBytes
	return &resp, s.http.CallContext(ctx, &resp, requestWithdrawStake, req.Signer, req.StakedSuiId, req.Gas, req.GasBudget)
}

type SplitCoinRequest struct {
	Signer       *sui.Address
	Coin         *sui.ObjectId
	SplitAmounts []*sui.BigInt
	Gas          *sui.ObjectId // optional
	GasBudget    *sui.BigInt
}

// SplitCoin Creates an unsigned transaction to split a coin object into multiple coins.
// better to replace with unsafe_pay API which consumes less gas
func (s *ClientImpl) SplitCoin(
	ctx context.Context,
	req *SplitCoinRequest,
) (*TransactionBytes, error) {
	resp := TransactionBytes{}
	return &resp, s.http.CallContext(ctx, &resp, splitCoin, req.Signer, req.Coin, req.SplitAmounts, req.Gas, req.GasBudget)
}

type SplitCoinEqualRequest struct {
	Signer     *sui.Address
	Coin       *sui.ObjectId
	SplitCount *sui.BigInt
	Gas        *sui.ObjectId // optional
	GasBudget  *sui.BigInt
}

// SplitCoinEqual Creates an unsigned transaction to split a coin object into multiple equal-size coins.
// better to replace with unsafe_pay API which consumes less gas
func (s *ClientImpl) SplitCoinEqual(
	ctx context.Context,
	req *SplitCoinEqualRequest,
) (*TransactionBytes, error) {
	resp := TransactionBytes{}
	return &resp, s.http.CallContext(ctx, &resp, splitCoinEqual, req.Signer, req.Coin, req.SplitCount, req.Gas, req.GasBudget)
}

type TransferObjectRequest struct {
	Signer    *sui.Address
	ObjectId  *sui.ObjectId
	Gas       *sui.ObjectId // optional
	GasBudget *sui.BigInt
	Recipient *sui.Address
}

// TransferObject Create an unsigned transaction to transfer an object from one address to another. The object's type must allow public transfers
func (s *ClientImpl) TransferObject(
	ctx context.Context,
	req *TransferObjectRequest,
) (*TransactionBytes, error) {
	resp := TransactionBytes{}
	return &resp, s.http.CallContext(ctx, &resp, transferObject, req.Signer, req.ObjectId, req.Gas, req.GasBudget, req.Recipient)
}

type TransferSuiRequest struct {
	Signer    *sui.Address
	ObjectId  *sui.ObjectId
	GasBudget *sui.BigInt
	Recipient *sui.Address
	Amount    *sui.BigInt // optional
}

// TransferSui Create an unsigned transaction to send SUI coin object to a Sui address. The SUI object is also used as the gas object.
func (s *ClientImpl) TransferSui(
	ctx context.Context,
	req *TransferSuiRequest,
) (*TransactionBytes, error) {
	resp := TransactionBytes{}
	return &resp, s.http.CallContext(ctx, &resp, transferSui, req.Signer, req.ObjectId, req.GasBudget, req.Recipient, req.Amount)
}
