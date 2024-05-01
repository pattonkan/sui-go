package client

import (
	"context"
	"errors"

	"github.com/howjmay/go-sui-sdk/sui_types"
	"github.com/howjmay/go-sui-sdk/types"
)

func (c *Client) GetDynamicFieldObject(
	ctx context.Context, parentObjectID *sui_types.ObjectID,
	name sui_types.DynamicFieldName,
) (*types.SuiObjectResponse, error) {
	var resp types.SuiObjectResponse
	return &resp, c.CallContext(ctx, &resp, getDynamicFieldObject, parentObjectID, name)
}

func (c *Client) GetDynamicFields(
	ctx context.Context, parentObjectID *sui_types.ObjectID, cursor *sui_types.ObjectID,
	limit *uint,
) (*types.DynamicFieldPage, error) {
	var resp types.DynamicFieldPage
	return &resp, c.CallContext(ctx, &resp, getDynamicFields, parentObjectID, cursor, limit)
}

// address : <SuiAddress> - the owner's Sui address
// query : <ObjectResponseQuery> - the objects query criteria.
// cursor : <CheckpointedObjectID> - An optional paging cursor. If provided, the query will start from the next item after the specified cursor. Default to start from the first item if not specified.
// limit : <uint> - Max number of items returned per page, default to [QUERY_MAX_RESULT_LIMIT_OBJECTS] if is 0
func (c *Client) GetOwnedObjects(
	ctx context.Context,
	address *sui_types.SuiAddress,
	query *types.SuiObjectResponseQuery,
	cursor *types.CheckpointedObjectID,
	limit *uint,
) (*types.ObjectsPage, error) {
	var resp types.ObjectsPage
	return &resp, c.CallContext(ctx, &resp, getOwnedObjects, address, query, cursor, limit)
}

func (c *Client) QueryEvents(
	ctx context.Context, query types.EventFilter, cursor *types.EventId, limit *uint,
	descendingOrder bool,
) (*types.EventPage, error) {
	var resp types.EventPage
	return &resp, c.CallContext(ctx, &resp, queryEvents, query, cursor, limit, descendingOrder)
}

func (c *Client) QueryTransactionBlocks(
	ctx context.Context, query types.SuiTransactionBlockResponseQuery,
	cursor *sui_types.TransactionDigest, limit *uint, descendingOrder bool,
) (*types.TransactionBlocksPage, error) {
	resp := types.TransactionBlocksPage{}
	return &resp, c.CallContext(ctx, &resp, queryTransactionBlocks, query, cursor, limit, descendingOrder)
}

func (c *Client) ResolveNameServiceAddress(ctx context.Context, suiName string) (*sui_types.SuiAddress, error) {
	var resp sui_types.SuiAddress
	err := c.CallContext(ctx, &resp, resolveNameServiceAddress, suiName)
	if err != nil && err.Error() == "nil address" {
		return nil, errors.New("sui name not found")
	}
	return &resp, nil
}

func (c *Client) ResolveNameServiceNames(ctx context.Context,
	owner *sui_types.SuiAddress, cursor *sui_types.ObjectID, limit *uint) (*types.SuiNamePage, error) {
	var resp types.SuiNamePage
	return &resp, c.CallContext(ctx, &resp, resolveNameServiceNames, owner, cursor, limit)
}

// TODO SubscribeEvent

// TODO SubscribeTransaction
