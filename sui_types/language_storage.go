package sui_types

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/howjmay/sui-go/sui_types/serialization"
)

// refer BCS doc https://github.com/diem/bcs/blob/master/README.md#externally-tagged-enumerations
// IMPORTANT! The order of the fields MATTERS! DON'T CHANGE!
// this is enum `TypeTag` in `external-crates/move/crates/move-core-types/src/language_storage.rs`
// each field should be the same as enum `TypeTag` there
type TypeTag struct {
	Bool    *serialization.EmptyEnum
	U8      *serialization.EmptyEnum
	U64     *serialization.EmptyEnum
	U128    *serialization.EmptyEnum
	Address *serialization.EmptyEnum
	Signer  *serialization.EmptyEnum
	Vector  *TypeTag
	Struct  *StructTag

	U16  *serialization.EmptyEnum
	U32  *serialization.EmptyEnum
	U256 *serialization.EmptyEnum
}

func (t TypeTag) IsBcsEnum() {}

// refer TypeTagSerializer.parseFromStr() at 'sdk/typescript/src/bcs/type-tag-serializer.ts'
func TypeTagFromString(data string) (*TypeTag, error) {
	switch data {
	case "address":
		return &TypeTag{Address: &serialization.EmptyEnum{}}, nil
	case "signer":
		return &TypeTag{Signer: &serialization.EmptyEnum{}}, nil
	case "bool":
		return &TypeTag{Bool: &serialization.EmptyEnum{}}, nil
	case "u8":
		return &TypeTag{U8: &serialization.EmptyEnum{}}, nil
	case "u16":
		return &TypeTag{U16: &serialization.EmptyEnum{}}, nil
	case "u32":
		return &TypeTag{U32: &serialization.EmptyEnum{}}, nil
	case "u64":
		return &TypeTag{U64: &serialization.EmptyEnum{}}, nil
	case "u128":
		return &TypeTag{U128: &serialization.EmptyEnum{}}, nil
	case "u256":
		return &TypeTag{U256: &serialization.EmptyEnum{}}, nil
	}

	vectorRegex := regexp.MustCompile(`^vector<(.+)>$`)
	structRegex := regexp.MustCompile(`^([^:]+)::([^:]+)::([^<]+)(<(.+)>)?$`)

	vectorMatches := vectorRegex.FindStringSubmatch(data)
	if len(vectorMatches) != 0 {
		typeTag, err := TypeTagFromString(vectorMatches[1])
		if err != nil {
			return nil, fmt.Errorf("can't parse %s to TypeTag: %w", data, err)
		}
		return &TypeTag{Vector: typeTag}, nil
	}

	structMatches := structRegex.FindStringSubmatch(data)
	if len(structMatches) != 0 {
		structTag := &StructTag{
			Address: MustSuiAddressFromHex(structMatches[1]),
			Module:  Identifier(structMatches[2]),
			Name:    Identifier(structMatches[3]),
		}
		if len(structMatches) > 4 {
			typeTag, err := parseStructTypeArgs(structMatches[5])
			if err != nil {
				return nil, fmt.Errorf("can't parse TypeParams of a Struct in TypeParams: %w", err)
			}
			structTag.TypeParams = typeTag
		}
		return &TypeTag{Struct: structTag}, nil
	}
	return nil, fmt.Errorf("not found")
}

func parseStructTypeArgs(structTypeParams string) ([]TypeTag, error) {
	var retTypeTag []TypeTag
	tokens := splitGenericParameters(structTypeParams, nil)
	for _, tok := range tokens {
		elt, err := TypeTagFromString(tok)
		if err != nil {
			return nil, fmt.Errorf("can't parse struct tag args: %w", err)
		}
		retTypeTag = append(retTypeTag, *elt)
	}
	return retTypeTag, nil
}

type StructTag struct {
	Address    *SuiAddress
	Module     Identifier
	Name       Identifier
	TypeParams []TypeTag
}

func (s *StructTag) UnmarshalJSON(data []byte) error {
	str := string(data)
	str, _ = strings.CutPrefix(str, "\"")
	str, _ = strings.CutSuffix(str, "\"")
	parsedStructTag, err := StructTagFromString(str)
	if err != nil {
		return fmt.Errorf("can't unmarshal to %s StructTag: %w", str, err)
	}
	s.Address = parsedStructTag.Address
	s.Module = parsedStructTag.Module
	s.Name = parsedStructTag.Name
	s.TypeParams = parsedStructTag.TypeParams
	return nil
}

func StructTagFromString(data string) (*StructTag, error) {
	parts := strings.Split(data, "::")
	address, module := parts[0], parts[1]

	rest := data[len(address)+len(module)+4:]
	name := rest
	if strings.Contains(rest, "<") {
		name = rest[:strings.Index(rest, "<")]
	}

	typeParams := []TypeTag{}
	if strings.Contains(rest, "<") {
		typeParamsRawStr := rest[strings.Index(rest, "<")+1 : strings.LastIndex(rest, ">")]
		typeParamsTokens := splitGenericParameters(typeParamsRawStr, []string{"<", ">"})
		typeParams = make([]TypeTag, len(typeParamsTokens))
		for i, token := range typeParamsTokens {
			param := TypeTag{}
			if !strings.Contains(token, "::") {
				typeTag, err := TypeTagFromString(token)
				if err != nil {
					return nil, fmt.Errorf("can't parse TypeParams: %w", err)
				}
				param = *typeTag
			} else {
				typeParam, err := StructTagFromString(token)
				if err != nil {
					return nil, fmt.Errorf("can't parse StructTag TypeParams: %w", err)
				}
				param.Struct = typeParam
			}

			typeParams[i] = param
		}
	}

	if len(typeParams) == 0 {
		typeParams = nil
	}

	return &StructTag{
		Address:    MustSuiAddressFromHex(address),
		Module:     module,
		Name:       name,
		TypeParams: typeParams,
	}, nil
}

func splitGenericParameters(str string, genericSeparators []string) []string {
	var left, right string
	if genericSeparators != nil {
		left, right = genericSeparators[0], genericSeparators[1]
	} else {
		left, right = "", ""
	}

	var tokens []string
	word := ""
	nestedAngleBrackets := 0

	for i := 0; i < len(str); i++ {
		char := string(str[i])
		if char == left {
			nestedAngleBrackets++
		}
		if char == right {
			nestedAngleBrackets--
		}
		if nestedAngleBrackets == 0 && char == "," {
			tokens = append(tokens, strings.TrimSpace(word))
			word = ""
			continue
		}
		word += char
	}

	tokens = append(tokens, strings.TrimSpace(word))
	return tokens
}
