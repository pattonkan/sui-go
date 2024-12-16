package suiclient

import (
	"context"

	"github.com/pattonkan/sui-go/sui"
	"github.com/pattonkan/sui-go/suisigner"
)

type DevInspectTransactionBlockRequest struct {
	SenderAddress *sui.Address
	TxKindBytes   sui.Base64Data
	GasPrice      *sui.BigInt // optional
	Epoch         *uint64     // optional
	// additional_args // optional // FIXME
}

type DevInspectTransactionBlockResponse struct {
	Effects WrapperTaggedJson[SuiTransactionBlockEffects] `json:"effects"`
	Events  []Event                                       `json:"events"`
	Results []ExecutionResultType                         `json:"results,omitempty"`
	Error   string                                        `json:"error,omitempty"`
}

// The txKindBytes is `TransactionKind` in base64.
// When a `TransactionData` is given, error `Deserialization error: malformed utf8` will be returned.
// which is different from `DryRunTransaction` and `ExecuteTransactionBlock`
// `DryRunTransaction` and `ExecuteTransactionBlock` takes `TransactionData` in base64
func (s *ClientImpl) DevInspectTransactionBlock(
	ctx context.Context,
	req *DevInspectTransactionBlockRequest,
) (*DevInspectTransactionBlockResponse, error) {
	var resp DevInspectTransactionBlockResponse
	return &resp, s.http.CallContext(ctx, &resp, devInspectTransactionBlock, req.SenderAddress, req.TxKindBytes, req.GasPrice, req.Epoch)
}

func (s *ClientImpl) DryRunTransaction(
	ctx context.Context,
	txDataBytes sui.Base64Data,
) (*DryRunTransactionBlockResponse, error) {
	var resp DryRunTransactionBlockResponse
	return &resp, s.http.CallContext(ctx, &resp, dryRunTransactionBlock, txDataBytes)
}

type ExecuteTransactionBlockRequest struct {
	TxDataBytes sui.Base64Data
	Signatures  []*suisigner.Signature
	Options     *SuiTransactionBlockResponseOptions // optional
	RequestType ExecuteTransactionRequestType       // optional
}

func (s *ClientImpl) ExecuteTransactionBlock(
	ctx context.Context,
	req *ExecuteTransactionBlockRequest,
) (*SuiTransactionBlockResponse, error) {
	resp := SuiTransactionBlockResponse{}
	return &resp, s.http.CallContext(ctx, &resp, executeTransactionBlock, req.TxDataBytes, req.Signatures, req.Options, req.RequestType)
}
