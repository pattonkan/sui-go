package sui_types

var (
	SuiSystemAddress, _               = SuiAddressFromHex("0x3")
	SuiSystemPackageId                = SuiSystemAddress
	SuiSystemStateObjectID, _         = NewObjectIDFromHex("0x5")
	SuiSystemStateObjectSharedVersion = ObjectStartVersion
)
