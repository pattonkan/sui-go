package sui

import (
	"context"
	"encoding/json"
	"errors"
	"log"

	"github.com/howjmay/sui-go/sui_types"
	"github.com/howjmay/sui-go/types"
	"github.com/tidwall/gjson"
)

func (s *ImplSuiAPI) GetDynamicFieldObject(
	ctx context.Context, parentObjectID *sui_types.ObjectID,
	name sui_types.DynamicFieldName,
) (*types.SuiObjectResponse, error) {
	var resp types.SuiObjectResponse
	return &resp, s.http.CallContext(ctx, &resp, getDynamicFieldObject, parentObjectID, name)
}

func (s *ImplSuiAPI) GetDynamicFields(
	ctx context.Context, parentObjectID *sui_types.ObjectID, cursor *sui_types.ObjectID,
	limit *uint,
) (*types.DynamicFieldPage, error) {
	var resp types.DynamicFieldPage
	return &resp, s.http.CallContext(ctx, &resp, getDynamicFields, parentObjectID, cursor, limit)
}

// address : <SuiAddress> - the owner's Sui address
// query : <ObjectResponseQuery> - the objects query criteria.
// cursor : <CheckpointedObjectID> - An optional paging cursor. If provided, the query will start from the next item after the specified cursor. Default to start from the first item if not specified.
// limit : <uint> - Max number of items returned per page, default to [QUERY_MAX_RESULT_LIMIT_OBJECTS] if is 0
func (s *ImplSuiAPI) GetOwnedObjects(
	ctx context.Context,
	address *sui_types.SuiAddress,
	query *types.SuiObjectResponseQuery,
	cursor *types.CheckpointedObjectID,
	limit *uint,
) (*types.ObjectsPage, error) {
	var resp types.ObjectsPage
	return &resp, s.http.CallContext(ctx, &resp, getOwnedObjects, address, query, cursor, limit)
}

func (s *ImplSuiAPI) QueryEvents(
	ctx context.Context, query types.EventFilter, cursor *types.EventId, limit *uint,
	descendingOrder bool,
) (*types.EventPage, error) {
	var resp types.EventPage
	return &resp, s.http.CallContext(ctx, &resp, queryEvents, query, cursor, limit, descendingOrder)
}

func (s *ImplSuiAPI) QueryTransactionBlocks(
	ctx context.Context, query types.SuiTransactionBlockResponseQuery,
	cursor *sui_types.TransactionDigest, limit *uint, descendingOrder bool,
) (*types.TransactionBlocksPage, error) {
	resp := types.TransactionBlocksPage{}
	return &resp, s.http.CallContext(ctx, &resp, queryTransactionBlocks, query, cursor, limit, descendingOrder)
}

func (s *ImplSuiAPI) ResolveNameServiceAddress(ctx context.Context, suiName string) (*sui_types.SuiAddress, error) {
	var resp sui_types.SuiAddress
	err := s.http.CallContext(ctx, &resp, resolveNameServiceAddress, suiName)
	if err != nil && err.Error() == "nil address" {
		return nil, errors.New("sui name not found")
	}
	return &resp, nil
}

func (s *ImplSuiAPI) ResolveNameServiceNames(ctx context.Context,
	owner *sui_types.SuiAddress, cursor *sui_types.ObjectID, limit *uint) (*types.SuiNamePage, error) {
	var resp types.SuiNamePage
	return &resp, s.http.CallContext(ctx, &resp, resolveNameServiceNames, owner, cursor, limit)
}

func (s *ImplSuiAPI) SubscribeEvent(ctx context.Context, filter types.EventFilter, resultCh chan types.SuiEvent) error {
	resp := make(chan []byte, 10)
	err := s.websocket.CallContext(ctx, resp, subscribeEvent, filter)
	if err != nil {
		return err
	}
	go func() {
		for messageData := range resp {
			var result types.SuiEvent
			if gjson.ParseBytes(messageData).Get("error").Exists() {
				log.Fatal(gjson.ParseBytes(messageData).Get("error").String())
			}

			err := json.Unmarshal([]byte(gjson.ParseBytes(messageData).Get("params.result").String()), &result)
			if err != nil {
				log.Fatal(err)
			}

			resultCh <- result
		}

	}()
	return nil
}

// TODO SubscribeTransaction
