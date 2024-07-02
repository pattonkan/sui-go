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

func (s *ImplSuiAPI) GetCheckpoint(ctx context.Context, checkpointId *models.BigInt) (*models.Checkpoint, error) {
	var resp models.Checkpoint
	return &resp, s.http.CallContext(ctx, &resp, getCheckpoint, checkpointId)
}

func (s *ImplSuiAPI) GetCheckpoints(ctx context.Context, req *models.GetCheckpointsRequest) (*models.CheckpointPage, error) {
	var resp models.CheckpointPage
	return &resp, s.http.CallContext(ctx, &resp, getCheckpoints, req.Cursor, req.Limit, req.DescendingOrder)
}

func (s *ImplSuiAPI) GetEvents(ctx context.Context, digest *sui_types.TransactionDigest) ([]*models.SuiEvent, error) {
	var resp []*models.SuiEvent
	return resp, s.http.CallContext(ctx, &resp, getEvents, digest)
}

func (s *ImplSuiAPI) GetLatestCheckpointSequenceNumber(ctx context.Context) (string, error) {
	var resp string
	return resp, s.http.CallContext(ctx, &resp, getLatestCheckpointSequenceNumber)
}

func (s *ImplSuiAPI) GetObject(ctx context.Context, req *models.GetObjectRequest) (*models.SuiObjectResponse, error) {
	var resp models.SuiObjectResponse
	return &resp, s.http.CallContext(ctx, &resp, getObject, req.ObjectID, req.Options)
}

func (s *ImplSuiAPI) GetProtocolConfig(
	ctx context.Context,
	version *models.BigInt, // optional
) (*models.ProtocolConfig, error) {
	var resp models.ProtocolConfig
	return &resp, s.http.CallContext(ctx, &resp, getProtocolConfig, version)
}

func (s *ImplSuiAPI) GetTotalTransactionBlocks(ctx context.Context) (string, error) {
	var resp string
	return resp, s.http.CallContext(ctx, &resp, getTotalTransactionBlocks)
}

func (s *ImplSuiAPI) GetTransactionBlock(ctx context.Context, req *models.GetTransactionBlockRequest) (*models.SuiTransactionBlockResponse, error) {
	resp := models.SuiTransactionBlockResponse{}
	return &resp, s.http.CallContext(ctx, &resp, getTransactionBlock, req.Digest, req.Options)
}

func (s *ImplSuiAPI) MultiGetObjects(ctx context.Context, req *models.MultiGetObjectsRequest) ([]models.SuiObjectResponse, error) {
	var resp []models.SuiObjectResponse
	return resp, s.http.CallContext(ctx, &resp, multiGetObjects, req.ObjectIDs, req.Options)
}

func (s *ImplSuiAPI) MultiGetTransactionBlocks(
	ctx context.Context,
	req *models.MultiGetTransactionBlocksRequest,
) ([]*models.SuiTransactionBlockResponse, error) {
	resp := []*models.SuiTransactionBlockResponse{}
	return resp, s.http.CallContext(ctx, &resp, multiGetTransactionBlocks, req.Digests, req.Options)
}

func (s *ImplSuiAPI) TryGetPastObject(
	ctx context.Context,
	req *models.TryGetPastObjectRequest,
) (*models.SuiPastObjectResponse, error) {
	var resp models.SuiPastObjectResponse
	return &resp, s.http.CallContext(ctx, &resp, tryGetPastObject, req.ObjectID, req.Version, req.Options)
}

func (s *ImplSuiAPI) TryMultiGetPastObjects(
	ctx context.Context,
	req *models.TryMultiGetPastObjectsRequest,
) ([]*models.SuiPastObjectResponse, error) {
	var resp []*models.SuiPastObjectResponse
	return resp, s.http.CallContext(ctx, &resp, tryMultiGetPastObjects, req.PastObjects, req.Options)
}
