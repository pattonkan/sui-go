package signer_test

import (
	"encoding/json"
	"testing"

	"github.com/howjmay/sui-go/signer"
	"github.com/stretchr/testify/require"
)

func TestAccount(t *testing.T) {
	account, err := signer.NewAccountWithMnemonic(signer.TEST_MNEMONIC)
	require.NoError(t, err)

	t.Logf("addr: %v", account.Address)
}

func Test_Signature_Marshal_Unmarshal(t *testing.T) {
	account, err := signer.NewAccountWithMnemonic(signer.TEST_MNEMONIC)
	require.NoError(t, err)

	msg := "Coming chat is very good jopfpzf"
	msgBytes := []byte(msg)

	signature1, err := account.SignTransactionBlock(msgBytes, signer.DefaultIntent())
	require.NoError(t, err)

	marshaledData, err := json.Marshal(signature1)
	require.NoError(t, err)

	var signature2 signer.Signature
	err = json.Unmarshal(marshaledData, &signature2)
	require.NoError(t, err)

	require.Equal(t, signature1, signature2)
}
