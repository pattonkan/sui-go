package isc_test

import (
	"testing"

	"github.com/howjmay/sui-go/isc"
)

func TestGetPublishedPackageID(t *testing.T) {
	packageID := isc.GetPublishedPackageID(isc.GetGitRoot() + "/isc/anchor_contract/iscanchor/publish_receipt.json")
	t.Log(packageID)
}
