package sui_test

import (
	"testing"

	"github.com/howjmay/sui-go/sui"
	"github.com/howjmay/sui-go/sui_types"
	"github.com/stretchr/testify/require"
)

func TestToSuiJsonArg(t *testing.T) {
	tests := []struct {
		name string
		arg  any
		want any
	}{
		{
			name: "uint8",
			arg:  uint8(3),
			want: "0x03",
		},
		{
			name: "uint16",
			arg:  uint16(3),
			want: "0x0003",
		},
		{
			name: "uint32",
			arg:  uint32(3),
			want: "0x00000003",
		},
		{
			name: "uint64",
			arg:  uint64(3),
			want: "0x0000000000000003",
		},
		{
			name: "SuiAddress",
			arg:  sui_types.MustSuiAddressFromHex("04a0b8440240ebd7db9a5158332f496e4a1b18576bbe5a925879eacadc1b61ba"),
			want: "0x04a0b8440240ebd7db9a5158332f496e4a1b18576bbe5a925879eacadc1b61ba",
		},
		{
			name: "ObjectID",
			arg:  sui_types.MustObjectIDFromHex("04a0b8440240ebd7db9a5158332f496e4a1b18576bbe5a925879eacadc1b61ba"),
			want: "0x04a0b8440240ebd7db9a5158332f496e4a1b18576bbe5a925879eacadc1b61ba",
		},
		{
			name: "PackageID",
			arg:  sui_types.MustPackageIDFromHex("04a0b8440240ebd7db9a5158332f496e4a1b18576bbe5a925879eacadc1b61ba"),
			want: "0x04a0b8440240ebd7db9a5158332f496e4a1b18576bbe5a925879eacadc1b61ba",
		},
		{
			name: "byte array",
			arg:  []byte("abcd"),
			want: []byte("abcd"),
		},
		{
			name: "array of byte arrays",
			arg:  [][]byte{[]byte("abc"), []byte("1234")},
			want: [][]byte{[]byte("abc"), []byte("1234")},
		},
		{
			name: "address array",
			arg:  []*sui_types.SuiAddress{sui_types.MustSuiAddressFromHex("0x3a14a003e3e890906d10eb5069bc3048ad4c3a7cbd1733b1ff21cd1c7fbf032e"), sui_types.MustSuiAddressFromHex("0x0c19db98f4b5562b4832a0c1fb6f77a0e9f1ed971e74c696ee21cc08ddfa2e91"), sui_types.MustSuiAddressFromHex("0xf69443d142f99550481984e3402534b4f3046e33892fdbd46be190c9e2c7cca8")},
			want: []*sui_types.SuiAddress{sui_types.MustSuiAddressFromHex("0x3a14a003e3e890906d10eb5069bc3048ad4c3a7cbd1733b1ff21cd1c7fbf032e"), sui_types.MustSuiAddressFromHex("0x0c19db98f4b5562b4832a0c1fb6f77a0e9f1ed971e74c696ee21cc08ddfa2e91"), sui_types.MustSuiAddressFromHex("0xf69443d142f99550481984e3402534b4f3046e33892fdbd46be190c9e2c7cca8")},
		},
		{
			name: "string array",
			arg:  []string{"abc", "1234"},
			want: []string{"abc", "1234"},
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				got := sui.ToSuiJsonArg(tt.arg)
				require.Equal(t, tt.want, got)
			},
		)
	}
}
