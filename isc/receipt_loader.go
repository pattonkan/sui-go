package isc

import (
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/howjmay/sui-go/sui_types"
	"github.com/tidwall/gjson"
)

func GetPublishedPackageID(jsonPath string) *sui_types.PackageID {
	jsonData, err := os.ReadFile(jsonPath)
	if err != nil {
		log.Fatalf("failed to read file: %v", err)
	}

	var packageID string
	objectChanges := gjson.Get(string(jsonData), "objectChanges").Array()

	for _, change := range objectChanges {
		if change.Get("type").String() == "published" {
			packageID = change.Get("packageId").String()
		}
	}
	suiPackageID, err := sui_types.PackageIDFromHex(packageID)
	if err != nil {
		log.Fatalf("failed to decode hex package ID: %v", err)
	}
	return suiPackageID
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
