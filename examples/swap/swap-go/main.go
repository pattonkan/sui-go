package main

import (
	"context"
	"fmt"

	"github.com/howjmay/sui-go/examples/swap/swap-go/pkg"
	"github.com/howjmay/sui-go/models"
	"github.com/howjmay/sui-go/sui"
	"github.com/howjmay/sui-go/sui/conn"
	"github.com/howjmay/sui-go/sui_signer"
	"github.com/howjmay/sui-go/utils"
)

func main() {
	suiClient, signer := sui.NewSuiClient(conn.TestnetEndpointUrl).WithSignerAndFund(sui_signer.TEST_SEED, 0)
	_, swapper := sui.NewSuiClient(conn.TestnetEndpointUrl).WithSignerAndFund(sui_signer.TEST_SEED, 1)
	fmt.Println("signer: ", signer.Address)
	fmt.Println("swapper: ", swapper.Address)

	swapPackageID := pkg.BuildAndPublish(suiClient, signer, utils.GetGitRoot()+"/examples/swap/swap")
	testcoinID, _ := pkg.BuildDeployMintTestcoin(suiClient, signer)
	testcoinCoinType := fmt.Sprintf("%s::testcoin::TESTCOIN", testcoinID.String())

	fmt.Println("swapPackageID: ", swapPackageID)
	fmt.Println("testcoinCoinType: ", testcoinCoinType)

	testcoinCoins, err := suiClient.GetCoins(
		context.Background(),
		&models.GetCoinsRequest{
			Owner:    signer.Address,
			CoinType: &testcoinCoinType,
		},
	)
	if err != nil {
		panic(err)
	}

	signerSuiCoinPage, err := suiClient.GetCoins(
		context.Background(),
		&models.GetCoinsRequest{
			Owner: signer.Address,
		},
	)
	if err != nil {
		panic(err)
	}

	poolObjectID := pkg.CreatePool(suiClient, signer, swapPackageID, testcoinID, testcoinCoins.Data[0], signerSuiCoinPage.Data)

	swapperSuiCoinPage1, err := suiClient.GetAllCoins(
		context.Background(),
		&models.GetAllCoinsRequest{Owner: swapper.Address},
	)
	if err != nil {
		panic(err)
	}

	fmt.Println("poolObjectID: ", poolObjectID)

	pkg.SwapSui(suiClient, swapper, swapPackageID, testcoinID, poolObjectID, swapperSuiCoinPage1.Data)

	swapperSuiCoinPage2, err := suiClient.GetAllCoins(
		context.Background(),
		&models.GetAllCoinsRequest{Owner: swapper.Address},
	)
	if err != nil {
		panic(err)
	}
	fmt.Println("swapper now has")
	for _, coin := range swapperSuiCoinPage2.Data {
		fmt.Printf("object: %s in type: %s\n", coin.CoinObjectID, coin.CoinType)
	}
}
