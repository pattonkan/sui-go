package models

import "github.com/howjmay/sui-go/sui_types"

type GetDynamicFieldObjectRequest struct {
	ParentObjectID *sui_types.ObjectID
	Name           *sui_types.DynamicFieldName
}

type GetDynamicFieldsRequest struct {
	ParentObjectID *sui_types.ObjectID
	Cursor         *sui_types.ObjectID // optional
	Limit          *uint               // optional
}

type GetOwnedObjectsRequest struct {
	Address *sui_types.SuiAddress
	Query   *SuiObjectResponseQuery // optional
	Cursor  *CheckpointedObjectID   // optional
	Limit   *uint                   // optional
}

type QueryEventsRequest struct {
	Query           *EventFilter
	Cursor          *EventId // optional
	Limit           *uint    // optional
	DescendingOrder bool     // optional
}

type QueryTransactionBlocksRequest struct {
	Query           *SuiTransactionBlockResponseQuery
	Cursor          *sui_types.TransactionDigest // optional
	Limit           *uint                        // optional
	DescendingOrder bool                         // optional
}

type ResolveNameServiceNamesRequest struct {
	Owner  *sui_types.SuiAddress
	Cursor *sui_types.ObjectID // optional
	Limit  *uint               // optional
}
