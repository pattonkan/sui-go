package client

import (
	"context"

	"github.com/howjmay/go-sui-sdk/sui_types"
	"github.com/howjmay/go-sui-sdk/types"
)

// TODO getChainIdentifier

// TODO getCheckpoint

// TODO getCheckpoints

func (c *Client) GetEvents(ctx context.Context, digest sui_types.TransactionDigest) ([]types.SuiEvent, error) {
	var resp []types.SuiEvent
	return resp, c.CallContext(ctx, &resp, getEvents, digest)
}

func (c *Client) GetLatestCheckpointSequenceNumber(ctx context.Context) (string, error) {
	var resp string
	return resp, c.CallContext(ctx, &resp, getLatestCheckpointSequenceNumber)
}

// TODO getLoadedChildObjects

func (c *Client) GetObject(
	ctx context.Context,
	objID *sui_types.ObjectID,
	options *types.SuiObjectDataOptions,
) (*types.SuiObjectResponse, error) {
	var resp types.SuiObjectResponse
	return &resp, c.CallContext(ctx, &resp, getObject, objID, options)
}

// TODO getProtocolConfig

func (c *Client) GetTotalTransactionBlocks(ctx context.Context) (string, error) {
	var resp string
	return resp, c.CallContext(ctx, &resp, getTotalTransactionBlocks)
}

func (c *Client) GetTransactionBlock(
	ctx context.Context,
	digest sui_types.TransactionDigest,
	options types.SuiTransactionBlockResponseOptions,
) (*types.SuiTransactionBlockResponse, error) {
	resp := types.SuiTransactionBlockResponse{}
	return &resp, c.CallContext(ctx, &resp, getTransactionBlock, digest, options)
}

func (c *Client) MultiGetObjects(
	ctx context.Context,
	objIDs []sui_types.ObjectID,
	options *types.SuiObjectDataOptions,
) ([]types.SuiObjectResponse, error) {
	var resp []types.SuiObjectResponse
	return resp, c.CallContext(ctx, &resp, multiGetObjects, objIDs, options)
}

// TODO multiGetTransactionBlocks

func (c *Client) TryGetPastObject(
	ctx context.Context,
	objectId *sui_types.ObjectID,
	version uint64,
	options *types.SuiObjectDataOptions,
) (*types.SuiPastObjectResponse, error) {
	var resp types.SuiPastObjectResponse
	return &resp, c.CallContext(ctx, &resp, tryGetPastObject, objectId, version, options)
}

// TODO tryMultiGetPastObjects
