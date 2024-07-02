package models

import "github.com/howjmay/sui-go/sui_types"

type BatchTransactionRequest struct {
	Signer    *sui_types.SuiAddress
	TxnParams []RPCTransactionRequestParams
	Gas       *sui_types.ObjectID // optional
	GasBudget *BigInt
	// txnBuilderMode // optional // FIXME SuiTransactionBlockBuilderMode
}

type RPCTransactionRequestParams struct {
	TransferObjectRequestParams *TransferObjectParams `json:"transferObjectRequestParams,omitempty"`
	MoveCallRequestParams       *MoveCallParams       `json:"moveCallRequestParams,omitempty"`
}

type TransferObjectParams struct {
	Recipient *sui_types.SuiAddress `json:"recipient,omitempty"`
	ObjectId  *sui_types.ObjectID   `json:"objectId,omitempty"`
}

type SuiTypeTag string
type SuiJsonValue string

type MoveCallParams struct {
	PackageObjectId *sui_types.ObjectID  `json:"packageObjectId,omitempty"`
	Module          sui_types.Identifier `json:"module,omitempty"`
	Function        sui_types.Identifier `json:"function,omitempty"`
	TypeArguments   []SuiTypeTag         `json:"typeArguments,omitempty"`
	// MoveCall arguments in JSON format
	Arguments []SuiJsonValue `json:"arguments,omitempty"`
}

type MergeCoinsRequest struct {
	Signer      *sui_types.SuiAddress `json:"signer,omitempty"`
	PrimaryCoin *sui_types.ObjectID   `json:"primaryCoin,omitempty"`
	CoinToMerge *sui_types.ObjectID   `json:"coinToMerge,omitempty"`
	Gas         *sui_types.ObjectID   `json:"gas,omitempty"` // optional
	GasBudget   *BigInt               `json:"gasBudget,omitempty"`
}

type MoveCallRequest struct {
	Signer    *sui_types.SuiAddress
	PackageID *sui_types.PackageID
	Module    string
	Function  string
	TypeArgs  []string
	Arguments []any
	Gas       *sui_types.ObjectID // optional
	GasBudget *BigInt
	// txnBuilderMode // optional // FIXME SuiTransactionBlockBuilderMode
}

type PayRequest struct {
	Signer     *sui_types.SuiAddress
	InputCoins []*sui_types.ObjectID
	Recipients []*sui_types.SuiAddress
	Amount     []*BigInt
	Gas        *sui_types.ObjectID // optional
	GasBudget  *BigInt
}

type PayAllSuiRequest struct {
	Signer     *sui_types.SuiAddress
	Recipient  *sui_types.SuiAddress
	InputCoins []*sui_types.ObjectID
	GasBudget  *BigInt
}

type PaySuiRequest struct {
	Signer     *sui_types.SuiAddress
	InputCoins []*sui_types.ObjectID
	Recipients []*sui_types.SuiAddress
	Amount     []*BigInt
	GasBudget  *BigInt
}

type PublishRequest struct {
	Sender          *sui_types.SuiAddress
	CompiledModules []*sui_types.Base64Data
	Dependencies    []*sui_types.ObjectID
	Gas             *sui_types.ObjectID // optional
	GasBudget       *BigInt
}

type RequestAddStakeRequest struct {
	Signer    *sui_types.SuiAddress
	Coins     []*sui_types.ObjectID
	Amount    *BigInt // optional
	Validator *sui_types.SuiAddress
	Gas       *sui_types.ObjectID // optional
	GasBudget *BigInt
}

type RequestWithdrawStakeRequest struct {
	Signer      *sui_types.SuiAddress
	StakedSuiId *sui_types.ObjectID
	Gas         *sui_types.ObjectID // optional
	GasBudget   *BigInt
}

type SplitCoinRequest struct {
	Signer       *sui_types.SuiAddress
	Coin         *sui_types.ObjectID
	SplitAmounts []*BigInt
	Gas          *sui_types.ObjectID // optional
	GasBudget    *BigInt
}

type SplitCoinEqualRequest struct {
	Signer     *sui_types.SuiAddress
	Coin       *sui_types.ObjectID
	SplitCount *BigInt
	Gas        *sui_types.ObjectID // optional
	GasBudget  *BigInt
}

type TransferObjectRequest struct {
	Signer    *sui_types.SuiAddress
	ObjectID  *sui_types.ObjectID
	Gas       *sui_types.ObjectID // optional
	GasBudget *BigInt
	Recipient *sui_types.SuiAddress
}

type TransferSuiRequest struct {
	Signer    *sui_types.SuiAddress
	ObjectID  *sui_types.ObjectID
	GasBudget *BigInt
	Recipient *sui_types.SuiAddress
	Amount    *BigInt // optional
}
