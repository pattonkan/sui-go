package lib

import (
	"encoding/json"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/howjmay/sui-go/sui"
)

type PublishReceipt struct {
	Digest      string `json:"digest"`
	Transaction struct {
		Data struct {
			MessageVersion string `json:"messageVersion"`
			Transaction    struct {
				Kind   string `json:"kind"`
				Inputs []struct {
					Type      string `json:"type"`
					ValueType string `json:"valueType"`
					Value     string `json:"value"`
				} `json:"inputs"`
				Transactions []struct {
					Publish         []string      `json:"Publish,omitempty"`
					TransferObjects []interface{} `json:"TransferObjects,omitempty"`
				} `json:"transactions"`
			} `json:"transaction"`
			Sender  string `json:"sender"`
			GasData struct {
				Payment []struct {
					ObjectId string `json:"objectId"`
					Version  int    `json:"version"`
					Digest   string `json:"digest"`
				} `json:"payment"`
				Owner  string `json:"owner"`
				Price  string `json:"price"`
				Budget string `json:"budget"`
			} `json:"gasData"`
		} `json:"data"`
		TxSignatures []string `json:"txSignatures"`
	} `json:"transaction"`
	Effects struct {
		MessageVersion string `json:"messageVersion"`
		Status         struct {
			Status string `json:"status"`
		} `json:"status"`
		ExecutedEpoch string `json:"executedEpoch"`
		GasUsed       struct {
			ComputationCost         string `json:"computationCost"`
			StorageCost             string `json:"storageCost"`
			StorageRebate           string `json:"storageRebate"`
			NonRefundableStorageFee string `json:"nonRefundableStorageFee"`
		} `json:"gasUsed"`
		ModifiedAtVersions []struct {
			ObjectId       string `json:"objectId"`
			SequenceNumber string `json:"sequenceNumber"`
		} `json:"modifiedAtVersions"`
		TransactionDigest string `json:"transactionDigest"`
		Created           []struct {
			Owner struct {
				AddressOwner string `json:"AddressOwner"`
			} `json:"owner"`
			Reference struct {
				ObjectId string `json:"objectId"`
				Version  int    `json:"version"`
				Digest   string `json:"digest"`
			} `json:"reference"`
		} `json:"created"`
		Mutated []struct {
			Owner struct {
				AddressOwner string `json:"AddressOwner"`
			} `json:"owner"`
			Reference struct {
				ObjectId string `json:"objectId"`
				Version  int    `json:"version"`
				Digest   string `json:"digest"`
			} `json:"reference"`
		} `json:"mutated"`
		GasObject struct {
			Owner struct {
				AddressOwner string `json:"AddressOwner"`
			} `json:"owner"`
			Reference struct {
				ObjectId string `json:"objectId"`
				Version  int    `json:"version"`
				Digest   string `json:"digest"`
			} `json:"reference"`
		} `json:"gasObject"`
		Dependencies []string `json:"dependencies"`
	} `json:"effects"`
	Events        []interface{} `json:"events"`
	ObjectChanges []struct {
		Type   string `json:"type"`
		Sender string `json:"sender,omitempty"`
		Owner  struct {
			AddressOwner string `json:"AddressOwner"`
		} `json:"owner,omitempty"`
		ObjectType      string   `json:"objectType,omitempty"`
		ObjectId        string   `json:"objectId,omitempty"`
		Version         string   `json:"version"`
		PreviousVersion string   `json:"previousVersion,omitempty"`
		Digest          string   `json:"digest"`
		PackageId       string   `json:"packageId,omitempty"`
		Modules         []string `json:"modules,omitempty"`
	} `json:"objectChanges"`
	BalanceChanges []struct {
		Owner struct {
			AddressOwner string `json:"AddressOwner"`
		} `json:"owner"`
		CoinType string `json:"coinType"`
		Amount   string `json:"amount"`
	} `json:"balanceChanges"`
	ConfirmedLocalExecution bool `json:"confirmedLocalExecution"`
}

func GetPublishedPackageId(receiptJson string) *sui.PackageId {
	filePath := "path/to/your/file.json" + getGitRoot()
	jsonData, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatalf("failed to read file: %v", err)
	}

	var receipt PublishReceipt
	if err := json.Unmarshal(jsonData, &receipt); err != nil {
		log.Fatalf("error unmarshaling json: %v", err)
	}

	var packageId string
	for _, change := range receipt.ObjectChanges {
		if change.Type == "published" {
			packageId = change.PackageId
		}
	}
	suiPackageId, err := sui.PackageIdFromHex(packageId)
	if err != nil {
		log.Fatalf("failed to decode hex package Id: %v", err)
	}
	return suiPackageId
}

func getGitRoot() string {
	cmd := exec.Command("git", "rev-parse", "--show-toplevel")
	output, err := cmd.Output()
	if err != nil {
		log.Fatalf("Failed to execute command: %v", err)
	}
	// Trim the newline character from the output
	return strings.TrimSpace(string(output))
}
