package sui_types

var (
	SuiSystemAddress                  = MustSuiAddressFromHex("0x3")
	SuiSystemPackageId                = SuiSystemAddress
	SuiSystemModuleName               = "sui_system"
	SuiSystemStateObjectID            = MustObjectIDFromHex("0x5")
	SuiSystemStateObjectSharedVersion = ObjectStartVersion
)
