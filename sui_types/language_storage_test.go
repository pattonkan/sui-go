package sui_types_test

import (
	"testing"

	"github.com/fardream/go-bcs/bcs"
	"github.com/howjmay/sui-go/sui_types"
	"github.com/howjmay/sui-go/sui_types/serialization"
	"github.com/stretchr/testify/require"
)

func TestTypeTagEncoding(t *testing.T) {
	typeTagU64 := sui_types.TypeTag{U64: &serialization.EmptyEnum{}}
	typeTagMarshaled, err := bcs.Marshal(typeTagU64)
	require.NoError(t, err)
	require.Equal(t, []byte{2}, typeTagMarshaled)

	typeTagStruct := sui_types.TypeTag{Struct: &sui_types.StructTag{
		Address: *sui_types.MustObjectIDFromHex("0x2eeb551107032ae860d76661f3f4573dd0f8c701116137e6525dcd95d4f8e58"),
		Module:  "testcoin",
		Name:    "TESTCOIN",
	}}
	typeTagStructMarshaled, err := bcs.Marshal(typeTagStruct)
	require.NoError(t, err)
	var structTag sui_types.TypeTag
	_, err = bcs.Unmarshal(typeTagStructMarshaled, &structTag)
	require.NoError(t, err)
}
