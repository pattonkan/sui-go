package models

import (
	"github.com/howjmay/sui-go/sui_signer"
	"github.com/howjmay/sui-go/sui_types"
)

type DevInspectTransactionBlockRequest struct {
	SenderAddress *sui_types.SuiAddress
	TxKindBytes   sui_types.Base64Data
	GasPrice      *BigInt // optional
	Epoch         *uint64 // optional
	// additional_args // optional // FIXME
}

type ExecuteTransactionBlockRequest struct {
	TxDataBytes sui_types.Base64Data
	Signatures  []*sui_signer.Signature
	Options     *SuiTransactionBlockResponseOptions // optional
	RequestType ExecuteTransactionRequestType       // optional
}
