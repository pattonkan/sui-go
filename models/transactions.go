package models

import (
	"errors"
	"fmt"

	"github.com/howjmay/sui-go/sui_types"
	"github.com/howjmay/sui-go/sui_types/serialization"
)

type ExecuteTransactionRequestType string

const (
	TxnRequestTypeWaitForEffectsCert    ExecuteTransactionRequestType = "WaitForEffectsCert"
	TxnRequestTypeWaitForLocalExecution ExecuteTransactionRequestType = "WaitForLocalExecution"
)

type EpochId = uint64

type GasCostSummary struct {
	ComputationCost         *BigInt `json:"computationCost"`
	StorageCost             *BigInt `json:"storageCost"`
	StorageRebate           *BigInt `json:"storageRebate"`
	NonRefundableStorageFee *BigInt `json:"nonRefundableStorageFee"`
}

const (
	ExecutionStatusSuccess = "success"
	ExecutionStatusFailure = "failure"
)

type ExecutionStatus struct {
	Status string `json:"status"`
	Error  string `json:"error,omitempty"`
}

type OwnedObjectRef struct {
	Owner     serialization.TagJson[sui_types.Owner] `json:"owner"`
	Reference SuiObjectRef                           `json:"reference"`
}

type SuiTransactionBlockEffectsModifiedAtVersions struct {
	ObjectID       sui_types.ObjectID `json:"objectId"`
	SequenceNumber *BigInt            `json:"sequenceNumber"`
}

type SuiTransactionBlockEffectsV1 struct {
	/** The status of the execution */
	Status ExecutionStatus `json:"status"`
	/** The epoch when this transaction was executed */
	ExecutedEpoch *BigInt        `json:"executedEpoch"`
	GasUsed       GasCostSummary `json:"gasUsed"`
	/** The version that every modified (mutated or deleted) object had before it was modified by this transaction. **/
	ModifiedAtVersions []SuiTransactionBlockEffectsModifiedAtVersions `json:"modifiedAtVersions,omitempty"`
	/** The object references of the shared objects used in this transaction. Empty if no shared objects were used. */
	SharedObjects []SuiObjectRef `json:"sharedObjects,omitempty"`
	/** The transaction digest */
	TransactionDigest sui_types.TransactionDigest `json:"transactionDigest"`
	/** ObjectRef and owner of new objects created */
	Created []OwnedObjectRef `json:"created,omitempty"`
	/** ObjectRef and owner of mutated objects, including gas object */
	Mutated []OwnedObjectRef `json:"mutated,omitempty"`
	/**
	 * ObjectRef and owner of objects that are unwrapped in this transaction.
	 * Unwrapped objects are objects that were wrapped into other objects in the past,
	 * and just got extracted out.
	 */
	Unwrapped []OwnedObjectRef `json:"unwrapped,omitempty"`
	/** Object Refs of objects now deleted (the old refs) */
	Deleted []SuiObjectRef `json:"deleted,omitempty"`
	/** Object Refs of objects now deleted (the old refs) */
	UnwrappedThenDeleted []SuiObjectRef `json:"unwrapped_then_deleted,omitempty"`
	/** Object refs of objects now wrapped in other objects */
	Wrapped []SuiObjectRef `json:"wrapped,omitempty"`
	/**
	 * The updated gas object reference. Have a dedicated field for convenient access.
	 * It's also included in mutated.
	 */
	GasObject OwnedObjectRef `json:"gasObject"`
	/** The events emitted during execution. Note that only successful transactions emit events */
	EventsDigest *sui_types.TransactionEventsDigest `json:"eventsDigest,omitempty"`
	/** The set of transaction digests this transaction depends on */
	Dependencies []sui_types.TransactionDigest `json:"dependencies,omitempty"`
}

type SuiTransactionBlockEffects struct {
	V1 *SuiTransactionBlockEffectsV1 `json:"v1"`
}

func (t SuiTransactionBlockEffects) Tag() string {
	return "messageVersion"
}

func (t SuiTransactionBlockEffects) Content() string {
	return ""
}

func (t SuiTransactionBlockEffects) GasFee() int64 {
	if t.V1 == nil {
		return 0
	}
	fee := t.V1.GasUsed.StorageCost.Int64() -
		t.V1.GasUsed.StorageRebate.Int64() +
		t.V1.GasUsed.ComputationCost.Int64()
	return fee
}

