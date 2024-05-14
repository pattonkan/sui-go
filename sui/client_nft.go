package sui

import (
	"context"

	"github.com/howjmay/sui-go/models"
	"github.com/howjmay/sui-go/sui_types"
)

// MintNFT
// Create an unsigned transaction to mint a nft at devnet
// func (s *ImplSuiAPI) MintNFT(
// 	ctx context.Context,
// 	signer *sui_types.SuiAddress,
// 	nftName, nftDescription, nftUri string,
// 	gas *sui_types.ObjectID,
// 	gasBudget uint64,
// ) (*models.TransactionBytes, error) {
// 	packageId, _ := sui_types.SuiAddressFromHex("0x2")
// 	args := []any{
// 		nftName, nftDescription, nftUri,
// 	}
// 	return s.MoveCall(
// 		ctx,
// 		signer,
// 		packageId,
// 		"devnet_nft",
// 		"mint",
// 		[]string{},
// 		ToSuiJsonArg(args),
// 		gas,
// 		models.NewSafeSuiBigInt(gasBudget),
// 	)
// }

func (s *ImplSuiAPI) GetNFTsOwnedByAddress(ctx context.Context, address *sui_types.SuiAddress) ([]models.SuiObjectResponse, error) {
	return s.BatchGetObjectsOwnedByAddress(
		ctx, address, &models.SuiObjectDataOptions{
			ShowType:    true,
			ShowContent: true,
			ShowOwner:   true,
		}, "0x2::devnet_nft::DevNetNFT",
	)
}
