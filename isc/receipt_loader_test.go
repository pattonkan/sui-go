package isc_test

import (
	"testing"

	"github.com/howjmay/sui-go/isc"
)

func TestGetPublishedPackageID(t *testing.T) {
	packageID := isc.GetPublishedPackageID(isc.GetGitRoot() + "/contracts/testcoin/publish_receipt.json")
	t.Log(packageID)
}
