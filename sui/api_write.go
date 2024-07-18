package sui

import (
	"context"

	"github.com/howjmay/sui-go/models"
	"github.com/howjmay/sui-go/sui_signer"
	"github.com/howjmay/sui-go/sui_types"
)

type DevInspectTransactionBlockRequest struct {
	SenderAddress *sui_types.SuiAddress
	TxKindBytes   sui_types.Base64Data
	GasPrice      *models.BigInt // optional
	Epoch         *uint64        // optional
	// additional_args // optional // FIXME
}

// The txKindBytes is `TransactionKind` in base64.
// When a `TransactionData` is given, error `Deserialization error: malformed utf8` will be returned.
// which is different from `DryRunTransaction` and `ExecuteTransactionBlock`
// `DryRunTransaction` and `ExecuteTransactionBlock` takes `TransactionData` in base64
func (s *ImplSuiAPI) DevInspectTransactionBlock(
	ctx context.Context,
	req *DevInspectTransactionBlockRequest,
) (*models.DevInspectResults, error) {
	var resp models.DevInspectResults
	return &resp, s.http.CallContext(ctx, &resp, devInspectTransactionBlock, req.SenderAddress, req.TxKindBytes, req.GasPrice, req.Epoch)
}

func (s *ImplSuiAPI) DryRunTransaction(
	ctx context.Context,
	txDataBytes sui_types.Base64Data,
) (*models.DryRunTransactionBlockResponse, error) {
	var resp models.DryRunTransactionBlockResponse
	return &resp, s.http.CallContext(ctx, &resp, dryRunTransactionBlock, txDataBytes)
}

type ExecuteTransactionBlockRequest struct {
	TxDataBytes sui_types.Base64Data
	Signatures  []*sui_signer.Signature
	Options     *models.SuiTransactionBlockResponseOptions // optional
	RequestType models.ExecuteTransactionRequestType       // optional
}

func (s *ImplSuiAPI) ExecuteTransactionBlock(
	ctx context.Context,
	req *ExecuteTransactionBlockRequest,
) (*models.SuiTransactionBlockResponse, error) {
	resp := models.SuiTransactionBlockResponse{}
	return &resp, s.http.CallContext(ctx, &resp, executeTransactionBlock, req.TxDataBytes, req.Signatures, req.Options, req.RequestType)
}
