package isc

import "github.com/howjmay/sui-go/models"

type NormalizedID struct {
	ID string `json:"id"`
}

type NormalizedAssets struct {
	Type   string `json:"type"`
	Fields struct {
		Coins struct {
			Type   string `json:"type"`
			Fields struct {
				ID NormalizedID `json:"id"`
				// FIXME this should be int
				Size string `json:"size"`
			} `json:"fields"`
		} `json:"coins"`
		ID   NormalizedID  `json:"id"`
		Nfts []interface{} `json:"nfts"`
	} `json:"fields"`
}

type Assets struct {
	Coins models.Coins  // only some fields in models.Coin will have value here (having values in some fields are just don't make sense)
	Nfts  []interface{} `json:"nfts"` // FIXME find the definition of NFT object
}

// $sui client object 0x44ec93c956416bed9772abe8ecc63efa8fcceeb2a5838e4198d7a441852d3069 --json
// {
// 	"objectId": "0x44ec93c956416bed9772abe8ecc63efa8fcceeb2a5838e4198d7a441852d3069",
// 	"version": "51",
// 	"digest": "9EVTY4dg7vPkAgv5ye8dFYQ1x92nZ1tqjaakhX9T1yVU",
// 	"type": "0xe9d253fc021cc41ad4fbb3c673a550b4b0fad7d7d82ab3d17fb628a00671dbd9::anchor::Anchor",
// 	"owner": {
// 	  "AddressOwner": "0x1a02d61c6434b4d0ff252a880c04050b5f27c8b574026c98dd72268865c0ede5"
// 	},
// 	"previousTransaction": "2v1axtyHWzC9QWuQcStx7Gszq485GgRU6LGaDaVp86J5",
// 	"storageRebate": "1839200",
// 	"content": {
// 	  "dataType": "moveObject",
// 	  "type": "0xe9d253fc021cc41ad4fbb3c673a550b4b0fad7d7d82ab3d17fb628a00671dbd9::anchor::Anchor",
// 	  "hasPublicTransfer": true,
// 	  "fields": {
// 		"assets": {
// 		  "type": "0xe9d253fc021cc41ad4fbb3c673a550b4b0fad7d7d82ab3d17fb628a00671dbd9::assets::Assets",
// 		  "fields": {
// 			"coins": {
// 			  "type": "0x2::bag::Bag",
// 			  "fields": {
// 				"id": {
// 				  "id": "0x480fbc345e43e50cde1bbb087f32b9e3c79f9bb4ec4d74ca0aff29f9a8b22fd8"
// 				},
// 				"size": "1"
// 			  }
// 			},
// 			"id": {
// 			  "id": "0x8e0e76c9cc9d8c8b8f1a0884c53f5aec750cbd8d1b1c9cbb1ebd528db8c2402a"
// 			},
// 			"nfts": []
// 		  }
// 		},
// 		"id": {
// 		  "id": "0x44ec93c956416bed9772abe8ecc63efa8fcceeb2a5838e4198d7a441852d3069"
// 		}
// 	  }
// 	}
//   }

//   ```
//   $sui client dynamic-field 0x480fbc345e43e50cde1bbb087f32b9e3c79f9bb4ec4d74ca0aff29f9a8b22fd8 --json
//   {
// 	"data": [
// 	  {
// 		"name": {
// 		  "type": "0x1::ascii::String",
// 		  "value": "54cd5a30b529b7aa513232b15c2bacdd0060652a17477860fde8eb5b13816f99::testcoin::TESTCOIN"
// 		},
// 		"bcsName": "TWA81MAr59kXzDvS2JQP1UegBNVdjMtCNwuxy4bo8fN9wu91hmwggyWsbridQU4bfkvvCg2c1x3A2RsRR8eXEptzczFRbYS4qbLdKgJrW3et5ULZSsSM",
// 		"type": "DynamicField",
// 		"objectType": "0x2::balance::Balance<0x54cd5a30b529b7aa513232b15c2bacdd0060652a17477860fde8eb5b13816f99::testcoin::TESTCOIN>",
// 		"objectId": "0xe1df36d69212f815657f9416c98ddba78c9b6fb4007d928ca5a563d91cb01891",
// 		"version": 51,
// 		"digest": "4rVDt2Qzi1d5tnuU9Teiz3vNqBB25wtmxK4ZFULUTq2F"
// 	  }
// 	],
// 	"nextCursor": "0xe1df36d69212f815657f9416c98ddba78c9b6fb4007d928ca5a563d91cb01891",
// 	"hasNextPage": false
//   }
//   ```

//   ```
//   $sui client object 0xe1df36d69212f815657f9416c98ddba78c9b6fb4007d928ca5a563d91cb01891 --json
//   {
// 	"objectId": "0xe1df36d69212f815657f9416c98ddba78c9b6fb4007d928ca5a563d91cb01891",
// 	"version": "51",
// 	"digest": "4rVDt2Qzi1d5tnuU9Teiz3vNqBB25wtmxK4ZFULUTq2F",
// 	"type": "0x2::dynamic_field::Field<0x1::ascii::String, 0x2::balance::Balance<0x54cd5a30b529b7aa513232b15c2bacdd0060652a17477860fde8eb5b13816f99::testcoin::TESTCOIN>>",
// 	"owner": {
// 	  "ObjectOwner": "0x480fbc345e43e50cde1bbb087f32b9e3c79f9bb4ec4d74ca0aff29f9a8b22fd8"
// 	},
// 	"previousTransaction": "2v1axtyHWzC9QWuQcStx7Gszq485GgRU6LGaDaVp86J5",
// 	"storageRebate": "3169200",
// 	"content": {
// 	  "dataType": "moveObject",
// 	  "type": "0x2::dynamic_field::Field<0x1::ascii::String, 0x2::balance::Balance<0x54cd5a30b529b7aa513232b15c2bacdd0060652a17477860fde8eb5b13816f99::testcoin::TESTCOIN>>",
// 	  "hasPublicTransfer": false,
// 	  "fields": {
// 		"id": {
// 		  "id": "0xe1df36d69212f815657f9416c98ddba78c9b6fb4007d928ca5a563d91cb01891"
// 		},
// 		"name": "54cd5a30b529b7aa513232b15c2bacdd0060652a17477860fde8eb5b13816f99::testcoin::TESTCOIN",
// 		"value": "1000000"
// 	  }
// 	}
//   }
//   ```
