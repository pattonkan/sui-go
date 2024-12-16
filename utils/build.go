package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os/exec"
	"strings"

	"github.com/pattonkan/sui-go/sui"
)

type CompiledMoveModules struct {
	Modules      []*sui.Base64Data `json:"modules"`
	Dependencies []*sui.Address    `json:"dependencies"`
	Digest       []int             `json:"digest"`
}

func MoveBuild(contractPath string) (*CompiledMoveModules, error) {
	var err error
	cmd := exec.Command("sui", "move", "build", "--dump-bytecode-as-base64")
	// TODO skip to fetch latest deps if there is no internet
	// cmd := exec.Command("sui", "move", "build", "--dump-bytecode-as-base64", "--skip-fetch-latest-git-deps")
	cmd.Dir = contractPath

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("failed to execute sui cli: %w, stderr: '%s'", err, stderr.String())
	}

	var modules CompiledMoveModules
	err = json.Unmarshal(stdout.Bytes(), &modules)
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
