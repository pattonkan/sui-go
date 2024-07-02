package models

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/howjmay/sui-go/sui_types"
	"github.com/howjmay/sui-go/sui_types/serialization"
)

type MoveFunctionArgType int

const (
	MoveFunctionArgTypeNone MoveFunctionArgType = iota
	MoveFunctionArgTypePure
	MoveFunctionArgTypeImmutableReference
	MoveFunctionArgTypeMutableReference
	MoveFunctionArgTypeValue
)

func (m *MoveFunctionArgType) UnmarshalJSON(data []byte) error {
	rawData := string(data)
	if strings.Contains(rawData, "Pure") {
		*m = MoveFunctionArgTypePure
	}
	if strings.Contains(rawData, "ByImmutableReference") {
		*m = MoveFunctionArgTypeImmutableReference
	}
	if strings.Contains(rawData, "ByMutableReference") {
		*m = MoveFunctionArgTypeMutableReference
	}
	if strings.Contains(rawData, "ByValue") {
		*m = MoveFunctionArgTypeValue
	}
	return nil
}

type MoveNormalizedFunction struct {
	Visibility     MoveVisibility
	IsEntry        bool
	TypeParameters []MoveAbilitySet
	Parameters     []MoveNormalizedType
	Return         []MoveNormalizedType
}

type MoveVisibility int

const (
	MoveVisibilityNone = iota
	MoveVisibilityPrivate
	MoveVisibilityPublic
	MoveVisibilityFriend
)

func (v *MoveVisibility) UnmarshalJSON(data []byte) error {
	switch string(data)[1 : len(data)-1] {
	case "Private":
		*v = MoveVisibilityPrivate
	case "Public":
		*v = MoveVisibilityPublic
	case "Friend":
		*v = MoveVisibilityFriend
	default:
		return fmt.Errorf("invalid json: %s", string(data))
	}
	return nil
}

type MoveAbilitySet struct {
	Abilities []MoveAbility
}

type MoveAbility int

const (
	MoveAbilityNone = iota
	MoveAbilityCopy
	MoveAbilityDrop
	MoveAbilityStore
	MoveAbilityKey
)

func (a *MoveAbility) UnmarshalJSON(data []byte) error {
	switch string(data)[1 : len(data)-1] {
	case "Copy":
		*a = MoveAbilityCopy
	case "Drop":
		*a = MoveAbilityDrop
	case "Store":
		*a = MoveAbilityStore
	case "Key":
		*a = MoveAbilityKey
	default:
		return fmt.Errorf("invalid json: %s", string(data))
	}
	return nil
}

type MoveNormalizedType struct {
	Bool             *serialization.EmptyEnum
	U8               *serialization.EmptyEnum
	U16              *serialization.EmptyEnum
	U32              *serialization.EmptyEnum
	U64              *serialization.EmptyEnum
	U128             *serialization.EmptyEnum
	U256             *serialization.EmptyEnum
	Address          *serialization.EmptyEnum
	Signer           *serialization.EmptyEnum
	Struct           *MoveNormalizedTypeStructType
	Vector           *MoveNormalizedType
	TypeParameter    *MoveTypeParameterIndex
	Reference        *MoveNormalizedType
	MutableReference *MoveNormalizedType
}

type MoveTypeParameterIndex uint16

func NewMoveTypeParameterIndex(i MoveTypeParameterIndex) *MoveTypeParameterIndex {
	return &i
}

type MoveNormalizedTypeStructType struct {
	Address       *sui_types.SuiAddress `json:"address"`
	Module        sui_types.Identifier  `json:"module"`
	Name          sui_types.Identifier  `json:"name"`
	TypeArguments []MoveNormalizedType  `json:"typeArguments"`
}

