package utils

import (
	"encoding/json"
	"fmt"
	"log"
	"os/exec"
	"strings"

	"github.com/howjmay/sui-go/sui_types/serialization"
)

type CompiledMoveModules struct {
	Modules      []*serialization.Base64Data `json:"modules"`
	Dependencies []string                    `json:"dependencies"`
	Digest       []int                       `json:"digest"`
}

func MoveBuild(path string) (*CompiledMoveModules, error) {
	var err error
	// Setup the command to be executed
	cmd := exec.Command("sui", "move", "build", "--dump-bytecode-as-base64")
	cmd.Dir = path
	fmt.Println("path ", path)

	// Run the command and capture the output
	output, err := cmd.Output()
	if err != nil {
		// return nil, err
		panic(err)
	}

	var modules CompiledMoveModules
	err = json.Unmarshal(output, &modules)
	if err != nil {
		// return nil, err
		panic(err)
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