func (t SuiTransactionBlockEffects) IsSuccess() bool {
	return t.V1.Status.Status == ExecutionStatusSuccess
}

const (
	SuiTransactionBlockKindSuiChangeEpoch             = "ChangeEpoch"
	SuiTransactionBlockKindSuiConsensusCommitPrologue = "ConsensusCommitPrologue"
	SuiTransactionBlockKindGenesis                    = "Genesis"
	SuiTransactionBlockKindProgrammableTransaction    = "ProgrammableTransaction"
)

type SuiTransactionBlockKind = serialization.TagJson[TransactionBlockKind]

type TransactionBlockKind struct {
	// A system transaction that will update epoch information on-chain.
	ChangeEpoch *SuiChangeEpoch `json:"ChangeEpoch,omitempty"`
	// A system transaction used for initializing the initial state of the chain.
	Genesis *SuiGenesisTransaction `json:"Genesis,omitempty"`
	// A system transaction marking the start of a series of transactions scheduled as part of a
	// checkpoint
	ConsensusCommitPrologue *SuiConsensusCommitPrologue `json:"ConsensusCommitPrologue,omitempty"`
	// A series of transactions where the results of one transaction can be used in future
	// transactions
	ProgrammableTransaction *SuiProgrammableTransactionBlock `json:"ProgrammableTransaction,omitempty"`
	// .. more transaction types go here
}

func (t TransactionBlockKind) Tag() string {
	return "kind"
}

func (t TransactionBlockKind) Content() string {
	return ""
}

type SuiChangeEpoch struct {
	Epoch                 *BigInt `json:"epoch"`
	StorageCharge         uint64  `json:"storage_charge"`
	ComputationCharge     uint64  `json:"computation_charge"`
	StorageRebate         uint64  `json:"storage_rebate"`
	EpochStartTimestampMs uint64  `json:"epoch_start_timestamp_ms"`
}

type SuiGenesisTransaction struct {
	Objects []sui_types.ObjectID `json:"objects"`
}

type SuiConsensusCommitPrologue struct {
	Epoch             uint64 `json:"epoch"`
	Round             uint64 `json:"round"`
	CommitTimestampMs uint64 `json:"commit_timestamp_ms"`
}

type SuiProgrammableTransactionBlock struct {
	Inputs []interface{} `json:"inputs"`
	// The transactions to be executed sequentially. A failure in any transaction will
	// result in the failure of the entire programmable transaction block.
	Commands []interface{} `json:"transactions"`
}

type SuiTransactionBlockDataV1 struct {
	Transaction SuiTransactionBlockKind `json:"transaction"`
	Sender      *sui_types.SuiAddress   `json:"sender"`
	GasData     *SuiGasData             `json:"gasData"`
}

type SuiTransactionBlockData struct {
	V1 *SuiTransactionBlockDataV1 `json:"v1,omitempty"`
}

func (t SuiTransactionBlockData) Tag() string {
	return "messageVersion"
}

func (t SuiTransactionBlockData) Content() string {
	return ""
}

type SuiTransactionBlock struct {
	Data         serialization.TagJson[SuiTransactionBlockData] `json:"data"`
	TxSignatures []string                                       `json:"txSignatures"`
}

