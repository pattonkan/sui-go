package models_test

import (
	"encoding/json"
	"testing"

	"github.com/howjmay/sui-go/models"
	"github.com/howjmay/sui-go/sui_types"
	"github.com/stretchr/testify/require"
)

func TestSuiEventDecode(t *testing.T) {
	receivingMessage := []byte(`{
      "id": {
        "txDigest": "EJthSfz1GvtoJ17L8AfCWRTRsKRhsJgG5Q4Zxm7QY6Ag",
        "eventSeq": "0"
      },
      "packageId": "0x000000000000000000000000000000000000000000000000000000000000dee9",
      "transactionModule": "clob_v2",
      "sender": "0x4a5d59860daf389d0a34c321fc486aae3a4c1eb3666b684e6c51ff8163c58cc5",
      "type": "0xdee9::clob_v2::OrderPlaced<0x2::sui::SUI, 0x5d4b302506645c37ff133b98c4b50a5ae14841659738d6d733d59d0d217a93bf::coin::COIN>",
      "parsedJson": {
        "base_asset_quantity_placed": "1752800000000",
        "client_order_id": "4459280326793917790",
        "expire_timestamp": "1721201284399",
        "is_bid": false,
        "order_id": "9223372036863437941",
        "original_quantity": "1752800000000",
        "owner": "0xa54f886de28b9f23b10e6fa682393a698805984129cb8ed3a8dd42c7acf4285b",
        "pool_id": "0x4405b50d791fd3346754e8171aaab6bc2ed26c2c46efdd033c14b30ae507ac33",
        "price": "863800"
      },
      "bcs": "yNS5iDS3Gvdo3DhXdtFpuTS12RrSiNkrvjcm2rejntCt1jjM24MvS8o5s9dDHCe4pxwJvVbm1pBWymUgAGmG7SAZpxJSocaXSFj5oW1YT2KkZoEieUfPRbmcmf6fAHYHi5Nb1Hq73RFXjXG5soSUsgZyrP",
      "timestampMs": "1721197686017"
    }`)
	var event models.SuiEvent
	err := json.Unmarshal(receivingMessage, &event)
	require.NoError(t, err)
	require.Equal(t, sui_types.MustPackageIDFromHex("0x000000000000000000000000000000000000000000000000000000000000dee9"), event.PackageId)
}
