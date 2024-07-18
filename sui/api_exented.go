package sui

import (
	"context"
	"encoding/json"
	"errors"
	"log"

	"github.com/howjmay/sui-go/models"
	"github.com/howjmay/sui-go/sui_types"
	"github.com/howjmay/sui-go/sui_types/serialization"
)

type GetDynamicFieldObjectRequest struct {
	ParentObjectID *sui_types.ObjectID
	Name           *sui_types.DynamicFieldName
}

func (s *ImplSuiAPI) GetDynamicFieldObject(
	ctx context.Context,
	req *GetDynamicFieldObjectRequest,
) (*models.SuiObjectResponse, error) {
	var resp models.SuiObjectResponse
	return &resp, s.http.CallContext(ctx, &resp, getDynamicFieldObject, req.ParentObjectID, req.Name)
}

type GetDynamicFieldsRequest struct {
	ParentObjectID *sui_types.ObjectID
	Cursor         *sui_types.ObjectID // optional
	Limit          *uint               // optional
}

func (s *ImplSuiAPI) GetDynamicFields(
	ctx context.Context,
	req *GetDynamicFieldsRequest,
) (*models.DynamicFieldPage, error) {
	var resp models.DynamicFieldPage
	return &resp, s.http.CallContext(ctx, &resp, getDynamicFields, req.ParentObjectID, req.Cursor, req.Limit)
}

type GetOwnedObjectsRequest struct {
	Address *sui_types.SuiAddress
	Query   *models.SuiObjectResponseQuery // optional
	Cursor  *models.CheckpointedObjectID   // optional
	Limit   *uint                          // optional
}

// address : <SuiAddress> - the owner's Sui address
// query : <ObjectResponseQuery> - the objects query criteria.
// cursor : <CheckpointedObjectID> - An optional paging cursor. If provided, the query will start from the next item after the specified cursor. Default to start from the first item if not specified.
// limit : <uint> - Max number of items returned per page, default to [QUERY_MAX_RESULT_LIMIT_OBJECTS] if is 0
func (s *ImplSuiAPI) GetOwnedObjects(
	ctx context.Context,
	req *GetOwnedObjectsRequest,
) (*models.ObjectsPage, error) {
	var resp models.ObjectsPage
	return &resp, s.http.CallContext(ctx, &resp, getOwnedObjects, req.Address, req.Query, req.Cursor, req.Limit)
}

type QueryEventsRequest struct {
	Query           *models.EventFilter
	Cursor          *models.EventId // optional
	Limit           *uint           // optional
	DescendingOrder bool            // optional
}

func (s *ImplSuiAPI) QueryEvents(
	ctx context.Context,
	req *QueryEventsRequest,
) (*models.EventPage, error) {
	var resp models.EventPage
	return &resp, s.http.CallContext(ctx, &resp, queryEvents, req.Query, req.Cursor, req.Limit, req.DescendingOrder)
}

type QueryTransactionBlocksRequest struct {
	Query           *models.SuiTransactionBlockResponseQuery
	Cursor          *sui_types.TransactionDigest // optional
	Limit           *uint                        // optional
	DescendingOrder bool                         // optional
}

func (s *ImplSuiAPI) QueryTransactionBlocks(
	ctx context.Context,
	req *QueryTransactionBlocksRequest,
) (*models.TransactionBlocksPage, error) {
	resp := models.TransactionBlocksPage{}
	return &resp, s.http.CallContext(ctx, &resp, queryTransactionBlocks, req.Query, req.Cursor, req.Limit, req.DescendingOrder)
}

type ResolveNameServiceNamesRequest struct {
	Owner  *sui_types.SuiAddress
	Cursor *sui_types.ObjectID // optional
	Limit  *uint               // optional
}

func (s *ImplSuiAPI) ResolveNameServiceAddress(ctx context.Context, suiName string) (*sui_types.SuiAddress, error) {
	var resp sui_types.SuiAddress
	err := s.http.CallContext(ctx, &resp, resolveNameServiceAddress, suiName)
	if err != nil && err.Error() == "nil address" {
		return nil, errors.New("sui name not found")
	}
	return &resp, nil
}

func (s *ImplSuiAPI) ResolveNameServiceNames(
	ctx context.Context,
	req *ResolveNameServiceNamesRequest,
) (*models.SuiNamePage, error) {
	var resp models.SuiNamePage
	return &resp, s.http.CallContext(ctx, &resp, resolveNameServiceNames, req.Owner, req.Cursor, req.Limit)
}

func (s *ImplSuiAPI) SubscribeEvent(
	ctx context.Context,
	filter *models.EventFilter,
	resultCh chan models.SuiEvent,
) error {
	resp := make(chan []byte, 10)
	err := s.websocket.CallContext(ctx, resp, subscribeEvent, filter)
	if err != nil {
		return err
	}
	go func() {
		for messageData := range resp {
			var result models.SuiEvent
			if err := json.Unmarshal(messageData, &result); err != nil {
				log.Fatal(err)
			}

			resultCh <- result
		}

	}()
	return nil
}

func (s *ImplSuiAPI) SubscribeTransaction(
	ctx context.Context,
	filter *models.TransactionFilter,
	resultCh chan serialization.TagJson[models.SuiTransactionBlockEffects],
) error {
	resp := make(chan []byte, 10)
	err := s.websocket.CallContext(ctx, resp, subscribeTransaction, filter)
	if err != nil {
		return err
	}
	go func() {
		for messageData := range resp {
			var result serialization.TagJson[models.SuiTransactionBlockEffects]
			if err := json.Unmarshal(messageData, &result); err != nil {
				log.Fatal(err)
			}

			resultCh <- result
		}

	}()
	return nil
}
