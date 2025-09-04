package sui

import (
	"errors"
	"fmt"
	"strings"
)

type ResourceType struct {
	Address    *Address
	Module     Identifier
	ObjectName Identifier // it can be function name or struct name, etc.

	SubTypes []*ResourceType
}

func IsSameResource(a, b string) (bool, error) {
	ra, err := NewResourceType(a)
	if err != nil {
		return false, fmt.Errorf("failed to parse resource type: %w", err)
	}
	rb, err := NewResourceType(b)
	if err != nil {
		return false, fmt.Errorf("failed to parse resource type: %w", err)
	}
	return ra.String() == rb.String(), nil
}

func MustNewResourceType(s string) *ResourceType {
	resource, err := NewResourceType(s)
	if err != nil {
		panic(err)
	}
	return resource
}

func NewResourceType(str string) (*ResourceType, error) {
	// Find the generic part <...>
	ltIdx := strings.Index(str, "<")
	var subTypes []*ResourceType
	var baseStr string

	if ltIdx != -1 {
		// Validate that > is the last character
		gtIdx := strings.LastIndex(str, ">")
		if gtIdx != len(str)-1 {
			return nil, errors.New("invalid type string literal")
		}

		// Extract the base type (before <)
		baseStr = str[:ltIdx]

		// Extract and parse subtypes
		genericPart := str[ltIdx+1 : gtIdx]
		if genericPart == "" {
			return nil, errors.New("empty generic parameters")
		}

		// Split subtypes by comma at depth 0
		subtypeStrs, err := splitSubtypes(genericPart)
		if err != nil {
			return nil, err
		}

		// Parse each subtype
		for _, subtypeStr := range subtypeStrs {
			subtype, err := NewResourceType(strings.TrimSpace(subtypeStr))
			if err != nil {
				return nil, err
			}
			subTypes = append(subTypes, subtype)
		}
	} else {
		baseStr = str
	}

	// Parse the base type (address::module::name)
	parts := strings.Split(baseStr, "::")
	if len(parts) != 3 {
		return nil, fmt.Errorf("invalid resource type string: %q", str)
	}

	addr, err := AddressFromHex(parts[0])
	if err != nil {
		return nil, err
	}

	module := parts[1]
	objectName := parts[2]

	return &ResourceType{
		Address:    addr,
		Module:     module,
		ObjectName: objectName,
		SubTypes:   subTypes,
	}, nil
}

// splitSubtypes splits a string by commas at depth 0 only
// e.g. "A<B,C>, D" -> ["A<B,C>", "D"]
func splitSubtypes(s string) ([]string, error) {
	var result []string
	var current strings.Builder
	depth := 0

	for _, r := range s {
		switch r {
		case '<':
			depth++
		case '>':
			depth--
			if depth < 0 {
				return nil, errors.New("unmatched closing bracket")
			}
		case ',':
			if depth == 0 {
				// Split here
				part := strings.TrimSpace(current.String())
				if part == "" {
					return nil, errors.New("empty subtype entry")
				}
				result = append(result, part)
				current.Reset()
				continue
			}
		}
		current.WriteRune(r)
	}

	if depth != 0 {
		return nil, errors.New("unmatched brackets")
	}

	// Add the last part
	part := strings.TrimSpace(current.String())
	if part == "" {
		return nil, errors.New("empty subtype entry")
	}
	result = append(result, part)

	return result, nil
}

func (t *ResourceType) UnmarshalJSON(data []byte) error {
	resource, err := NewResourceType(string(data[1 : len(data)-1]))
	if err != nil {
		return err
	}
	*t = *resource
	return nil
}

func (t *ResourceType) Contains(address *Address, moduleName string, funcName string) bool {
	if t == nil {
		return false
	}

	// Check if this type matches
	if t.Module == moduleName && t.ObjectName == funcName {
		if address == nil {
			return true
		}
		if t.Address.String() == address.String() {
			return true
		}
	}

	// Recursively check all subtypes
	for _, subType := range t.SubTypes {
		if subType.Contains(address, moduleName, funcName) {
			return true
		}
	}

	return false
}

func (t *ResourceType) String() string {
	if len(t.SubTypes) == 0 {
		return fmt.Sprintf("%v::%v::%v", t.Address.String(), t.Module, t.ObjectName)
	}

	var subtypeStrs []string
	for _, subtype := range t.SubTypes {
		subtypeStrs = append(subtypeStrs, subtype.String())
	}

	return fmt.Sprintf("%v::%v::%v<%v>",
		t.Address.String(),
		t.Module,
		t.ObjectName,
		strings.Join(subtypeStrs, ", "))
}

func (t *ResourceType) ShortString() string {
	if len(t.SubTypes) == 0 {
		return fmt.Sprintf("%v::%v::%v", t.Address.ShortString(), t.Module, t.ObjectName)
	}

	var subtypeStrs []string
	for _, subtype := range t.SubTypes {
		subtypeStrs = append(subtypeStrs, subtype.ShortString())
	}

	return fmt.Sprintf("%v::%v::%v<%v>",
		t.Address.ShortString(),
		t.Module,
		t.ObjectName,
		strings.Join(subtypeStrs, ", "))
}
