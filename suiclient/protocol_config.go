package suiclient

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/pattonkan/sui-go/sui"
)

type ProtocolConfig struct {
	MaxSupportedProtocolVersion *sui.BigInt                    `json:"maxSupportedProtocolVersion,omitempty"`
	MinSupportedProtocolVersion *sui.BigInt                    `json:"minSupportedProtocolVersion,omitempty"`
	ProtocolVersion             *sui.BigInt                    `json:"protocolVersion,omitempty"`
	Attributes                  map[string]ProtocolConfigValue `json:"attributes,omitempty"`
	FeatureFlags                map[string]bool                `json:"featureFlags,omitempty"`
}

type ProtocolConfigValue struct {
	U16 *uint16  `json:"u16,omitempty"`
	U32 *uint32  `json:"u32,omitempty"`
	U64 *uint64  `json:"u64,omitempty"`
	F64 *float64 `json:"f64,omitempty"`
}

func (p *ProtocolConfigValue) UnmarshalJSON(data []byte) error {
	var temp map[string]string
	if err := json.Unmarshal(data, &temp); err != nil {
		return err
	}

	if u16str, ok := temp["u16"]; ok {
		_u16, err := strconv.ParseUint(u16str, 10, 16)
		if err != nil {
			return fmt.Errorf("can't parse %s to u16", u16str)
		}
		u16 := uint16(_u16)
		p.U16 = &u16
	} else if u32str, ok := temp["u32"]; ok {
		_u32, err := strconv.ParseUint(u32str, 10, 32)
		if err != nil {
			return fmt.Errorf("can't parse %s to u32", u32str)
		}
		u32 := uint32(_u32)
		p.U32 = &u32
	} else if u64str, ok := temp["u64"]; ok {
		u64, err := strconv.ParseUint(u64str, 10, 64)
		if err != nil {
			return fmt.Errorf("can't parse %s to u64", u64str)
		}
		p.U64 = &u64
	} else if f64str, ok := temp["f64"]; ok {
		f64, err := strconv.ParseFloat(f64str, 64)
		if err != nil {
			return fmt.Errorf("can't parse %s to f64", f64str)
		}
		p.F64 = &f64
	}

	return nil
}
