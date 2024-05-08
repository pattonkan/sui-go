package utils_test

import (
	"fmt"
	"testing"

	"github.com/howjmay/sui-go/utils"

	"github.com/stretchr/testify/require"
)

func TestMoveBuild(t *testing.T) {
	t.Skip("FIXME install sui for ci to test")
	// FIXME add a testing contract for the localnet
	modules, err := utils.MoveBuild(utils.GetGitRoot() + "/isc/contracts/isc")
	require.NoError(t, err)
	fmt.Println("modules", modules)
}
