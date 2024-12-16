package suiclient

import (
	"github.com/pattonkan/sui-go/sui"
)

type DynamicFieldType struct {
	DynamicField  *sui.EmptyEnum `json:"DynamicField"`
	DynamicObject *sui.EmptyEnum `json:"DynamicObject"`
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

type DynamicFieldInfo struct {
	Name       DynamicFieldName                    `json:"name"`
	BcsName    sui.Base58                          `json:"bcsName"`
	Type       WrapperTaggedJson[DynamicFieldType] `json:"type"`
	ObjectType sui.ObjectType                      `json:"objectType"`
	ObjectId   sui.ObjectId                        `json:"objectId"`
	Version    sui.SequenceNumber                  `json:"version"`
	Digest     sui.ObjectDigest                    `json:"digest"`
}

type DynamicFieldPage = Page[DynamicFieldInfo, sui.ObjectId]
