package sui

import (
	"strings"
)

var (
	SuiCoinType = "0x2::sui::SUI"
)

func IsSameAddressString(addr1, addr2 string) bool {
	addr1 = strings.TrimPrefix(addr1, "0x")
	addr2 = strings.TrimPrefix(addr2, "0x")
	return strings.TrimLeft(addr1, "0") == strings.TrimLeft(addr2, "0")
}
