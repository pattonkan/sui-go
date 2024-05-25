package models

type Supply struct {
	Value SafeSuiBigInt[uint64] `json:"value"`
}
