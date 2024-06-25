package sui

import (
	"context"

	"github.com/howjmay/sui-go/models"
	"github.com/howjmay/sui-go/sui_signer"
	"github.com/howjmay/sui-go/sui_types"
)

// the txKindBytes is TransactionKind in base64
// which is different from `DryRunTransaction` and `ExecuteTransactionBlock`
// `DryRunTransaction` and `ExecuteTransactionBlock` takes `TransactionData` in base64
func (s *ImplSuiAPI) DevInspectTransactionBlock(
	ctx context.Context,
	senderAddress *sui_types.SuiAddress,
	txKindBytes sui_types.Base64Data,
	gasPrice *models.BigInt, // optional
	epoch *uint64, // optional
) (*models.DevInspectResults, error) {
	var resp models.DevInspectResults
	return &resp, s.http.CallContext(ctx, &resp, devInspectTransactionBlock, senderAddress, txKindBytes, gasPrice, epoch)
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
	txDataBytes sui_types.Base64Data,
	signatures []*sui_signer.Signature,
	options *models.SuiTransactionBlockResponseOptions,
	requestType models.ExecuteTransactionRequestType,
) (*models.SuiTransactionBlockResponse, error) {
	resp := models.SuiTransactionBlockResponse{}
	return &resp, s.http.CallContext(ctx, &resp, executeTransactionBlock, txDataBytes, signatures, options, requestType)
}
