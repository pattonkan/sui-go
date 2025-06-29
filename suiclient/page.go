package suiclient

import (
	"fmt"

	"github.com/fardream/go-bcs/bcs"
	"github.com/pattonkan/sui-go/sui"
)

type Page[T SuiTransactionBlockResponse | Event | Coin | *Coin | SuiObjectResponse | DynamicFieldInfo | string | *Checkpoint,
	C sui.TransactionDigest | EventId | sui.ObjectId | sui.BigInt | string] struct {
	Data []T `json:"data"`
	// 'NextCursor' points to the last item in the page.
	// Reading with next_cursor will start from the next item after next_cursor
	NextCursor  *C   `json:"nextCursor,omitempty"`
	HasNextPage bool `json:"hasNextPage"`
}

// Sui has 'BalanceCursor' and 'ObjectCursor' in the wrapper type 'BcsCursor'
// Here we implement 'Cursor' as BalanceCursor, which has only one more field 'CoinBalanceBucket' than ObjectCursor.
type Cursor struct {
	ObjectId          *sui.ObjectId
	CpSequenceNumber  uint64
	CoinBalanceBucket *uint64 `bcs:"optional"` //
}

func (c *Cursor) Bytes() []byte {
	if c == nil {
		return nil
	}
	b, err := bcs.Marshal(c)
	if err != nil {
		return nil
	}
	return b
}

func (c *Cursor) Base64() sui.Base64 {
	return sui.Base64(c.Bytes())
}

func DecodeCursor(rawCursor string) (*Cursor, error) {
	var cursor Cursor
	b64, err := sui.NewBase64(rawCursor)
	if err != nil {
		return nil, fmt.Errorf("failed to decode page's cursor from base64: %w", err)
	}
	_, err = bcs.Unmarshal(b64.Data(), &cursor)
	if err != nil {
		return nil, fmt.Errorf("failed to decode page's cursor from BCS: %w", err)
	}
	return &cursor, nil

}
