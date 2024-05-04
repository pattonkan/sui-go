package sui

import (
	"context"

	"github.com/howjmay/sui-go/models"
	"github.com/howjmay/sui-go/sui_types"
)

// TODO getChainIdentifier

// TODO getCheckpoint

// TODO getCheckpoints

func (s *ImplSuiAPI) GetEvents(ctx context.Context, digest sui_types.TransactionDigest) ([]models.SuiEvent, error) {
	var resp []models.SuiEvent
	return resp, s.http.CallContext(ctx, &resp, getEvents, digest)
}

func (s *ImplSuiAPI) GetLatestCheckpointSequenceNumber(ctx context.Context) (string, error) {
	var resp string
	return resp, s.http.CallContext(ctx, &resp, getLatestCheckpointSequenceNumber)
}

// TODO getLoadedChildObjects

func (s *ImplSuiAPI) GetObject(
	ctx context.Context,
	objID *sui_types.ObjectID,
	options *models.SuiObjectDataOptions,
) (*models.SuiObjectResponse, error) {
	var resp models.SuiObjectResponse
	return &resp, s.http.CallContext(ctx, &resp, getObject, objID, options)
}

// TODO getProtocolConfig

func (s *ImplSuiAPI) GetTotalTransactionBlocks(ctx context.Context) (string, error) {
	var resp string
	return resp, s.http.CallContext(ctx, &resp, getTotalTransactionBlocks)
}

func (s *ImplSuiAPI) GetTransactionBlock(
	ctx context.Context,
	digest sui_types.TransactionDigest,
	options models.SuiTransactionBlockResponseOptions,
) (*models.SuiTransactionBlockResponse, error) {
	resp := models.SuiTransactionBlockResponse{}
	return &resp, s.http.CallContext(ctx, &resp, getTransactionBlock, digest, options)
}

func (s *ImplSuiAPI) MultiGetObjects(
	ctx context.Context,
	objIDs []sui_types.ObjectID,
	options *models.SuiObjectDataOptions,
) ([]models.SuiObjectResponse, error) {
	var resp []models.SuiObjectResponse
	return resp, s.http.CallContext(ctx, &resp, multiGetObjects, objIDs, options)
}

// TODO multiGetTransactionBlocks

func (s *ImplSuiAPI) TryGetPastObject(
	ctx context.Context,
	objectId *sui_types.ObjectID,
	version uint64,
	options *models.SuiObjectDataOptions,
) (*models.SuiPastObjectResponse, error) {
	var resp models.SuiPastObjectResponse
	return &resp, s.http.CallContext(ctx, &resp, tryGetPastObject, objectId, version, options)
}

// TODO tryMultiGetPastObjects
