package movebcs

import (
	"github.com/fardream/go-bcs/bcs"
	"github.com/howjmay/sui-go/sui"
)

const (
	// The `Borrow` does not match the `Referent`.
	EBorrowWrongBorrow = 0
	// An attempt to swap the `Referent.value` with another object of the same type.
	EBorrowWrongValue = 1
)

// An object wrapping a `T` and providing the borrow API.
type Referent[T any] struct {
	Id    *sui.Address
	Value bcs.Option[T]
}

// A hot potato making sure the object is put back once borrowed.
type Borrow struct {
	Ref *sui.Address
	Obj *sui.ObjectId
}
