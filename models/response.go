package models

import (
	"github.com/howjmay/sui-go/sui_types"
	"github.com/howjmay/sui-go/sui_types/serialization"
)

type AuthSignInfo interface{}

type CertifiedTransaction struct {
	TransactionDigest *sui_types.TransactionDigest `json:"transactionDigest"`
	TxSignature       string                       `json:"txSignature"`
	AuthSignInfo      *AuthSignInfo                `json:"authSignInfo"`

	Data *SenderSignedData `json:"data"`
}

type ParsedTransactionResponse interface{}

type ExecuteTransactionEffects struct {
	TransactionEffectsDigest *sui_types.TransactionDigest `json:"transactionEffectsDigest"`

	Effects      serialization.TagJson[SuiTransactionBlockEffects] `json:"effects"`
	AuthSignInfo *AuthSignInfo                                     `json:"authSignInfo"`
}

type ExecuteTransactionResponse struct {
	Certificate CertifiedTransaction      `json:"certificate"`
	Effects     ExecuteTransactionEffects `json:"effects"`

	ConfirmedLocalExecution bool `json:"confirmed_local_execution"`
}

func (r *ExecuteTransactionResponse) TransactionDigest() *sui_types.TransactionDigest {
	return r.Certificate.TransactionDigest
}

type SuiCoinMetadata struct {
	Decimals    uint8               `json:"decimals"`
	Name        string              `json:"name"`
	Symbol      string              `json:"symbol"`
	Description string              `json:"description"`
	IconUrl     string              `json:"iconUrl,omitempty"`
	Id          *sui_types.ObjectID `json:"id"`
}
