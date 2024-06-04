package sui

import (
	"context"

	"github.com/howjmay/sui-go/models"
	"github.com/howjmay/sui-go/sui_types"
)

func (s *ImplSuiAPI) GetChainIdentifier(ctx context.Context) (string, error) {
	var resp string
	return resp, s.http.CallContext(ctx, &resp, getChainIdentifier)
}

func (s *ImplSuiAPI) GetCheckpoint(ctx context.Context, checkpointId models.SafeSuiBigInt[uint64]) (*models.Checkpoint, error) {
	var resp models.Checkpoint
	return &resp, s.http.CallContext(ctx, &resp, getCheckpoint, checkpointId)
}

func (s *ImplSuiAPI) GetCheckpoints(ctx context.Context, cursor *models.SafeSuiBigInt[uint64], limit *uint64, descendingOrder bool) (*models.CheckpointPage, error) {
	var resp models.CheckpointPage
	return &resp, s.http.CallContext(ctx, &resp, getCheckpoints, cursor, limit, descendingOrder)
}

func (s *ImplSuiAPI) GetEvents(ctx context.Context, digest *sui_types.TransactionDigest) ([]*models.SuiEvent, error) {
	var resp []*models.SuiEvent
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

func (s *ImplSuiAPI) GetProtocolConfig(
	ctx context.Context,
	version *models.SafeSuiBigInt[uint64],
) (*models.ProtocolConfig, error) {
	var resp models.ProtocolConfig
	return &resp, s.http.CallContext(ctx, &resp, getProtocolConfig, version)
}

func (s *ImplSuiAPI) GetTotalTransactionBlocks(ctx context.Context) (string, error) {
	var resp string
	return resp, s.http.CallContext(ctx, &resp, getTotalTransactionBlocks)
}

func (s *ImplSuiAPI) GetTransactionBlock(
	ctx context.Context,
	digest *sui_types.TransactionDigest,
	options *models.SuiTransactionBlockResponseOptions,
) (*models.SuiTransactionBlockResponse, error) {
	resp := models.SuiTransactionBlockResponse{}
	return &resp, s.http.CallContext(ctx, &resp, getTransactionBlock, digest, options)
}

func (s *ImplSuiAPI) MultiGetObjects(
	ctx context.Context,
	objIDs []*sui_types.ObjectID,
	options *models.SuiObjectDataOptions,
) ([]models.SuiObjectResponse, error) {
	var resp []models.SuiObjectResponse
	return resp, s.http.CallContext(ctx, &resp, multiGetObjects, objIDs, options)
}

func (s *ImplSuiAPI) MultiGetTransactionBlocks(
	ctx context.Context,
	digests []*sui_types.Digest,
	options *models.SuiTransactionBlockResponseOptions,
) ([]*models.SuiTransactionBlockResponse, error) {
	resp := []*models.SuiTransactionBlockResponse{}
	return resp, s.http.CallContext(ctx, &resp, multiGetTransactionBlocks, digests, options)
}

func (s *ImplSuiAPI) TryGetPastObject(
	ctx context.Context,
	objectId *sui_types.ObjectID,
	version uint64,
	options *models.SuiObjectDataOptions,
) (*models.SuiPastObjectResponse, error) {
	var resp models.SuiPastObjectResponse
	return &resp, s.http.CallContext(ctx, &resp, tryGetPastObject, objectId, version, options)
}

func (s *ImplSuiAPI) TryMultiGetPastObjects(
	ctx context.Context,
	pastObjects []*models.SuiGetPastObjectRequest,
	options *models.SuiObjectDataOptions,
) ([]*models.SuiPastObjectResponse, error) {
	var resp []*models.SuiPastObjectResponse
	return resp, s.http.CallContext(ctx, &resp, tryMultiGetPastObjects, pastObjects, options)
}