type ObjectChange struct {
	Published *struct {
		PackageId sui_types.ObjectID     `json:"packageId"`
		Version   *BigInt                `json:"version"`
		Digest    sui_types.ObjectDigest `json:"digest"`
		Nodules   []string               `json:"nodules"`
	} `json:"published,omitempty"`
	// Transfer objects to new address / wrap in another object
	Transferred *struct {
		Sender     sui_types.SuiAddress   `json:"sender"`
		Recipient  ObjectOwner            `json:"recipient"`
		ObjectType sui_types.ObjectType   `json:"objectType"`
		ObjectID   sui_types.ObjectID     `json:"objectId"`
		Version    *BigInt                `json:"version"`
		Digest     sui_types.ObjectDigest `json:"digest"`
	} `json:"transferred,omitempty"`
	// Object mutated.
	Mutated *struct {
		Sender          sui_types.SuiAddress   `json:"sender"`
		Owner           ObjectOwner            `json:"owner"`
		ObjectType      sui_types.ObjectType   `json:"objectType"`
		ObjectID        sui_types.ObjectID     `json:"objectId"`
		Version         *BigInt                `json:"version"`
		PreviousVersion *BigInt                `json:"previousVersion"`
		Digest          sui_types.ObjectDigest `json:"digest"`
	} `json:"mutated,omitempty"`
	// Delete object j
	Deleted *struct {
		Sender     sui_types.SuiAddress `json:"sender"`
		ObjectType sui_types.ObjectType `json:"objectType"`
		ObjectID   sui_types.ObjectID   `json:"objectId"`
		Version    *BigInt              `json:"version"`
	} `json:"deleted,omitempty"`
	// Wrapped object
	Wrapped *struct {
		Sender     sui_types.SuiAddress `json:"sender"`
		ObjectType sui_types.ObjectType `json:"objectType"`
		ObjectID   sui_types.ObjectID   `json:"objectId"`
		Version    *BigInt              `json:"version"`
	} `json:"wrapped,omitempty"`
	// New object creation
	Created *struct {
		Sender     sui_types.SuiAddress   `json:"sender"`
		Owner      ObjectOwner            `json:"owner"`
		ObjectType sui_types.ObjectType   `json:"objectType"`
		ObjectID   sui_types.ObjectID     `json:"objectId"`
		Version    *BigInt                `json:"version"`
		Digest     sui_types.ObjectDigest `json:"digest"`
	} `json:"created,omitempty"`
}

func (o ObjectChange) Tag() string {
	return "type"
}

func (o ObjectChange) Content() string {
	return ""
}

type BalanceChange struct {
	Owner    ObjectOwner          `json:"owner"`
	CoinType sui_types.ObjectType `json:"coinType"`
	/* Coin balance change(positive means receive, negative means send) */
	Amount string `json:"amount"`
}

type SuiTransactionBlockResponse struct {
	Digest                  sui_types.TransactionDigest                        `json:"digest"`
	Transaction             *SuiTransactionBlock                               `json:"transaction,omitempty"`
	RawTransaction          sui_types.Base64Data                               `json:"rawTransaction,omitempty"` // enable by show_raw_input
	Effects                 *serialization.TagJson[SuiTransactionBlockEffects] `json:"effects,omitempty"`
	Events                  []*SuiEvent                                        `json:"events,omitempty"`
	TimestampMs             *BigInt                                            `json:"timestampMs,omitempty"`
	Checkpoint              *BigInt                                            `json:"checkpoint,omitempty"`
	ConfirmedLocalExecution *bool                                              `json:"confirmedLocalExecution,omitempty"`
	ObjectChanges           []serialization.TagJson[ObjectChange]              `json:"objectChanges,omitempty"`
	BalanceChanges          []BalanceChange                                    `json:"balanceChanges,omitempty"`
	Errors                  []string                                           `json:"errors,omitempty"` // Errors that occurred in fetching/serializing the transaction.
	// FIXME datatype may be wrong
	RawEffects []string `json:"rawEffects,omitempty"` // enable by show_raw_effects
}

// requires to set 'SuiTransactionBlockResponseOptions.ShowObjectChanges' to true
func (r *SuiTransactionBlockResponse) GetPublishedPackageID() (*sui_types.PackageID, error) {
	if r.ObjectChanges == nil {
		return nil, errors.New("no 'SuiTransactionBlockResponse.ObjectChanges' object")
	}
	var packageID sui_types.PackageID
	for _, change := range r.ObjectChanges {
		if change.Data.Published != nil {
			packageID = change.Data.Published.PackageId
			return &packageID, nil
		}
	}
	return nil, fmt.Errorf("not found")
}

func (r *SuiTransactionBlockResponse) GetCreatedObjectInfo(module string, objName string) (*sui_types.ObjectID, sui_types.ObjectType, error) {
	if r.ObjectChanges == nil {
		return nil, "", fmt.Errorf("no ObjectChanges")
	}
	for _, change := range r.ObjectChanges {
		if change.Data.Created != nil {
			// some possible examples
			// * 0x2::coin::TreasuryCap<0x14c12b454ac6996024342312769e00bb98c70ad2f3546a40f62516c83aa0f0d4::testcoin::TESTCOIN>
			// * 0x14c12b454ac6996024342312769e00bb98c70ad2f3546a40f62516c83aa0f0d4::anchor::Anchor
			resource, err := NewResourceType(change.Data.Created.ObjectType)
			if err != nil {
				return nil, "", fmt.Errorf("invalid resource target string")
			}
			if resource.Module == module && resource.ObjectName == objName {
				return &change.Data.Created.ObjectID, change.Data.Created.ObjectType, nil
			}
			for ; resource != nil; resource = resource.SubType {
				if resource.Module == module && resource.ObjectName == objName {
					return &change.Data.Created.ObjectID, change.Data.Created.ObjectType, nil
				}
			}
		}
	}
	return nil, "", fmt.Errorf("not found")
}

