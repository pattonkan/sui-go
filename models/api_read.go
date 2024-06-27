package models

import "github.com/howjmay/sui-go/sui_types"

type GetCheckpointsRequest struct {
	Cursor          *BigInt // optional
	Limit           *uint64 // optional
	DescendingOrder bool
}

type GetObjectRequest struct {
	ObjectID *sui_types.ObjectID
	Options  *SuiObjectDataOptions // optional
}

type GetTransactionBlockRequest struct {
	Digest  *sui_types.TransactionDigest
	Options *SuiTransactionBlockResponseOptions // optional
}

type MultiGetObjectsRequest struct {
	ObjectIDs []*sui_types.ObjectID
	Options   *SuiObjectDataOptions // optional
}

type MultiGetTransactionBlocksRequest struct {
	Digests []*sui_types.Digest
	Options *SuiTransactionBlockResponseOptions // optional
}

type TryGetPastObjectRequest struct {
	ObjectID *sui_types.ObjectID
	Version  uint64
	Options  *SuiObjectDataOptions // optional
}

type TryMultiGetPastObjectsRequest struct {
	PastObjects []*SuiGetPastObjectRequest
	Options     *SuiObjectDataOptions // optional
}
