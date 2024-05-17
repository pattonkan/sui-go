package sui_types_test

import (
	"testing"

	"github.com/fardream/go-bcs/bcs"
	"github.com/howjmay/sui-go/sui_types"
	"github.com/stretchr/testify/require"
)

func TestTransferSui(t *testing.T) {
	recipient, err := sui_types.SuiAddressFromHex("0x7e875ea78ee09f08d72e2676cf84e0f1c8ac61d94fa339cc8e37cace85bebc6e")
	require.NoError(t, err)

	ptb := sui_types.NewProgrammableTransactionBuilder()
	amount := uint64(100000)
	err = ptb.TransferSui(recipient, &amount)
	require.NoError(t, err)
	pt := ptb.Finish()
	digest := sui_types.MustNewDigest("HvbE2UZny6cP4KukaXetmj4jjpKTDTjVo23XEcu7VgSn")
	objectId, err := sui_types.SuiAddressFromHex("0x13c1c3d0e15b4039cec4291c75b77c972c10c8e8e70ab4ca174cf336917cb4db")
	require.NoError(t, err)
	tx := sui_types.NewProgrammable(
		recipient,
		pt,
		[]*sui_types.ObjectRef{
			{
				ObjectID: objectId,
				Version:  14924029,
				Digest:   digest,
			},
		},
		10000000,
		1000,
	)
	txByte, err := bcs.Marshal(tx)
	require.NoError(t, err)
	t.Logf("%x", txByte)
}
