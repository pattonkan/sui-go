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

type GetCheckpointsRequest struct {
	Cursor          *models.BigInt // optional
	Limit           *uint64        // optional
	DescendingOrder bool
}

func (s *ImplSuiAPI) GetCheckpoints(ctx context.Context, req *GetCheckpointsRequest) (*models.CheckpointPage, error) {
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

type GetObjectRequest struct {
	ObjectID *sui_types.ObjectID
	Options  *models.SuiObjectDataOptions // optional
}

func (s *ImplSuiAPI) GetObject(ctx context.Context, req *GetObjectRequest) (*models.SuiObjectResponse, error) {
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

type GetTransactionBlockRequest struct {
	Digest  *sui_types.TransactionDigest
	Options *models.SuiTransactionBlockResponseOptions // optional
}

func (s *ImplSuiAPI) GetTransactionBlock(ctx context.Context, req *GetTransactionBlockRequest) (*models.SuiTransactionBlockResponse, error) {
	resp := models.SuiTransactionBlockResponse{}
	return &resp, s.http.CallContext(ctx, &resp, getTransactionBlock, req.Digest, req.Options)
}

type MultiGetObjectsRequest struct {
	ObjectIDs []*sui_types.ObjectID
	Options   *models.SuiObjectDataOptions // optional
}

func (s *ImplSuiAPI) MultiGetObjects(ctx context.Context, req *MultiGetObjectsRequest) ([]models.SuiObjectResponse, error) {
	var resp []models.SuiObjectResponse
	return resp, s.http.CallContext(ctx, &resp, multiGetObjects, req.ObjectIDs, req.Options)
}

type MultiGetTransactionBlocksRequest struct {
	Digests []*sui_types.Digest
	Options *models.SuiTransactionBlockResponseOptions // optional
}

func (s *ImplSuiAPI) MultiGetTransactionBlocks(
	ctx context.Context,
	req *MultiGetTransactionBlocksRequest,
) ([]*models.SuiTransactionBlockResponse, error) {
	resp := []*models.SuiTransactionBlockResponse{}
	return resp, s.http.CallContext(ctx, &resp, multiGetTransactionBlocks, req.Digests, req.Options)
}

type TryGetPastObjectRequest struct {
	ObjectID *sui_types.ObjectID
	Version  uint64
	Options  *models.SuiObjectDataOptions // optional
}

func (s *ImplSuiAPI) TryGetPastObject(
	ctx context.Context,
	req *TryGetPastObjectRequest,
) (*models.SuiPastObjectResponse, error) {
	var resp models.SuiPastObjectResponse
	return &resp, s.http.CallContext(ctx, &resp, tryGetPastObject, req.ObjectID, req.Version, req.Options)
}

type TryMultiGetPastObjectsRequest struct {
	PastObjects []*models.SuiGetPastObjectRequest
	Options     *models.SuiObjectDataOptions // optional
}

func (s *ImplSuiAPI) TryMultiGetPastObjects(
	ctx context.Context,
	req *TryMultiGetPastObjectsRequest,
) ([]*models.SuiPastObjectResponse, error) {
	var resp []*models.SuiPastObjectResponse
	return resp, s.http.CallContext(ctx, &resp, tryMultiGetPastObjects, req.PastObjects, req.Options)
}
