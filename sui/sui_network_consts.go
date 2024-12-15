package sui

var (
	SuiPackageIdMoveStdlib   = MustPackageIdFromHex("0x1")
	SuiPackageIdSuiFramework = MustPackageIdFromHex("0x2")
	SuiPackageIdSuiSystem    = MustPackageIdFromHex("0x3")
	SuiPackageIdBridge       = MustPackageIdFromHex("0xb")
	SuiPackageIdDeepbook     = MustPackageIdFromHex("0xdee9")
)

var (
	SuiObjectIdSystemState        = MustObjectIdFromHex("0x5")
	SuiObjectIdClock              = MustObjectIdFromHex("0x6")
	SuiObjectIdAuthenticatorState = MustObjectIdFromHex("0x7")
	SuiObjectIdRandomnessState    = MustObjectIdFromHex("0x8")
	SuiObjectIdBridge             = MustObjectIdFromHex("0x9")
	SuiObjectIdDenyList           = MustObjectIdFromHex("0x403")
)

var (
	SuiSystemModuleName Identifier = "sui_system"
)

var (
	SuiSystemStateObjectSharedVersion        = SequenceNumber(1)
	SuiClockObjectSharedVersion              = SequenceNumber(1)
	SuiAuthenticatorStateObjectSharedVersion = SequenceNumber(1)
)
