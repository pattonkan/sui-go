package models

import (
	"github.com/howjmay/sui-go/sui_types"
)

type ObjectOwnerInternal struct {
	AddressOwner *sui_types.SuiAddress `json:"AddressOwner,omitempty"`
	ObjectOwner  *sui_types.SuiAddress `json:"ObjectOwner,omitempty"`
	SingleOwner  *sui_types.SuiAddress `json:"SingleOwner,omitempty"`
	Shared       *struct {
		InitialSharedVersion *sui_types.SequenceNumber `json:"initial_shared_version"`
	} `json:"Shared,omitempty"`
}

type ObjectOwner struct {
	*ObjectOwnerInternal
	*string
}

type Page[T SuiTransactionBlockResponse | SuiEvent | Coin | *Coin | SuiObjectResponse | DynamicFieldInfo | string | *Checkpoint,
	C sui_types.TransactionDigest | EventId | sui_types.ObjectID | BigInt] struct {
	Data []T `json:"data"`
	// 'NextCursor' points to the last item in the page.
	// Reading with next_cursor will start from the next item after next_cursor
	NextCursor  *C   `json:"nextCursor,omitempty"`
	HasNextPage bool `json:"hasNextPage"`
}
