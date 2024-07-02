package sui

import (
	"context"
	"encoding/json"
	"errors"
	"log"

	"github.com/howjmay/sui-go/models"
	"github.com/howjmay/sui-go/sui_types"
	"github.com/howjmay/sui-go/sui_types/serialization"
	"github.com/tidwall/gjson"
)

func (s *ImplSuiAPI) GetDynamicFieldObject(
	ctx context.Context,
	req *models.GetDynamicFieldObjectRequest,
) (*models.SuiObjectResponse, error) {
	var resp models.SuiObjectResponse
	return &resp, s.http.CallContext(ctx, &resp, getDynamicFieldObject, req.ParentObjectID, req.Name)
}

func (s *ImplSuiAPI) GetDynamicFields(
	ctx context.Context,
	req *models.GetDynamicFieldsRequest,
) (*models.DynamicFieldPage, error) {
	var resp models.DynamicFieldPage
	return &resp, s.http.CallContext(ctx, &resp, getDynamicFields, req.ParentObjectID, req.Cursor, req.Limit)
}

// address : <SuiAddress> - the owner's Sui address
// query : <ObjectResponseQuery> - the objects query criteria.
// cursor : <CheckpointedObjectID> - An optional paging cursor. If provided, the query will start from the next item after the specified cursor. Default to start from the first item if not specified.
// limit : <uint> - Max number of items returned per page, default to [QUERY_MAX_RESULT_LIMIT_OBJECTS] if is 0
func (s *ImplSuiAPI) GetOwnedObjects(
	ctx context.Context,
	req *models.GetOwnedObjectsRequest,
) (*models.ObjectsPage, error) {
	var resp models.ObjectsPage
	return &resp, s.http.CallContext(ctx, &resp, getOwnedObjects, req.Address, req.Query, req.Cursor, req.Limit)
}

func (s *ImplSuiAPI) QueryEvents(
	ctx context.Context,
	req *models.QueryEventsRequest,
) (*models.EventPage, error) {
	var resp models.EventPage
	return &resp, s.http.CallContext(ctx, &resp, queryEvents, req.Query, req.Cursor, req.Limit, req.DescendingOrder)
}

func (s *ImplSuiAPI) QueryTransactionBlocks(
	ctx context.Context,
	req *models.QueryTransactionBlocksRequest,
) (*models.TransactionBlocksPage, error) {
	resp := models.TransactionBlocksPage{}
	return &resp, s.http.CallContext(ctx, &resp, queryTransactionBlocks, req.Query, req.Cursor, req.Limit, req.DescendingOrder)
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
	req *models.ResolveNameServiceNamesRequest,
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
