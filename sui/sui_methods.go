package sui

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
	subscribeEvent            SuiXMethod = "suix_subscribeEvent"       // TODO
	subscribeTransaction      SuiXMethod = "suix_subscribeTransaction" // TODO

	// Governance Read API
	getCommitteeInfo        SuiXMethod = "suix_getCommitteeInfo" // TODO
	getLatestSuiSystemState SuiXMethod = "suix_getLatestSuiSystemState"
	getReferenceGasPrice    SuiXMethod = "suix_getReferenceGasPrice"
	getStakes               SuiXMethod = "suix_getStakes"
	getStakesByIds          SuiXMethod = "suix_getStakesByIds"
	getValidatorsApy        SuiXMethod = "suix_getValidatorsApy"

	// Move Utils
	getMoveFunctionArgTypes           SuiMethod = "sui_getMoveFunctionArgTypes"           // TODO
	getNormalizedMoveFunction         SuiMethod = "sui_getNormalizedMoveFunction"         // TODO
	getNormalizedMoveModule           SuiMethod = "sui_getNormalizedMoveModule"           // TODO
	getNormalizedMoveModulesByPackage SuiMethod = "sui_getNormalizedMoveModulesByPackage" // TODO
	getNormalizedMoveStruct           SuiMethod = "sui_getNormalizedMoveStruct"           // TODO

	// Read API
	getChainIdentifier                SuiMethod = "sui_getChainIdentifier" // TODO
	getCheckpoint                     SuiMethod = "sui_getCheckpoint"      // TODO
	getCheckpoints                    SuiMethod = "sui_getCheckpoints"     // TODO
	getEvents                         SuiMethod = "sui_getEvents"
	getLatestCheckpointSequenceNumber SuiMethod = "sui_getLatestCheckpointSequenceNumber"
	getLoadedChildObjects             SuiMethod = "sui_getLoadedChildObjects" // TODO
	getObject                         SuiMethod = "sui_getObject"
	getProtocolConfig                 SuiMethod = "sui_getProtocolConfig" // TODO
	getTotalTransactionBlocks         SuiMethod = "sui_getTotalTransactionBlocks"
	getTransactionBlock               SuiMethod = "sui_getTransactionBlock"
	multiGetObjects                   SuiMethod = "sui_multiGetObjects"
	multiGetTransactionBlocks         SuiMethod = "sui_multiGetTransactionBlocks" // TODO
	tryGetPastObject                  SuiMethod = "sui_tryGetPastObject"
	tryMultiGetPastObjects            SuiMethod = "sui_tryMultiGetPastObjects" // TODO

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
