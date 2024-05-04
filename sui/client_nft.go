package sui

import (
	"context"

	"github.com/howjmay/sui-go/sui_types"
	"github.com/howjmay/sui-go/types"
)

// MintNFT
// Create an unsigned transaction to mint a nft at devnet
func (s *ImplSuiAPI) MintNFT(
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
	return s.MoveCall(
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

func (s *ImplSuiAPI) GetNFTsOwnedByAddress(ctx context.Context, address *sui_types.SuiAddress) ([]types.SuiObjectResponse, error) {
	return s.BatchGetObjectsOwnedByAddress(
		ctx, address, &types.SuiObjectDataOptions{
			ShowType:    true,
			ShowContent: true,
			ShowOwner:   true,
		}, "0x2::devnet_nft::DevNetNFT",
	)
}
