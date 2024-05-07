package sui_types

import "github.com/howjmay/sui-go/sui_types/serialization"

type StructTag struct {
	Address    SuiAddress
	Module     string
	Name       string
	TypeParams []TypeTag
}

type TypeTag struct {
	Bool    *serialization.EmptyEnum
	U8      *serialization.EmptyEnum
	U16     *serialization.EmptyEnum
	U32     *serialization.EmptyEnum
	U64     *serialization.EmptyEnum
	U128    *serialization.EmptyEnum
	U256    *serialization.EmptyEnum
	Address *serialization.EmptyEnum
	Signer  *serialization.EmptyEnum
	Vector  *TypeTag
	Struct  *StructTag
}

func (t TypeTag) IsBcsEnum() {}
