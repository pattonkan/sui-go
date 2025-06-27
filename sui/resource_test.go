package sui_test

import (
	"reflect"
	"testing"

	"github.com/pattonkan/sui-go/sui"
	"github.com/stretchr/testify/require"
)

func TestNewResourceType(t *testing.T) {
	tests := []struct {
		name    string
		str     string
		want    *sui.ResourceType
		wantErr bool
	}{
		{
			name: "sample",
			str:  "0x23::coin::Xxxx",
			want: &sui.ResourceType{sui.MustAddressFromHex("0x23"), "coin", "Xxxx", nil},
		},
		{
			name: "three level",
			str:  "0xabc::Coin::Xxxx<0x789::AAA::ppp<0x111::mod3::func3>>",
			want: &sui.ResourceType{
				sui.MustAddressFromHex("0xabc"), "Coin", "Xxxx",
				&sui.ResourceType{
					sui.MustAddressFromHex("0x789"), "AAA", "ppp",
					&sui.ResourceType{sui.MustAddressFromHex("0x111"), "mod3", "func3", nil},
				},
			},
		},
		{
			name:    "error address",
			str:     "0x123abcg::coin::Xxxx",
			wantErr: true,
		},
		{
			name:    "error format",
			str:     "0x1::m1::f1<0x2::m2::f2>x",
			wantErr: true,
		},
		{
			name:    "error format2",
			str:     "0x1::m1::f1<<0x3::m3::f3>0x2::m2::f2>",
			wantErr: true,
		},
		{
			name:    "error format3",
			str:     "<0x3::m3::f3>0x1::m1::f1<0x2::m2::f2>",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				got, err := sui.NewResourceType(tt.str)
				if (err != nil) != tt.wantErr {
					t.Errorf("NewResourceType() error: %v, wantErr %v", err, tt.wantErr)
					return
				}
				if !reflect.DeepEqual(got, tt.want) {
					t.Errorf("NewResourceType(): %v, want %v", got, tt.want)
				}
			},
		)
	}
}

func TestContains(t *testing.T) {
	tests := []struct {
		name   string
		str    string
		target *sui.ResourceType
		want   bool
	}{
		{
			name:   "successful, two levels",
			str:    "0xe87e::swap::Pool<0x2f63::testcoin::TESTCOIN>",
			target: &sui.ResourceType{Module: "swap", ObjectName: "Pool"},
			want:   true,
		},
		{
			name:   "unsuccessful, two levels",
			str:    "0xe87e::swap::Pool<0x2f63::testcoin::TESTCOIN>",
			target: &sui.ResourceType{Module: "name", ObjectName: "Pool"},
			want:   false,
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				src, err := sui.NewResourceType(tt.str)
				require.NoError(t, err)
				require.Equal(t, tt.want, src.Contains(tt.target.Address, tt.target.Module, tt.target.ObjectName))
			},
		)
	}
}

func TestResourceTypeString(t *testing.T) {
	typeString := "0x1::mmm1::fff1<0x123abcdef::mm2::ff3>"

	resourceType, err := sui.NewResourceType(typeString)
	require.NoError(t, err)
	res := "0x0000000000000000000000000000000000000000000000000000000000000001::mmm1::fff1<0x0000000000000000000000000000000000000000000000000000000123abcdef::mm2::ff3>"
	require.Equal(t, resourceType.String(), res)
}

func TestResourceTypeShortString(t *testing.T) {
	tests := []struct {
		name string
		arg  *sui.ResourceType
		want string
	}{
		{
			arg:  &sui.ResourceType{sui.MustAddressFromHex("0x1"), "m1", "f1", nil},
			want: "0x1::m1::f1",
		},
		{
			arg: &sui.ResourceType{
				sui.MustAddressFromHex("0x1"), "m1", "f1",
				&sui.ResourceType{
					sui.MustAddressFromHex("2"), "m2", "f2",
					&sui.ResourceType{sui.MustAddressFromHex("0x123abcdef"), "m3", "f3", nil},
				},
			},
			want: "0x1::m1::f1<0x2::m2::f2<0x123abcdef::m3::f3>>",
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				if got := tt.arg.ShortString(); got != tt.want {
					t.Errorf("ResourceType.ShortString(): %v, want %v", got, tt.want)
				}
			},
		)
	}
}
