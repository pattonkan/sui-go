package models

import "encoding/json"

type ProtocolConfig struct {
	MaxSupportedProtocolVersion SafeSuiBigInt[uint64]          `json:"maxSupportedProtocolVersion,omitempty"`
	MinSupportedProtocolVersion SafeSuiBigInt[uint64]          `json:"minSupportedProtocolVersion,omitempty"`
	ProtocolVersion             SafeSuiBigInt[uint64]          `json:"protocolVersion,omitempty"`
	Attributes                  map[string]ProtocolConfigValue `json:"attributes,omitempty"`
	FeatureFlags                map[string]bool                `json:"featureFlags,omitempty"`
}

type ProtocolConfigValue struct {
	U64 *string `json:"u64,omitempty"`
	U32 *string `json:"u32,omitempty"`
	F64 *string `json:"f64,omitempty"`
}

func (p *ProtocolConfigValue) UnmarshalJSON(data []byte) error {
	var temp map[string]string
	if err := json.Unmarshal(data, &temp); err != nil {
		return err
	}

	if u64, ok := temp["u64"]; ok {
		p.U64 = &u64
	} else if u32, ok := temp["u32"]; ok {
		p.U32 = &u32
	} else if f64, ok := temp["f64"]; ok {
		p.F64 = &f64
	}

	return nil
}
