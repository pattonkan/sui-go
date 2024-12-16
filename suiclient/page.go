package suiclient

import "github.com/pattonkan/sui-go/sui"

type Page[T SuiTransactionBlockResponse | Event | Coin | *Coin | SuiObjectResponse | DynamicFieldInfo | string | *Checkpoint,
	C sui.TransactionDigest | EventId | sui.ObjectId | sui.BigInt] struct {
	Data []T `json:"data"`
	// 'NextCursor' points to the last item in the page.
	// Reading with next_cursor will start from the next item after next_cursor
	NextCursor  *C   `json:"nextCursor,omitempty"`
	HasNextPage bool `json:"hasNextPage"`
}
