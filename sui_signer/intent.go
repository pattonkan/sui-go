package sui_signer

import (
	"bytes"

	"github.com/fardream/go-bcs/bcs"

	"github.com/howjmay/sui-go/lib"
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
			TransactionData: &lib.EmptyEnum{},
		},
		Version: IntentVersion{
			V0: &lib.EmptyEnum{},
		},
		AppID: AppID{
			Sui: &lib.EmptyEnum{},
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

type AppID struct {
	Sui     *lib.EmptyEnum
	Narwhal *lib.EmptyEnum
}

func (a AppID) IsBcsEnum() {}

func MessageWithIntent(intent Intent, message []byte) []byte {
	intentMessage := bytes.NewBuffer(intent.Bytes())
	intentMessage.Write(message)
	return intentMessage.Bytes()
}
