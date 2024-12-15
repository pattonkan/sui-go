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

	SubType *ResourceType
}

func NewResourceType(str string) (*ResourceType, error) {
	ltIdx := strings.Index(str, "<")
	var subType *ResourceType
	var err error
	if ltIdx != -1 {
		gtIdx := strings.LastIndex(str, ">")
		if gtIdx != len(str)-1 {
			return nil, errors.New("invalid type string literal")
		}
		subType, err = NewResourceType(str[ltIdx+1 : gtIdx])
		if err != nil {
			return nil, err
		}
		str = str[:ltIdx]
	}

	parts := strings.Split(str, "::")
	if len(parts) != 3 {
		return nil, errors.New("invalid type string literal")
	}
	addr, err := AddressFromHex(parts[0])
	if err != nil {
		return nil, err
	}
	return &ResourceType{
		Address:    addr,
		Module:     parts[1],
		ObjectName: parts[2],
		SubType:    subType,
	}, nil
}

func (r ResourceType) Contains(address *Address, moduleName string, funcName string) bool {
	if r.Module == moduleName && r.ObjectName == funcName {
		if address == nil {
			return true
		}
		if r.Address.String() == address.String() {
			return true
		}
	}
	for r = *r.SubType; r.SubType != nil; r = *r.SubType {
		if r.Module == moduleName && r.ObjectName == funcName {
			if address == nil {
				return true
			}
			if r.Address.String() == address.String() {
				return true
			}
		}
	}
	return false
}

func (t *ResourceType) String() string {
	if t.SubType != nil {
		return fmt.Sprintf("%v::%v::%v<%v>", t.Address.String(), t.Module, t.ObjectName, t.SubType.String())
	} else {
		return fmt.Sprintf("%v::%v::%v", t.Address.String(), t.Module, t.ObjectName)
	}
}

func (t *ResourceType) ShortString() string {
	if t.SubType != nil {
		return fmt.Sprintf("%v::%v::%v<%v>", t.Address.ShortString(), t.Module, t.ObjectName, t.SubType.ShortString())
	} else {
		return fmt.Sprintf("%v::%v::%v", t.Address.ShortString(), t.Module, t.ObjectName)
	}
}
