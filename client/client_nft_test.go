package client_test

// FIXME implement nft contract

// func TestMintNFT(t *testing.T) {
// 	cli := TestnetClient(t)

// 	var (
// 		timeNow = time.Now().Format("06-01-02 15:04")
// 		nftName = "ComingChat NFT at " + timeNow
// 		nftDesc = "This is a NFT created by ComingChat"
// 		nftUrl  = "https://coming.chat/favicon.ico"
// 	)
// 	coins, err := cli.GetSuiCoinsOwnedByAddress(context.TODO(), account.TEST_ADDRESS)
// 	require.NoError(t, err)

// 	firstCoin, err := coins.PickCoinNoLess(12000)
// 	require.NoError(t, err)

// 	txnBytes, err := cli.MintNFT(context.TODO(), account.TEST_ADDRESS, nftName, nftDesc, nftUrl, &firstCoin.CoinObjectID, 12000)
// 	require.NoError(t, err)
// 	t.Log(txnBytes.TxBytes)

// 	resp, err := cli.DryRunTransaction(context.TODO(), txnBytes.TxBytes)
// 	require.NoError(t, err)
// 	require.True(t, resp.Effects.Data.IsSuccess())
// 	require.Empty(t, resp.Effects.Data.V1.Status.Error)
// 	t.Logf("%#v", resp)
// }

// func TestGetDevNFTs(t *testing.T) {
// 	cli := TestnetClient(t)

// 	nfts, err := cli.GetNFTsOwnedByAddress(context.TODO(), account.TEST_ADDRESS)
// 	require.NoError(t, err)
// 	for _, nft := range nfts {
// 		t.Log(nft.Data)
// 	}
// }
