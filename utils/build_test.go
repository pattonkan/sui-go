package utils_test

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/howjmay/sui-go/utils"

	"github.com/stretchr/testify/require"
)

func TestMoveBuild(t *testing.T) {
	t.Skip("FIXME install sui for ci to test")
	modules, err := utils.MoveBuild(utils.GetGitRoot() + "/contracts/testcoin/")
	require.NoError(t, err)

	jsonData, err := os.ReadFile(utils.GetGitRoot() + "/contracts/testcoin/contract_base64.json")
	require.NoError(t, err)
	var expectedModules utils.CompiledMoveModules
	err = json.Unmarshal(jsonData, &expectedModules)
	require.NoError(t, err)
	require.Equal(t, &expectedModules, modules)
}
