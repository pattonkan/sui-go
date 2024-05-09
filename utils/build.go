package utils

import (
	"encoding/json"
	"fmt"
	"log"
	"os/exec"
	"strings"

	"github.com/howjmay/sui-go/sui_types"
)

type CompiledMoveModules struct {
	Modules      []*sui_types.Base64Data `json:"modules"`
	Dependencies []*sui_types.SuiAddress `json:"dependencies"`
	Digest       []int                   `json:"digest"`
}

func MoveBuild(contractPath string) (*CompiledMoveModules, error) {
	var err error
	cmd := exec.Command("sui", "move", "build", "--dump-bytecode-as-base64")
	// TODO skip to fetch latest deps if there is no internet
	// cmd := exec.Command("sui", "move", "build", "--dump-bytecode-as-base64", "--skip-fetch-latest-git-deps")
	cmd.Dir = contractPath

	output, err := cmd.Output()
	if err != nil {
		// FIXME propagate the stderr message
		return nil, fmt.Errorf("failed to execute sui cli: %w", err)
	}

	var modules CompiledMoveModules
	err = json.Unmarshal(output, &modules)
	if err != nil {
		return nil, fmt.Errorf("failed to parse move build result: %w", err)
	}

	return &modules, nil
}

func GetGitRoot() string {
	cmd := exec.Command("git", "rev-parse", "--show-toplevel")
	output, err := cmd.Output()
	if err != nil {
		log.Fatalf("Failed to execute command: %v", err)
	}
	// Trim the newline character from the output
	return strings.TrimSpace(string(output))
}
