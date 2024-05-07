package sui_types

import "github.com/howjmay/sui-go/lib"

type StructTag struct {
	Address    SuiAddress
	Module     string
	Name       string
	TypeParams []TypeTag
}

type TypeTag struct {
	Bool    *lib.EmptyEnum
	U8      *lib.EmptyEnum
	U16     *lib.EmptyEnum
	U32     *lib.EmptyEnum
	U64     *lib.EmptyEnum
	U128    *lib.EmptyEnum
	U256    *lib.EmptyEnum
	Address *lib.EmptyEnum
	Signer  *lib.EmptyEnum
	Vector  *TypeTag
	Struct  *StructTag
}

func (t TypeTag) IsBcsEnum() {}
