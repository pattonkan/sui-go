package sui

import (
	"encoding/json"
	"fmt"
	"math/big"
	"strings"
)

type BigInt struct {
	*big.Int
}

func NewBigInt(v uint64) *BigInt {
	return &BigInt{new(big.Int).SetUint64(v)}
}

func NewBigIntInt64(v int64) *BigInt {
	return &BigInt{new(big.Int).SetInt64(v)}
}

func (b *BigInt) UnmarshalText(data []byte) error {
	return b.UnmarshalJSON(data)
}

func (b *BigInt) UnmarshalJSON(data []byte) error {
	rawData := strings.TrimSpace(string(data))
	if strings.HasPrefix(rawData, `"`) && strings.HasSuffix(rawData, `"`) {
		rawData = rawData[1 : len(rawData)-1]
	}
	if b.Int == nil {
		b.Int = new(big.Int)
	}
	if rawData == "null" {
		b.SetInt64(0)
		return nil
	}
	_, ok := b.SetString(rawData, 10)
	if ok {
		return nil
	}
	return fmt.Errorf("json data [%s] is not T", string(data))
}

func (b *BigInt) MarshalJSON() ([]byte, error) {
	return json.Marshal(b.String())
}

func (b *BigInt) BigInt() *big.Int {
	return b.Int
}
