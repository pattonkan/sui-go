package models

import (
	"math/big"

	"github.com/howjmay/sui-go/sui_types"
)

type BigInt = big.Int

// type BigInt struct{ big.Int }

// type *BigInt = decimal.Decimal

// type SafeBigInt interface {
// 	~int64 | ~uint64
// }

// func New*BigInt[T SafeBigInt](num T) *BigInt[T] {
// 	return *BigInt[T]{
// 		data: num,
// 	}
// }

// type *BigInt[T SafeBigInt] struct {
// 	data T
// }

// func (s *models.BigInt[T]) UnmarshalText(data []byte) error {
// 	return s.UnmarshalJSON(data)
// }

// func (s *models.BigInt[T]) UnmarshalJSON(data []byte) error {
// 	num := decimal.NewFromInt(0)
// 	err := num.UnmarshalJSON(data)
// 	if err != nil {
// 		return err
// 	}

// 	if num.BigInt().IsInt64() {
// 		s.data = T(num.BigInt().Int64())
// 		return nil
// 	}

// 	if num.BigInt().IsUint64() {
// 		s.data = T(num.BigInt().Uint64())
// 		return nil
// 	}
// 	return fmt.Errorf("json data [%s] is not T", string(data))
// }

// func (s *BigInt[T]) MarshalJSON() ([]byte, error) {
// 	return decimal.NewFromInt(int64(s.data)).MarshalJSON()
// }

// func (s *BigInt[T]) Int64() int64 {
// 	return int64(s.data)
// }

// func (s *BigInt[T]) Uint64() uint64 {
// 	return uint64(s.data)
// }

// func (s *models.BigInt[T]) Decimal() decimal.Decimal {
// 	return decimal.NewFromBigInt(new(big.Int).SetUint64(s.Uint64()), 0)
// }

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

type Page[T SuiTransactionBlockResponse | SuiEvent | Coin | *Coin | SuiObjectResponse | DynamicFieldInfo | string,
	C sui_types.TransactionDigest | EventId | sui_types.ObjectID] struct {
	Data []T `json:"data"`
	// 'NextCursor' points to the last item in the page.
	// Reading with next_cursor will start from the next item after next_cursor
	NextCursor  *C   `json:"nextCursor,omitempty"`
	HasNextPage bool `json:"hasNextPage"`
}
