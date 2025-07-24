package randtypes

import (
	crypto_rand "crypto/rand"
	"math/rand"

	"github.com/pattonkan/sui-go/sui"
	"github.com/pattonkan/sui-go/sui/suiptb"
)

func RandomObjectRef() *sui.ObjectRef {
	return &sui.ObjectRef{
		ObjectId: RandomAddress(),
		Version:  rand.Uint64(),
		Digest:   RandomDigest(),
	}
}

func RandomAddress() *sui.Address {
	var a sui.Address
	_, _ = crypto_rand.Read(a[:])
	return &a
}

func RandomDigest() *sui.Digest {
	var d sui.Digest
	_, _ = crypto_rand.Read(d[:])
	return &d
}

func RandomTransactionData() *suiptb.TransactionData {
	ptb := suiptb.NewTransactionDataTransactionBuilder()
	ptb.Command(
		suiptb.Command{
			MoveCall: &suiptb.ProgrammableMoveCall{
				Package:       RandomAddress(),
				Module:        "test_module",
				Function:      "test_func",
				TypeArguments: []sui.TypeTag{},
				Arguments:     []suiptb.Argument{},
			},
		},
	)
	pt := ptb.Finish()
	tx := suiptb.NewTransactionData(
		RandomAddress(),
		pt,
		[]*sui.ObjectRef{},
		10000,
		100,
	)
	return &tx
}