type ReturnValueType interface{}
type MutableReferenceOutputType interface{}
type ExecutionResultType struct {
	MutableReferenceOutputs []MutableReferenceOutputType `json:"mutableReferenceOutputs,omitempty"`
	ReturnValues            []ReturnValueType            `json:"returnValues,omitempty"`
}

type DevInspectResults struct {
	Effects serialization.TagJson[SuiTransactionBlockEffects] `json:"effects"`
	Events  []SuiEvent                                        `json:"events"`
	Results []ExecutionResultType                             `json:"results,omitempty"`
	Error   *string                                           `json:"error,omitempty"`
}

type TransactionFilter struct {
	// Query by checkpoint.
	Checkpoint *BigInt `json:"Checkpoint,omitempty"`
	// Query by move function.
	MoveFunction *TransactionFilterMoveFunction `json:"MoveFunction,omitempty"`
	// Query by input object.
	InputObject *sui_types.ObjectID `json:"InputObject,omitempty"`
	// Query by changed object, including created, mutated and unwrapped objects.
	ChangedObject *sui_types.ObjectID `json:"ChangedObject,omitempty"`
	// Query by sender address.
	FromAddress *sui_types.SuiAddress `json:"FromAddress,omitempty"`
	// Query by recipient address.
	ToAddress *sui_types.SuiAddress `json:"ToAddress,omitempty"`
	// Query by sender and recipient address.
	FromAndToAddress *TransactionFilterFromAndToAddress `json:"FromAndToAddress,omitempty"`
	// Query txs that have a given address as sender or recipient.
	FromOrToAddress *sui_types.SuiAddress `json:"FromOrToAddress,omitempty"`
	// Query by transaction kind
	TransactionKind *string `json:"TransactionKind,omitempty"`
	// Query transactions of any given kind in the input.
	TransactionKindIn []string `json:"TransactionKindIn,omitempty"`
}

type TransactionFilterMoveFunction struct {
	Package  *sui_types.ObjectID   `json:"package"`
	Module   *sui_types.Identifier `json:"module,omitempty"`
	Function *sui_types.Identifier `json:"function,omitempty"`
}

type TransactionFilterFromAndToAddress struct {
	From *sui_types.SuiAddress `json:"from"`
	To   *sui_types.SuiAddress `json:"to"`
}

type SuiTransactionBlockResponseOptions struct {
	// Whether to show transaction input data. Default to be False
	ShowInput bool `json:"showInput,omitempty"`
	// Whether to show bcs-encoded transaction input data
	ShowRawInput bool `json:"showRawInput,omitempty"`
	// Whether to show transaction effects. Default to be False
	ShowEffects bool `json:"showEffects,omitempty"`
	// Whether to show transaction events. Default to be False
	ShowEvents bool `json:"showEvents,omitempty"`
	// Whether to show object_changes. Default to be False
	ShowObjectChanges bool `json:"showObjectChanges,omitempty"`
	// Whether to show balance_changes. Default to be False
	ShowBalanceChanges bool `json:"showBalanceChanges,omitempty"`
	// Whether to show raw transaction effects. Default to be False
	ShowRawEffects bool `json:"showRawEffects,omitempty"`
}

type SuiTransactionBlockResponseQuery struct {
	Filter  *TransactionFilter                  `json:"filter,omitempty"`
	Options *SuiTransactionBlockResponseOptions `json:"options,omitempty"`
}

type TransactionBlocksPage = Page[SuiTransactionBlockResponse, sui_types.TransactionDigest]

type DryRunTransactionBlockResponse struct {
	Effects        serialization.TagJson[SuiTransactionBlockEffects] `json:"effects"`
	Events         []SuiEvent                                        `json:"events"`
	ObjectChanges  []serialization.TagJson[ObjectChange]             `json:"objectChanges"`
	BalanceChanges []BalanceChange                                   `json:"balanceChanges"`
	Input          serialization.TagJson[SuiTransactionBlockData]    `json:"input"`
}
