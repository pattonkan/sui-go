package suiptb

import (
	"github.com/fardream/go-bcs/bcs"
	"github.com/pattonkan/sui-go/sui"
	"github.com/pattonkan/sui-go/suisigner"
	"golang.org/x/crypto/blake2b"
)

var (
	SuiSystemMut = CallArg{
		Object: &SuiSystemMutObj,
	}

	SuiSystemMutObj = ObjectArg{
		SharedObject: &SharedObjectArg{
			Id:                   sui.SuiObjectIdSystemState,
			InitialSharedVersion: sui.SuiSystemStateObjectSharedVersion,
			Mutable:              true,
		},
	}
)

type TransactionData struct {
	V1 *TransactionDataV1
}

func (t TransactionData) IsBcsEnum() {}

func (t TransactionData) SigningDigest() ([]byte, error) {
	intent := suisigner.IntentTransaction()
	msg, err := bcs.Marshal(t)
	if err != nil {
		return nil, err
	}
	data := suisigner.MessageWithIntent(intent, msg)
	hash := blake2b.Sum256(data)
	return hash[:], nil
}

type TransactionDataV1 struct {
	Kind       TransactionKind
	Sender     sui.Address
	GasData    GasData
	Expiration TransactionExpiration
}

type TransactionKind struct {
	ProgrammableTransaction *ProgrammableTransaction
	ChangeEpoch             *ChangeEpoch
	Genesis                 *GenesisTransaction
	ConsensusCommitPrologue *ConsensusCommitPrologue
}

func (t TransactionKind) IsBcsEnum() {}

func NewTransactionData(
	sender *sui.Address,
	pt ProgrammableTransaction,
	gasPayment []*sui.ObjectRef,
	gasBudget uint64,
	gasPrice uint64,
) TransactionData {
	return NewTransactionDataAllowSponsor(*sender, pt, gasPayment, gasBudget, gasPrice, sender)
}

// This 'TransactionData' will need to be signed by both 'sponsor' and 'sender'
func NewTransactionDataAllowSponsor(
	sender sui.Address,
	pt ProgrammableTransaction,
	gasPayment []*sui.ObjectRef,
	gasBudget uint64,
	gasPrice uint64,
	sponsor *sui.Address,
) TransactionData {
	kind := TransactionKind{
		ProgrammableTransaction: &pt,
	}
	return TransactionData{
		V1: &TransactionDataV1{
			Kind:   kind,
			Sender: sender,
			GasData: GasData{
				Payment: gasPayment,
				Owner:   sponsor,
				Price:   gasPrice,
				Budget:  gasBudget,
			},
			Expiration: TransactionExpiration{
				None: &sui.EmptyEnum{},
			},
		},
	}
}
