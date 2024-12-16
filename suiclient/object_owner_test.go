package suiclient_test

import (
	"encoding/json"
	"testing"

	"github.com/pattonkan/sui-go/sui"
	"github.com/pattonkan/sui-go/suiclient"
	"github.com/stretchr/testify/require"
)

func TestObjectOwnerMarshal(t *testing.T) {
	{
		var dataStruct struct {
			Owner *suiclient.ObjectOwner `json:"owner"`
		}
		jsonString := []byte(`{"owner":"Immutable"}`)

		err := json.Unmarshal(jsonString, &dataStruct)
		require.NoError(t, err)
		enData, err := json.Marshal(dataStruct)
		require.NoError(t, err)
		require.Equal(t, jsonString, enData)
	}
	{
		var dataStruct struct {
			Owner *suiclient.ObjectOwner `json:"owner"`
		}
		jsonString := []byte(`{"owner":{"AddressOwner":"0xfb1f678fcfe31c7c1924319e49614ffbe3a984842ceed559aa2d772e60a2ef8f"}}`)

		err := json.Unmarshal(jsonString, &dataStruct)
		require.NoError(t, err)
		enData, err := json.Marshal(dataStruct)
		require.NoError(t, err)
		require.Equal(t, jsonString, enData)
	}
}

func TestIsSameStringAddress(t *testing.T) {
	type args struct {
		addr1 string
		addr2 string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "same address",
			args: args{
				"0x00000123",
				"0x000000123",
			},
			want: true,
		},
		{
			name: "not same address",
			args: args{
				"0x123f",
				"0x00000000123",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := sui.IsSameAddressString(tt.args.addr1, tt.args.addr2); got != tt.want {
				t.Errorf("IsSameStringAddress(): %v, want %v", got, tt.want)
			}
		})
	}
}
