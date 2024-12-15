package sui

import (
	"encoding/json"
	"fmt"
	"strings"
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
	Bool             *EmptyEnum
	U8               *EmptyEnum
	U16              *EmptyEnum
	U32              *EmptyEnum
	U64              *EmptyEnum
	U128             *EmptyEnum
	U256             *EmptyEnum
	Address          *EmptyEnum
	Signer           *EmptyEnum
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
	Address       *Address             `json:"address"`
	Module        Identifier           `json:"module"`
	Name          Identifier           `json:"name"`
	TypeArguments []MoveNormalizedType `json:"typeArguments"`
}

func (s *MoveNormalizedType) UnmarshalJSON(data []byte) error {
	var raw interface{}
	if err := json.Unmarshal(data, &raw); err != nil {
		return fmt.Errorf("failed to unmarshal %s to MoveNormalizedType: %w", string(data), err)
	}

	switch v := raw.(type) {
	case string:
		var emptyEnum EmptyEnum
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
	Address           *Address
	Name              Identifier
	Friends           []*MoveModuleId
	Structs           map[Identifier]*MoveNormalizedStruct
	ExposedFunctions  map[Identifier]*MoveNormalizedFunction
}

type MoveModule struct {
	Package *ObjectId  `json:"package"`
	Module  Identifier `json:"module"`
}

type MoveModuleId struct {
	Address *Address
	Name    Identifier
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
	Name Identifier
	Type *MoveNormalizedType
}
