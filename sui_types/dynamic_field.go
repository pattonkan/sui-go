package sui_types

import "github.com/howjmay/sui-go/sui_types/serialization"

type DynamicFieldType struct {
	DynamicField  *serialization.EmptyEnum `json:"DynamicField"`
	DynamicObject *serialization.EmptyEnum `json:"DynamicObject"`
}

func (d DynamicFieldType) Tag() string {
	return ""
}

func (d DynamicFieldType) Content() string {
	return ""
}

type DynamicFieldName struct {
	Type  string `json:"type"` // TODO Maybe sui_types.ObjectType type
	Value any    `json:"value"`
}
