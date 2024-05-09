package sui_types

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"strings"
)

const SuiAddressLen = 32

type SuiAddress [SuiAddressLen]uint8

func SuiAddressFromHex(str string) (*SuiAddress, error) {
	if strings.HasPrefix(str, "0x") || strings.HasPrefix(str, "0X") {
		str = str[2:]
	}
	if len(str)%2 != 0 {
		str = "0" + str
	}
	data, err := hex.DecodeString(str)
	if err != nil {
		return nil, err
	}
	if len(data) > SuiAddressLen {
		return nil, errors.New("the len is invalid")
	}
	var address SuiAddress
	copy(address[SuiAddressLen-len(data):], data[:])
	return &address, nil
}

func MustSuiAddressFromHex(str string) *SuiAddress {
	addr, err := SuiAddressFromHex(str)
	if err != nil {
		panic(err)
	}
	return addr
}

func (a SuiAddress) Data() []byte {
	return a[:]
}
func (a SuiAddress) Length() int {
	return len(a)
}
func (a SuiAddress) String() string {
	return "0x" + hex.EncodeToString(a[:])
}

func (a SuiAddress) ShortString() string {
	return "0x" + strings.TrimLeft(hex.EncodeToString(a[:]), "0")
}

func (a SuiAddress) MarshalJSON() ([]byte, error) {
	return json.Marshal(a.String())
}

func (a *SuiAddress) UnmarshalJSON(data []byte) error {
	var str *string
	err := json.Unmarshal(data, &str)
	if err != nil {
		return err
	}
	if str == nil {
		return errors.New("nil address")
	}
	tmp, err := SuiAddressFromHex(*str)
	if err == nil {
		*a = *tmp
	}
	return err
}

// FIXME may need to be pointer
func (a SuiAddress) MarshalBCS() ([]byte, error) {
	return a[:], nil
}
