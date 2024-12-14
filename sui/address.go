package sui

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"strings"
)

const AddressLen = 32

type Address [AddressLen]uint8

func AddressFromHex(str string) (*Address, error) {
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
	if len(data) > AddressLen {
		return nil, errors.New("the len is invalid")
	}
	var address Address
	copy(address[AddressLen-len(data):], data[:])
	return &address, nil
}

func MustAddressFromHex(str string) *Address {
	addr, err := AddressFromHex(str)
	if err != nil {
		panic(err)
	}
	return addr
}

func (a Address) Bytes() []byte {
	return a[:]
}
func (a Address) Length() int {
	return len(a)
}
func (a Address) String() string {
	return "0x" + hex.EncodeToString(a[:])
}

// ShortString Returns the address with leading zeros trimmed, e.g. 0x2
func (a Address) ShortString() string {
	return "0x" + strings.TrimLeft(hex.EncodeToString(a[:]), "0")
}

func (a Address) MarshalJSON() ([]byte, error) {
	return json.Marshal(a.String())
}

func (a *Address) UnmarshalJSON(data []byte) error {
	var str *string
	err := json.Unmarshal(data, &str)
	if err != nil {
		return err
	}
	if str == nil {
		return errors.New("nil address")
	}
	tmp, err := AddressFromHex(*str)
	if err == nil {
		*a = *tmp
	}
	return err
}

func (a Address) MarshalBCS() ([]byte, error) {
	return a[:], nil
}
