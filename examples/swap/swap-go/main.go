package main

import (
	"context"
	"fmt"

	"github.com/pattonkan/sui-go/examples/swap/swap-go/pkg"
	"github.com/pattonkan/sui-go/suiclient"
	"github.com/pattonkan/sui-go/suiclient/conn"
	"github.com/pattonkan/sui-go/suisigner"
	"github.com/pattonkan/sui-go/utils"
)

func main() {
	suiClient, signer := suiclient.NewClient(conn.TestnetEndpointUrl).WithSignerAndFund(suisigner.TEST_SEED, 0)
	_, swapper := suiclient.NewClient(conn.TestnetEndpointUrl).WithSignerAndFund(suisigner.TEST_SEED, 1)
	fmt.Println("signer: ", signer.Address)
	fmt.Println("swapper: ", swapper.Address)

	swapPackageId := pkg.BuildAndPublish(suiClient, signer, utils.GetGitRoot()+"/examples/swap/swap")
	testcoinId, _ := pkg.BuildDeployMintTestcoin(suiClient, signer)
	testcoinCoinType := fmt.Sprintf("%s::testcoin::TESTCOIN", testcoinId.String())

	fmt.Println("swapPackageId: ", swapPackageId)
	fmt.Println("testcoinCoinType: ", testcoinCoinType)

	testcoinCoins, err := suiClient.GetCoins(
		context.Background(),
		&suiclient.GetCoinsRequest{
			Owner:    signer.Address,
			CoinType: &testcoinCoinType,
		},
	)
	if err != nil {
		panic(err)
	}

	signerSuiCoinPage, err := suiClient.GetCoins(
		context.Background(),
		&suiclient.GetCoinsRequest{
			Owner: signer.Address,
		},
	)
	if err != nil {
		panic(err)
	}

	poolObjectId := pkg.CreatePool(suiClient, signer, swapPackageId, testcoinId, testcoinCoins.Data[0], signerSuiCoinPage.Data)

	swapperSuiCoinPage1, err := suiClient.GetAllCoins(
		context.Background(),
		&suiclient.GetAllCoinsRequest{Owner: swapper.Address},
	)
	if err != nil {
		panic(err)
	}

	fmt.Println("poolObjectId: ", poolObjectId)

	pkg.SwapSui(suiClient, swapper, swapPackageId, testcoinId, poolObjectId, swapperSuiCoinPage1.Data)

	swapperSuiCoinPage2, err := suiClient.GetAllCoins(
		context.Background(),
		&suiclient.GetAllCoinsRequest{Owner: swapper.Address},
	)
	if err != nil {
		panic(err)
	}
	fmt.Println("swapper now has")
	for _, coin := range swapperSuiCoinPage2.Data {
		fmt.Printf("object: %s in type: %s\n", coin.CoinObjectId, coin.CoinType)
	}
}
