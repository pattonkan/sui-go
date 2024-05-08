package isc_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/howjmay/sui-go/isc"
	"github.com/howjmay/sui-go/sui"
	"github.com/howjmay/sui-go/sui/conn"
	"github.com/howjmay/sui-go/sui_signer"
	"github.com/stretchr/testify/require"
)

type Client struct {
	API sui.ImplSuiAPI
}

func TestStartNewChain(t *testing.T) {
	t.Skip("only for localnet")
	client := isc.NewIscClient(sui.NewSuiClient(conn.LocalnetEndpointUrl))

	signer, err := sui_signer.NewSignerWithMnemonic(sui_signer.TEST_MNEMONIC)
	require.NoError(t, err)

	t.Log("sui_signer: ", signer.Address)
	digest, err := sui.RequestFundFromFaucet(signer.Address, conn.LocalnetFaucetUrl)
	require.NoError(t, err)
	t.Log("digest: ", digest)

	packageID, anchorCap := isc.GetIscPackageIDAndAnchor(isc.GetGitRoot() + "/isc/contracts/isc/publish_receipt.json")

	res, err := client.StartNewChain(context.Background(), signer, packageID, anchorCap)
	require.NoError(t, err)
	t.Logf("StartNewChain response: %#v\n", res)
	for _, change := range res.ObjectChanges {
		fmt.Println("change.Data.Created: ", change.Data.Created)
	}
}
