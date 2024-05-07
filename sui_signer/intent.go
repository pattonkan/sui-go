package sui_signer

import (
	"bytes"

	"github.com/howjmay/sui-go/sui_types/serialization"

	"github.com/fardream/go-bcs/bcs"
)

type Intent struct {
	// the type of the IntentMessage
	Scope IntentScope
	// version the network supports
	Version IntentVersion
	// application that the signature refers to
	AppID AppID
}

func DefaultIntent() Intent {
	return Intent{
		Scope: IntentScope{
			TransactionData: &serialization.EmptyEnum{},
		},
		Version: IntentVersion{
			V0: &serialization.EmptyEnum{},
		},
		AppID: AppID{
			Sui: &serialization.EmptyEnum{},
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
	TransactionData         *serialization.EmptyEnum // Used for a user signature on a transaction data.
	TransactionEffects      *serialization.EmptyEnum // Used for an authority signature on transaction effects.
	CheckpointSummary       *serialization.EmptyEnum // Used for an authority signature on a checkpoint summary.
	PersonalMessage         *serialization.EmptyEnum // Used for a user signature on a personal message.
	SenderSignedTransaction *serialization.EmptyEnum // Used for an authority signature on a user signed transaction.
	ProofOfPossession       *serialization.EmptyEnum // Used as a signature representing an authority's proof of possession of its authority protocol key.
	HeaderDigest            *serialization.EmptyEnum // Used for narwhal authority signature on header digest.
}

func (i IntentScope) IsBcsEnum() {}

type IntentVersion struct {
	V0 *serialization.EmptyEnum
}

func (i IntentVersion) IsBcsEnum() {}

type AppID struct {
	Sui     *serialization.EmptyEnum
	Narwhal *serialization.EmptyEnum
}

func (a AppID) IsBcsEnum() {}

func MessageWithIntent(intent Intent, message []byte) []byte {
	intentMessage := bytes.NewBuffer(intent.Bytes())
	intentMessage.Write(message)
	return intentMessage.Bytes()
}
