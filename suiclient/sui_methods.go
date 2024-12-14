package suiclient

type SuiMethod string

func (s SuiMethod) String() string {
	return string(s)
}

type SuiXMethod string

func (s SuiXMethod) String() string {
	return string(s)
}

type UnsafeMethod string

func (u UnsafeMethod) String() string {
	return string(u)
}

const (
	// Coin Query API
	getAllBalances  SuiXMethod = "suix_getAllBalances"
	getAllCoins     SuiXMethod = "suix_getAllCoins"
	getBalance      SuiXMethod = "suix_getBalance"
	getCoinMetadata SuiXMethod = "suix_getCoinMetadata"
	getCoins        SuiXMethod = "suix_getCoins"
	getTotalSupply  SuiXMethod = "suix_getTotalSupply"

	// Extended API
	getDynamicFieldObject     SuiXMethod = "suix_getDynamicFieldObject"
	getDynamicFields          SuiXMethod = "suix_getDynamicFields"
	getOwnedObjects           SuiXMethod = "suix_getOwnedObjects"
	queryEvents               SuiXMethod = "suix_queryEvents"
	queryTransactionBlocks    SuiXMethod = "suix_queryTransactionBlocks"
	resolveNameServiceAddress SuiXMethod = "suix_resolveNameServiceAddress"
	resolveNameServiceNames   SuiXMethod = "suix_resolveNameServiceNames"
	subscribeEvent            SuiXMethod = "suix_subscribeEvent"
	subscribeTransaction      SuiXMethod = "suix_subscribeTransaction"

	// Governance Read API
	getCommitteeInfo        SuiXMethod = "suix_getCommitteeInfo"
	getLatestSuiSystemState SuiXMethod = "suix_getLatestSuiSystemState"
	getReferenceGasPrice    SuiXMethod = "suix_getReferenceGasPrice"
	getStakes               SuiXMethod = "suix_getStakes"
	getStakesByIds          SuiXMethod = "suix_getStakesByIds"
	getValidatorsApy        SuiXMethod = "suix_getValidatorsApy"

	// Move Utils
	getMoveFunctionArgTypes           SuiMethod = "sui_getMoveFunctionArgTypes"
	getNormalizedMoveFunction         SuiMethod = "sui_getNormalizedMoveFunction"
	getNormalizedMoveModule           SuiMethod = "sui_getNormalizedMoveModule"
	getNormalizedMoveModulesByPackage SuiMethod = "sui_getNormalizedMoveModulesByPackage"
	getNormalizedMoveStruct           SuiMethod = "sui_getNormalizedMoveStruct"

	// Read API
	getChainIdentifier                SuiMethod = "sui_getChainIdentifier"
	getCheckpoint                     SuiMethod = "sui_getCheckpoint"
	getCheckpoints                    SuiMethod = "sui_getCheckpoints"
	getEvents                         SuiMethod = "sui_getEvents"
	getLatestCheckpointSequenceNumber SuiMethod = "sui_getLatestCheckpointSequenceNumber"
	getObject                         SuiMethod = "sui_getObject"
	getProtocolConfig                 SuiMethod = "sui_getProtocolConfig"
	getTotalTransactionBlocks         SuiMethod = "sui_getTotalTransactionBlocks"
	getTransactionBlock               SuiMethod = "sui_getTransactionBlock"
	multiGetObjects                   SuiMethod = "sui_multiGetObjects"
	multiGetTransactionBlocks         SuiMethod = "sui_multiGetTransactionBlocks"
	tryGetPastObject                  SuiMethod = "sui_tryGetPastObject"
	tryMultiGetPastObjects            SuiMethod = "sui_tryMultiGetPastObjects"

	// Transaction Builder API
	batchTransaction     UnsafeMethod = "unsafe_batchTransaction"
	mergeCoins           UnsafeMethod = "unsafe_mergeCoins"
	moveCall             UnsafeMethod = "unsafe_moveCall"
	pay                  UnsafeMethod = "unsafe_pay"
	payAllSui            UnsafeMethod = "unsafe_payAllSui"
	paySui               UnsafeMethod = "unsafe_paySui"
	publish              UnsafeMethod = "unsafe_publish"
	requestAddStake      UnsafeMethod = "unsafe_requestAddStake"
	requestWithdrawStake UnsafeMethod = "unsafe_requestWithdrawStake"
	splitCoin            UnsafeMethod = "unsafe_splitCoin"
	splitCoinEqual       UnsafeMethod = "unsafe_splitCoinEqual"
	transferObject       UnsafeMethod = "unsafe_transferObject"
	transferSui          UnsafeMethod = "unsafe_transferSui"

	// Write API
	devInspectTransactionBlock SuiMethod = "sui_devInspectTransactionBlock"
	dryRunTransactionBlock     SuiMethod = "sui_dryRunTransactionBlock"
	executeTransactionBlock    SuiMethod = "sui_executeTransactionBlock"
)
