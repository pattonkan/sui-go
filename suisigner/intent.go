package suisigner

import (
	"bytes"

	"github.com/fardream/go-bcs/bcs"
	"github.com/pattonkan/sui-go/sui"
)

type Intent struct {
	// the type of the IntentMessage
	Scope IntentScope
	// version the network supports
	Version IntentVersion
	// application that the signature refers to
	AppId AppId
}

func DefaultIntent() Intent {
	return Intent{
		Scope: IntentScope{
			TransactionData: &sui.EmptyEnum{},
		},
		Version: IntentVersion{
			V0: &sui.EmptyEnum{},
		},
		AppId: AppId{
			Sui: &sui.EmptyEnum{},
		},
	}
}

func (i *Intent) Bytes() []byte {
	b, err := bcs.Marshal(i)
	if err != nil {
		return nil
	}
	return b
}

// the type of the IntentMessage
type IntentScope struct {
	TransactionData         *sui.EmptyEnum // Used for a user signature on a transaction data.
	TransactionEffects      *sui.EmptyEnum // Used for an authority signature on transaction effects.
	CheckpointSummary       *sui.EmptyEnum // Used for an authority signature on a checkpoint summary.
	PersonalMessage         *sui.EmptyEnum // Used for a user signature on a personal message.
	SenderSignedTransaction *sui.EmptyEnum // Used for an authority signature on a user signed transaction.
	ProofOfPossession       *sui.EmptyEnum // Used as a signature representing an authority's proof of possession of its authority protocol key.
	HeaderDigest            *sui.EmptyEnum // Used for narwhal authority signature on header digest.
}

func (i IntentScope) IsBcsEnum() {}

type IntentVersion struct {
	V0 *sui.EmptyEnum
}

func (i IntentVersion) IsBcsEnum() {}

type AppId struct {
	Sui     *sui.EmptyEnum
	Narwhal *sui.EmptyEnum
}

func (a AppId) IsBcsEnum() {}

func MessageWithIntent(intent Intent, message []byte) []byte {
	intentMessage := bytes.NewBuffer(intent.Bytes())
	intentMessage.Write(message)
	return intentMessage.Bytes()
}
