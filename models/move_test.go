package models_test

import (
	"reflect"
	"testing"

	"github.com/howjmay/sui-go/models"
	"github.com/howjmay/sui-go/sui_types"
	"github.com/stretchr/testify/require"
)

func TestNewResourceType(t *testing.T) {
	tests := []struct {
		name    string
		str     string
		want    *models.ResourceType
		wantErr bool
	}{
		{
			name: "sample",
			str:  "0x23::coin::Xxxx",
			want: &models.ResourceType{addressFromHex(t, "0x23"), "coin", "Xxxx", nil},
		},
		{
			name: "three level",
			str:  "0xabc::Coin::Xxxx<0x789::AAA::ppp<0x111::mod3::func3>>",
			want: &models.ResourceType{
				addressFromHex(t, "0xabc"), "Coin", "Xxxx",
				&models.ResourceType{
					addressFromHex(t, "0x789"), "AAA", "ppp",
					&models.ResourceType{addressFromHex(t, "0x111"), "mod3", "func3", nil},
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
				got, err := models.NewResourceType(tt.str)
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

func TestResourceType_String(t *testing.T) {
	typeString := "0x1::mmm1::fff1<0x123abcdef::mm2::ff3>"

	resourceType, err := models.NewResourceType(typeString)
	require.NoError(t, err)
	res := "0x0000000000000000000000000000000000000000000000000000000000000001::mmm1::fff1<0x0000000000000000000000000000000000000000000000000000000123abcdef::mm2::ff3>"
	require.Equal(t, resourceType.String(), res)
}

func TestResourceType_ShortString(t *testing.T) {
	tests := []struct {
		name string
		arg  *models.ResourceType
		want string
	}{
		{
			arg:  &models.ResourceType{addressFromHex(t, "0x1"), "m1", "f1", nil},
			want: "0x1::m1::f1",
		},
		{
			arg: &models.ResourceType{
				addressFromHex(t, "0x1"), "m1", "f1",
				&models.ResourceType{
					addressFromHex(t, "2"), "m2", "f2",
					&models.ResourceType{addressFromHex(t, "0x123abcdef"), "m3", "f3", nil},
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

func addressFromHex(t *testing.T, hex string) *sui_types.SuiAddress {
	addr, err := sui_types.AddressFromHex(hex)
	require.NoError(t, err)
	return addr
}
