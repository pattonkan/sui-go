package sui_types

var (
	SuiPackageIdMoveStdlib   = MustPackageIDFromHex("0x1")
	SuiPackageIdSuiFramework = MustPackageIDFromHex("0x2")
	SuiPackageIdSuiSystem    = MustPackageIDFromHex("0x3")
	SuiPackageIdBridge       = MustPackageIDFromHex("0xb")
	SuiPackageIdDeepbook     = MustPackageIDFromHex("0xdee9")
)

var (
	SuiObjectIdSystemState        = MustObjectIDFromHex("0x5")
	SuiObjectIdClock              = MustObjectIDFromHex("0x6")
	SuiObjectIdAuthenticatorState = MustObjectIDFromHex("0x7")
	SuiObjectIdRandomnessState    = MustObjectIDFromHex("0x8")
	SuiObjectIdBridge             = MustObjectIDFromHex("0x9")
	SuiObjectIdDenyList           = MustObjectIDFromHex("0x403")
)

var (
	SuiSystemModule Identifier = "sui_system"
)

var (
	SuiSystemStateObjectSharedVersion        = SequenceNumber(1)
	SuiClockObjectSharedVersion              = SequenceNumber(1)
	SuiAuthenticatorStateObjectSharedVersion = SequenceNumber(1)
)
