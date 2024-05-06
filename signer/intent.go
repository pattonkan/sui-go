// package signer

// type AppId int

// const (
// 	Sui AppId = 0
// )

// type IntentVersion int

// const (
// 	V0 IntentVersion = 0
// )

// type IntentScope int

// const (
// 	TransactionData    IntentScope = 0
// 	TransactionEffects IntentScope = 1
// 	CheckpointSummary  IntentScope = 2
// 	PersonalMessage    IntentScope = 3
// )

// func IntentWithScope(intentScope IntentScope) []int {
// 	return []int{int(intentScope), int(V0), int(Sui)}
// }

package signer

import (
	"github.com/fardream/go-bcs/bcs"
	"github.com/howjmay/sui-go/lib"
)

type IntentScope struct {
	TransactionData         *lib.EmptyEnum // Used for a user signature on a transaction data.
	TransactionEffects      *lib.EmptyEnum // Used for an authority signature on transaction effects.
	CheckpointSummary       *lib.EmptyEnum // Used for an authority signature on a checkpoint summary.
	PersonalMessage         *lib.EmptyEnum // Used for a user signature on a personal message.
	SenderSignedTransaction *lib.EmptyEnum // Used for an authority signature on a user signed transaction.
	ProofOfPossession       *lib.EmptyEnum // Used as a signature representing an authority's proof of possession of its authority protocol key.
	HeaderDigest            *lib.EmptyEnum // Used for narwhal authority signature on header digest.
}

func (i IntentScope) IsBcsEnum() {}

type IntentVersion struct {
	V0 *lib.EmptyEnum
}

func (i IntentVersion) IsBcsEnum() {}

type AppId struct {
	Sui     *lib.EmptyEnum
	Narwhal *lib.EmptyEnum
}

func (a AppId) IsBcsEnum() {}

type Intent struct {
	Scope   IntentScope
	Version IntentVersion
	AppId   AppId
}

func DefaultIntent() Intent {
	return Intent{
		Scope: IntentScope{
			TransactionData: &lib.EmptyEnum{},
		},
		Version: IntentVersion{
			V0: &lib.EmptyEnum{},
		},
		AppId: AppId{
			Sui: &lib.EmptyEnum{},
		},
	}
}

// type IntentValue interface {
// 	TransactionData | ~[]byte
// }

// type IntentMessage[T IntentValue] struct {
// 	Intent Intent
// 	Value  T
// }

// func MessageWithIntent[T IntentValue](intent Intent, value T) IntentMessage[T] {
// 	return IntentMessage[T]{
// 		Intent: intent,
// 		Value:  value,
// 	}
// }

func (i *Intent) Bytes() []byte {
	b, err := bcs.Marshal(i)
	if err != nil {
		return nil
	}
	return b
}

func MessageWithIntent(intent Intent, message []byte) []byte {
	intentBytes := intent.Bytes()
	intentMessage := make([]byte, len(intentBytes)+len(message))
	copy(intentMessage, intentBytes)
	copy(intentMessage[len(intentBytes):], message)
	return intentMessage
}
