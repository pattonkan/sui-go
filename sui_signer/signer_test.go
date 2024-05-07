package sui_signer_test

import (
	"encoding/json"
	"testing"

	"github.com/howjmay/sui-go/sui_signer"

	"github.com/stretchr/testify/require"
)

func TestAccount(t *testing.T) {
	signer, err := sui_signer.NewSignerWithMnemonic(sui_signer.TEST_MNEMONIC)
	require.NoError(t, err)

	t.Logf("addr: %v", signer.Address)
}

func Test_Signature_Marshal_Unmarshal(t *testing.T) {
	signer, err := sui_signer.NewSignerWithMnemonic(sui_signer.TEST_MNEMONIC)
	require.NoError(t, err)

	msg := "I want to have some bubble tea"
	msgBytes := []byte(msg)

	signature1, err := signer.SignTransactionBlock(msgBytes, sui_signer.DefaultIntent())
	require.NoError(t, err)

	marshaledData, err := json.Marshal(signature1)
	require.NoError(t, err)

	var signature2 sui_signer.Signature
	err = json.Unmarshal(marshaledData, &signature2)
	require.NoError(t, err)

	require.Equal(t, signature1, signature2)
}
