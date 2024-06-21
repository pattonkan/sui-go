package sui_types_test

import (
	"testing"

	"github.com/fardream/go-bcs/bcs"
	"github.com/howjmay/sui-go/sui_types"
	"github.com/howjmay/sui-go/sui_types/serialization"
	"github.com/stretchr/testify/require"
)

func TestTypeTagBcs(t *testing.T) {
	typeTagU64 := sui_types.TypeTag{U64: &serialization.EmptyEnum{}}
	typeTagMarshaled, err := bcs.Marshal(typeTagU64)
	require.NoError(t, err)
	require.Equal(t, []byte{2}, typeTagMarshaled)

	typeTagStruct := sui_types.TypeTag{Struct: &sui_types.StructTag{
		Address: sui_types.MustObjectIDFromHex("0x2eeb551107032ae860d76661f3f4573dd0f8c701116137e6525dcd95d4f8e58"),
		Module:  "testcoin",
		Name:    "TESTCOIN",
	}}
	typeTagStructMarshaled, err := bcs.Marshal(typeTagStruct)
	require.NoError(t, err)
	var structTag sui_types.TypeTag
	_, err = bcs.Unmarshal(typeTagStructMarshaled, &structTag)
	require.NoError(t, err)
}

func TestStructTagEncoding(t *testing.T) {
	s1 := `0x2::foo::bar<0x3::baz::qux<0x4::nested::result, 0x5::funny::other>, bool>`
	structTag, err := sui_types.StructTagFromString(s1)
	require.NoError(t, err)

	require.Equal(t, sui_types.MustObjectIDFromHex("0x2"), structTag.Address)
	require.Equal(t, sui_types.Identifier("foo"), structTag.Module)
	require.Equal(t, sui_types.Identifier("bar"), structTag.Name)

	typeParam0 := structTag.TypeParams[0].Struct
	require.Equal(t, sui_types.MustObjectIDFromHex("0x3"), typeParam0.Address)
	require.Equal(t, sui_types.Identifier("baz"), typeParam0.Module)
	require.Equal(t, sui_types.Identifier("qux"), typeParam0.Name)
	typeParam00 := structTag.TypeParams[0].Struct.TypeParams[0].Struct
	require.Equal(t, sui_types.MustObjectIDFromHex("0x4"), typeParam00.Address)
	require.Equal(t, sui_types.Identifier("nested"), typeParam00.Module)
	require.Equal(t, sui_types.Identifier("result"), typeParam00.Name)

	typeParam01 := structTag.TypeParams[0].Struct.TypeParams[1].Struct
	require.Equal(t, sui_types.MustObjectIDFromHex("0x5"), typeParam01.Address)
	require.Equal(t, sui_types.Identifier("funny"), typeParam01.Module)
	require.Equal(t, sui_types.Identifier("other"), typeParam01.Name)

	require.NotNil(t, structTag.TypeParams[1].Bool)
}