func (s *MoveNormalizedType) UnmarshalJSON(data []byte) error {
	var raw interface{}
	if err := json.Unmarshal(data, &raw); err != nil {
		return fmt.Errorf("failed to unmarshal %s to MoveNormalizedType: %w", string(data), err)
	}

	switch v := raw.(type) {
	case string:
		var emptyEnum serialization.EmptyEnum
		switch v {
		case "Bool":
			s.Bool = &emptyEnum
		case "U8":
			s.U8 = &emptyEnum
		case "U16":
			s.U16 = &emptyEnum
		case "U32":
			s.U32 = &emptyEnum
		case "U64":
			s.U64 = &emptyEnum
		case "U128":
			s.U128 = &emptyEnum
		case "U256":
			s.U256 = &emptyEnum
		case "Address":
			s.Address = &emptyEnum
		case "Signer":
			s.Signer = &emptyEnum
		default:
			return fmt.Errorf("unknown type: %s", v)
		}
	case map[string]interface{}:
		if rawVal, ok := v["Struct"]; ok {
			tmp, err := json.Marshal(rawVal)
			if err != nil {
				return fmt.Errorf("failed to decode rawVal: %v in Struct: %w", rawVal, err)
			}
			var val MoveNormalizedTypeStructType
			if err := json.Unmarshal(tmp, &val); err != nil {
				return fmt.Errorf("failed to unmarshal val: %v in Struct: %w", val, err)
			}
			s.Struct = &val
		} else if rawVal, ok := v["Vector"]; ok {
			tmp, err := json.Marshal(rawVal)
			if err != nil {
				return fmt.Errorf("failed to decode rawVal: %v in Vector: %w", rawVal, err)
			}
			var val MoveNormalizedType
			if err := json.Unmarshal(tmp, &val); err != nil {
				return fmt.Errorf("failed to unmarshal val: %v in Vector: %w", val, err)
			}
			s.Vector = &val
		} else if rawVal, ok := v["TypeParameter"]; ok {
			tmp := MoveTypeParameterIndex(rawVal.(float64))
			s.TypeParameter = &tmp
		} else if rawVal, ok := v["Reference"]; ok {
			tmp, err := json.Marshal(rawVal)
			if err != nil {
				return fmt.Errorf("failed to decode rawVal: %v in Reference: %w", rawVal, err)
			}
			var val MoveNormalizedType
			if err := json.Unmarshal(tmp, &val); err != nil {
				return fmt.Errorf("failed to unmarshal val: %v in Reference: %w", val, err)
			}
			s.Reference = &val
		} else if rawVal, ok := v["MutableReference"]; ok {
			tmp, err := json.Marshal(rawVal)
			if err != nil {
				return fmt.Errorf("failed to decode rawVal: %v in MutableReference: %w", rawVal, err)
			}
			var val MoveNormalizedType
			if err := json.Unmarshal(tmp, &val); err != nil {
				return fmt.Errorf("failed to unmarshal val: %v in MutableReference: %w", val, err)
			}
			s.MutableReference = &val
		} else {
			return fmt.Errorf("unknown type object: %v", v)
		}
	default:
		return fmt.Errorf("unexpected type: %T", v)
	}
	return nil
}

type MoveNormalizedModule struct {
	FileFormatVersion uint32
	Address           *sui_types.SuiAddress
	Name              sui_types.Identifier
	Friends           []*MoveModuleID
	Structs           map[sui_types.Identifier]*MoveNormalizedStruct
	ExposedFunctions  map[sui_types.Identifier]*MoveNormalizedFunction
}

type MoveModuleID struct {
	Address *sui_types.SuiAddress
	Name    sui_types.Identifier
}

type MoveNormalizedStruct struct {
	Abilities      MoveAbilitySet
	TypeParameters []*MoveStructTypeParameter
	Fields         []*MoveNormalizedField
}

type MoveStructTypeParameter struct {
	Constraints MoveAbilitySet
	IsPhantom   bool
}

type MoveNormalizedField struct {
	Name sui_types.Identifier
	Type *MoveNormalizedType
}
