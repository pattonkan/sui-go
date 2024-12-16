package suiclient

import (
	"context"
	"encoding/json"
	"errors"
	"log"

	"github.com/pattonkan/sui-go/sui"
)

type GetDynamicFieldObjectRequest struct {
	ParentObjectId *sui.ObjectId
	Name           *DynamicFieldName
}

func (s *ClientImpl) GetDynamicFieldObject(
	ctx context.Context,
	req *GetDynamicFieldObjectRequest,
) (*SuiObjectResponse, error) {
	var resp SuiObjectResponse
	return &resp, s.http.CallContext(ctx, &resp, getDynamicFieldObject, req.ParentObjectId, req.Name)
}

type GetDynamicFieldsRequest struct {
	ParentObjectId *sui.ObjectId
	Cursor         *sui.ObjectId // optional
	Limit          *uint         // optional
}

func (s *ClientImpl) GetDynamicFields(
	ctx context.Context,
	req *GetDynamicFieldsRequest,
) (*DynamicFieldPage, error) {
	var resp DynamicFieldPage
	return &resp, s.http.CallContext(ctx, &resp, getDynamicFields, req.ParentObjectId, req.Cursor, req.Limit)
}

type GetOwnedObjectsRequest struct {
	Address *sui.Address
	Query   *SuiObjectResponseQuery // optional
	Cursor  *CheckpointedObjectId   // optional
	Limit   *uint                   // optional
}

// address : <Address> - the owner's Sui address
// query : <ObjectResponseQuery> - the objects query criteria.
// cursor : <CheckpointedObjectId> - An optional paging cursor. If provided, the query will start from the next item after the specified cursor. Default to start from the first item if not specified.
// limit : <uint> - Max number of items returned per page, default to [QUERY_MAX_RESULT_LIMIT_OBJECTS] if is 0
func (s *ClientImpl) GetOwnedObjects(
	ctx context.Context,
	req *GetOwnedObjectsRequest,
) (*ObjectsPage, error) {
	var resp ObjectsPage
	return &resp, s.http.CallContext(ctx, &resp, getOwnedObjects, req.Address, req.Query, req.Cursor, req.Limit)
}

type QueryEventsRequest struct {
	Query           *EventFilter
	Cursor          *EventId // optional
	Limit           *uint    // optional
	DescendingOrder bool     // optional
}

func (s *ClientImpl) QueryEvents(
	ctx context.Context,
	req *QueryEventsRequest,
) (*EventPage, error) {
	var resp EventPage
	return &resp, s.http.CallContext(ctx, &resp, queryEvents, req.Query, req.Cursor, req.Limit, req.DescendingOrder)
}

type QueryTransactionBlocksRequest struct {
	Query           *TransactionBlockResponseQuery
	Cursor          *sui.TransactionDigest // optional
	Limit           *uint                  // optional
	DescendingOrder bool                   // optional
}

func (s *ClientImpl) QueryTransactionBlocks(
	ctx context.Context,
	req *QueryTransactionBlocksRequest,
) (*TransactionBlocksPage, error) {
	resp := TransactionBlocksPage{}
	return &resp, s.http.CallContext(ctx, &resp, queryTransactionBlocks, req.Query, req.Cursor, req.Limit, req.DescendingOrder)
}

type ResolveNameServiceNamesRequest struct {
	Owner  *sui.Address
	Cursor *sui.ObjectId // optional
	Limit  *uint         // optional
}

func (s *ClientImpl) ResolveNameServiceAddress(ctx context.Context, suiName string) (*sui.Address, error) {
	var resp sui.Address
	err := s.http.CallContext(ctx, &resp, resolveNameServiceAddress, suiName)
	if err != nil && err.Error() == "nil address" {
		return nil, errors.New("sui name not found")
	}
	return &resp, nil
}

func (s *ClientImpl) ResolveNameServiceNames(
	ctx context.Context,
	req *ResolveNameServiceNamesRequest,
) (*SuiNamePage, error) {
	var resp SuiNamePage
	return &resp, s.http.CallContext(ctx, &resp, resolveNameServiceNames, req.Owner, req.Cursor, req.Limit)
}

func (s *ClientImpl) SubscribeEvent(
	ctx context.Context,
	filter *EventFilter,
	resultCh chan Event,
) error {
	resp := make(chan []byte, 10)
	err := s.websocket.CallContext(ctx, resp, subscribeEvent, filter)
	if err != nil {
		return err
	}
	go func() {
		for messageData := range resp {
			var result Event
			if err := json.Unmarshal(messageData, &result); err != nil {
				log.Fatal(err)
			}

			resultCh <- result
		}

	}()
	return nil
}

func (s *ClientImpl) SubscribeTransaction(
	ctx context.Context,
	filter *TransactionFilter,
	resultCh chan WrapperTaggedJson[SuiTransactionBlockEffects],
) error {
	resp := make(chan []byte, 10)
	err := s.websocket.CallContext(ctx, resp, subscribeTransaction, filter)
	if err != nil {
		return err
	}
	go func() {
		for messageData := range resp {
			var result WrapperTaggedJson[SuiTransactionBlockEffects]
			if err := json.Unmarshal(messageData, &result); err != nil {
				log.Fatal(err)
			}

			resultCh <- result
		}

	}()
	return nil
}
