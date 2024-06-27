package sui

import (
	"context"

	"github.com/howjmay/sui-go/models"
	"github.com/howjmay/sui-go/sui_types"
)

// The txKindBytes is `TransactionKind` in base64.
// When a `TransactionData` is given, error `Deserialization error: malformed utf8` will be returned.
// which is different from `DryRunTransaction` and `ExecuteTransactionBlock`
// `DryRunTransaction` and `ExecuteTransactionBlock` takes `TransactionData` in base64
func (s *ImplSuiAPI) DevInspectTransactionBlock(
	ctx context.Context,
	req *models.DevInspectTransactionBlockRequest,
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

func (s *ImplSuiAPI) ExecuteTransactionBlock(
	ctx context.Context,
	req *models.ExecuteTransactionBlockRequest,
) (*models.SuiTransactionBlockResponse, error) {
	resp := models.SuiTransactionBlockResponse{}
	return &resp, s.http.CallContext(ctx, &resp, executeTransactionBlock, req.TxDataBytes, req.Signatures, req.Options, req.RequestType)
}
