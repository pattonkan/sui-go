package sui

import (
	"fmt"
	"strings"

	"github.com/howjmay/sui-go/models"
	"github.com/howjmay/sui-go/sui_types"
)

// requires `ShowObjectChanges: true`
func GetCreatedObjectIdAndType(txRes *models.SuiTransactionBlockResponse, moduleName string, objectName string) (*sui_types.ObjectID, string, error) {
	if txRes.ObjectChanges == nil {
		return nil, "", fmt.Errorf("no ObjectChanges")
	}
	for _, change := range txRes.ObjectChanges {
		if change.Data.Created != nil {
			// FIXME error-prone, we need to parse the object type to check
			// some possible examples
			// * 0x2::coin::TreasuryCap<0x14c12b454ac6996024342312769e00bb98c70ad2f3546a40f62516c83aa0f0d4::testcoin::TESTCOIN>
			// * 0x14c12b454ac6996024342312769e00bb98c70ad2f3546a40f62516c83aa0f0d4::anchor::Anchor
			if strings.Contains(change.Data.Created.ObjectType, fmt.Sprintf("%s::%s", moduleName, objectName)) {
				return &change.Data.Created.ObjectID, change.Data.Created.ObjectType, nil
			}
		}
	}
	return nil, "", fmt.Errorf("not found")
}
