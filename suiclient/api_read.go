package suiclient

import (
	"context"

	"github.com/pattonkan/sui-go/sui"
)

func (s *ClientImpl) GetChainIdentifier(ctx context.Context) (string, error) {
	var resp string
	return resp, s.http.CallContext(ctx, &resp, getChainIdentifier)
}

func (s *ClientImpl) GetCheckpoint(ctx context.Context, checkpointId *sui.BigInt) (*Checkpoint, error) {
	var resp Checkpoint
	return &resp, s.http.CallContext(ctx, &resp, getCheckpoint, checkpointId)
}

type GetCheckpointsRequest struct {
	Cursor          *sui.BigInt // optional
	Limit           *uint64     // optional
	DescendingOrder bool
}

func (s *ClientImpl) GetCheckpoints(ctx context.Context, req *GetCheckpointsRequest) (*CheckpointPage, error) {
	var resp CheckpointPage
	return &resp, s.http.CallContext(ctx, &resp, getCheckpoints, req.Cursor, req.Limit, req.DescendingOrder)
}

func (s *ClientImpl) GetEvents(ctx context.Context, digest *sui.TransactionDigest) ([]*Event, error) {
	var resp []*Event
	return resp, s.http.CallContext(ctx, &resp, getEvents, digest)
}

func (s *ClientImpl) GetLatestCheckpointSequenceNumber(ctx context.Context) (string, error) {
	var resp string
	return resp, s.http.CallContext(ctx, &resp, getLatestCheckpointSequenceNumber)
}

type GetObjectRequest struct {
	ObjectId *sui.ObjectId
	Options  *SuiObjectDataOptions // optional
}

func (s *ClientImpl) GetObject(ctx context.Context, req *GetObjectRequest) (*SuiObjectResponse, error) {
	var resp SuiObjectResponse
	return &resp, s.http.CallContext(ctx, &resp, getObject, req.ObjectId, req.Options)
}

func (s *ClientImpl) GetProtocolConfig(
	ctx context.Context,
	version *sui.BigInt, // optional
) (*ProtocolConfig, error) {
	var resp ProtocolConfig
	return &resp, s.http.CallContext(ctx, &resp, getProtocolConfig, version)
}

func (s *ClientImpl) GetTotalTransactionBlocks(ctx context.Context) (string, error) {
	var resp string
	return resp, s.http.CallContext(ctx, &resp, getTotalTransactionBlocks)
}

type GetTransactionBlockRequest struct {
	Digest  *sui.TransactionDigest
	Options *SuiTransactionBlockResponseOptions // optional
}

func (s *ClientImpl) GetTransactionBlock(ctx context.Context, req *GetTransactionBlockRequest) (*SuiTransactionBlockResponse, error) {
	resp := SuiTransactionBlockResponse{}
	return &resp, s.http.CallContext(ctx, &resp, getTransactionBlock, req.Digest, req.Options)
}

type MultiGetObjectsRequest struct {
	ObjectIds []*sui.ObjectId
	Options   *SuiObjectDataOptions // optional
}

func (s *ClientImpl) MultiGetObjects(ctx context.Context, req *MultiGetObjectsRequest) ([]SuiObjectResponse, error) {
	var resp []SuiObjectResponse
	return resp, s.http.CallContext(ctx, &resp, multiGetObjects, req.ObjectIds, req.Options)
}

type MultiGetTransactionBlocksRequest struct {
	Digests []*sui.Digest
	Options *SuiTransactionBlockResponseOptions // optional
}

func (s *ClientImpl) MultiGetTransactionBlocks(
	ctx context.Context,
	req *MultiGetTransactionBlocksRequest,
) ([]*SuiTransactionBlockResponse, error) {
	resp := []*SuiTransactionBlockResponse{}
	return resp, s.http.CallContext(ctx, &resp, multiGetTransactionBlocks, req.Digests, req.Options)
}

type TryGetPastObjectRequest struct {
	ObjectId *sui.ObjectId
	Version  uint64
	Options  *SuiObjectDataOptions // optional
}

func (s *ClientImpl) TryGetPastObject(
	ctx context.Context,
	req *TryGetPastObjectRequest,
) (*SuiPastObjectResponse, error) {
	var resp SuiPastObjectResponse
	return &resp, s.http.CallContext(ctx, &resp, tryGetPastObject, req.ObjectId, req.Version, req.Options)
}

type TryMultiGetPastObjectsRequest struct {
	PastObjects []*SuiGetPastObjectRequest
	Options     *SuiObjectDataOptions // optional
}

func (s *ClientImpl) TryMultiGetPastObjects(
	ctx context.Context,
	req *TryMultiGetPastObjectsRequest,
) ([]*SuiPastObjectResponse, error) {
	var resp []*SuiPastObjectResponse
	return resp, s.http.CallContext(ctx, &resp, tryMultiGetPastObjects, req.PastObjects, req.Options)
}
