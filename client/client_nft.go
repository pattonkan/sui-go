package client

import (
	"context"

	"github.com/howjmay/go-sui-sdk/sui_types"
	"github.com/howjmay/go-sui-sdk/types"
)

// MintNFT
// Create an unsigned transaction to mint a nft at devnet
func (c *Client) MintNFT(
	ctx context.Context,
	signer *sui_types.SuiAddress,
	nftName, nftDescription, nftUri string,
	gas *sui_types.ObjectID,
	gasBudget uint64,
) (*types.TransactionBytes, error) {
	packageId, _ := sui_types.NewAddressFromHex("0x2")
	args := []any{
		nftName, nftDescription, nftUri,
	}
	return c.MoveCall(
		ctx,
		signer,
		packageId,
		"devnet_nft",
		"mint",
		[]string{},
		args,
		gas,
		types.NewSafeSuiBigInt(gasBudget),
	)
}

func (c *Client) GetNFTsOwnedByAddress(ctx context.Context, address *sui_types.SuiAddress) ([]types.SuiObjectResponse, error) {
	return c.BatchGetObjectsOwnedByAddress(
		ctx, address, &types.SuiObjectDataOptions{
			ShowType:    true,
			ShowContent: true,
			ShowOwner:   true,
		}, "0x2::devnet_nft::DevNetNFT",
	)
}
