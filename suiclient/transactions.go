package suiclient

import (
	"errors"
	"fmt"

	"github.com/pattonkan/sui-go/sui"
)

type RPCTransactionRequestParams struct {
	TransferObjectRequestParams *TransferObjectParams `json:"transferObjectRequestParams,omitempty"`
	MoveCallRequestParams       *MoveCallParams       `json:"moveCallRequestParams,omitempty"`
}

type TransferObjectParams struct {
	Recipient *sui.Address  `json:"recipient,omitempty"`
	ObjectId  *sui.ObjectId `json:"objectId,omitempty"`
}

type MoveCallParams struct {
	PackageObjectId *sui.ObjectId  `json:"packageObjectId,omitempty"`
	Module          sui.Identifier `json:"module,omitempty"`
	Function        sui.Identifier `json:"function,omitempty"`
	TypeArguments   []SuiTypeTag   `json:"typeArguments,omitempty"`
	// MoveCall arguments in JSON format
	Arguments []SuiJsonValue `json:"arguments,omitempty"`
}

type ExecuteTransactionRequestType string

const (
	TxnRequestTypeWaitForEffectsCert    ExecuteTransactionRequestType = "WaitForEffectsCert"
	TxnRequestTypeWaitForLocalExecution ExecuteTransactionRequestType = "WaitForLocalExecution"
)

// type EpochId = uint64

type GasCostSummary struct {
	ComputationCost         *sui.BigInt `json:"computationCost"`
	StorageCost             *sui.BigInt `json:"storageCost"`
	StorageRebate           *sui.BigInt `json:"storageRebate"`
	NonRefundableStorageFee *sui.BigInt `json:"nonRefundableStorageFee"`
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
	Owner     WrapperTaggedJson[sui.Owner] `json:"owner"`
	Reference SuiObjectRef                 `json:"reference"`
}

type SuiTransactionBlockEffectsModifiedAtVersions struct {
	ObjectId       sui.ObjectId `json:"objectId"`
	SequenceNumber *sui.BigInt  `json:"sequenceNumber"`
}

type SuiTransactionBlockEffectsV1 struct {
	/** The status of the execution */
	Status ExecutionStatus `json:"status"`
	/** The epoch when this transaction was executed */
	ExecutedEpoch *sui.BigInt    `json:"executedEpoch"`
	GasUsed       GasCostSummary `json:"gasUsed"`
	/** The version that every modified (mutated or deleted) object had before it was modified by this transaction. **/
	ModifiedAtVersions []SuiTransactionBlockEffectsModifiedAtVersions `json:"modifiedAtVersions,omitempty"`
	/** The object references of the shared objects used in this transaction. Empty if no shared objects were used. */
	SharedObjects []SuiObjectRef `json:"sharedObjects,omitempty"`
	/** The transaction digest */
	TransactionDigest sui.TransactionDigest `json:"transactionDigest"`
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
	EventsDigest *sui.TransactionEventsDigest `json:"eventsDigest,omitempty"`
	/** The set of transaction digests this transaction depends on */
	Dependencies []sui.TransactionDigest `json:"dependencies,omitempty"`
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

type SuiTransactionBlockKind = WrapperTaggedJson[TransactionBlockKind]

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
	Epoch                 *sui.BigInt `json:"epoch"`
	StorageCharge         uint64      `json:"storage_charge"`
	ComputationCharge     uint64      `json:"computation_charge"`
	StorageRebate         uint64      `json:"storage_rebate"`
	EpochStartTimestampMs uint64      `json:"epoch_start_timestamp_ms"`
}

type SuiGenesisTransaction struct {
	Objects []sui.ObjectId `json:"objects"`
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
	Sender      *sui.Address            `json:"sender"`
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
	Data         WrapperTaggedJson[SuiTransactionBlockData] `json:"data"`
	TxSignatures []string                                   `json:"txSignatures"`
}

type ObjectChange struct {
	Published *struct {
		PackageId sui.ObjectId     `json:"packageId"`
		Version   *sui.BigInt      `json:"version"`
		Digest    sui.ObjectDigest `json:"digest"`
		Nodules   []string         `json:"nodules"`
	} `json:"published,omitempty"`
	// Transfer objects to new address / wrap in another object
	Transferred *struct {
		Sender     sui.Address      `json:"sender"`
		Recipient  ObjectOwner      `json:"recipient"`
		ObjectType sui.ObjectType   `json:"objectType"`
		ObjectId   sui.ObjectId     `json:"objectId"`
		Version    *sui.BigInt      `json:"version"`
		Digest     sui.ObjectDigest `json:"digest"`
	} `json:"transferred,omitempty"`
	// Object mutated.
	Mutated *struct {
		Sender          sui.Address      `json:"sender"`
		Owner           ObjectOwner      `json:"owner"`
		ObjectType      sui.ObjectType   `json:"objectType"`
		ObjectId        sui.ObjectId     `json:"objectId"`
		Version         *sui.BigInt      `json:"version"`
		PreviousVersion *sui.BigInt      `json:"previousVersion"`
		Digest          sui.ObjectDigest `json:"digest"`
	} `json:"mutated,omitempty"`
	// Delete object j
	Deleted *struct {
		Sender     sui.Address    `json:"sender"`
		ObjectType sui.ObjectType `json:"objectType"`
		ObjectId   sui.ObjectId   `json:"objectId"`
		Version    *sui.BigInt    `json:"version"`
	} `json:"deleted,omitempty"`
	// Wrapped object
	Wrapped *struct {
		Sender     sui.Address    `json:"sender"`
		ObjectType sui.ObjectType `json:"objectType"`
		ObjectId   sui.ObjectId   `json:"objectId"`
		Version    *sui.BigInt    `json:"version"`
	} `json:"wrapped,omitempty"`
	// New object creation
	Created *struct {
		Sender     sui.Address      `json:"sender"`
		Owner      ObjectOwner      `json:"owner"`
		ObjectType sui.ObjectType   `json:"objectType"`
		ObjectId   sui.ObjectId     `json:"objectId"`
		Version    *sui.BigInt      `json:"version"`
		Digest     sui.ObjectDigest `json:"digest"`
	} `json:"created,omitempty"`
}

func (o ObjectChange) Tag() string {
	return "type"
}

func (o ObjectChange) Content() string {
	return ""
}

type BalanceChange struct {
	Owner    ObjectOwner    `json:"owner"`
	CoinType sui.ObjectType `json:"coinType"`
	/* Coin balance change(positive means receive, negative means send) */
	Amount string `json:"amount"`
}

type SuiTransactionBlockResponse struct {
	Digest                  sui.TransactionDigest                          `json:"digest"`
	Transaction             *SuiTransactionBlock                           `json:"transaction,omitempty"`
	RawTransaction          sui.Base64Data                                 `json:"rawTransaction,omitempty"` // enable by show_raw_input
	Effects                 *WrapperTaggedJson[SuiTransactionBlockEffects] `json:"effects,omitempty"`
	Events                  []*Event                                       `json:"events,omitempty"`
	TimestampMs             *sui.BigInt                                    `json:"timestampMs,omitempty"`
	Checkpoint              *sui.BigInt                                    `json:"checkpoint,omitempty"`
	ConfirmedLocalExecution *bool                                          `json:"confirmedLocalExecution,omitempty"`
	ObjectChanges           []WrapperTaggedJson[ObjectChange]              `json:"objectChanges,omitempty"`
	BalanceChanges          []BalanceChange                                `json:"balanceChanges,omitempty"`
	Errors                  []string                                       `json:"errors,omitempty"` // Errors that occurred in fetching/serializing the transaction.
	// FIXME datatype may be wrong
	RawEffects []byte `json:"rawEffects,omitempty"` // enable by show_raw_effects
}

// requires to set 'SuiTransactionBlockResponseOptions.ShowObjectChanges' to true
func (r *SuiTransactionBlockResponse) GetPublishedPackageId() (*sui.PackageId, error) {
	if r.ObjectChanges == nil {
		return nil, errors.New("no 'SuiTransactionBlockResponse.ObjectChanges' object")
	}
	var packageId sui.PackageId
	for _, change := range r.ObjectChanges {
		if change.Data.Published != nil {
			packageId = change.Data.Published.PackageId
			return &packageId, nil
		}
	}
	return nil, fmt.Errorf("not found")
}

func (r *SuiTransactionBlockResponse) GetCreatedObjectInfo(module string, objName string) (*sui.ObjectId, sui.ObjectType, error) {
	if r.ObjectChanges == nil {
		return nil, "", fmt.Errorf("no ObjectChanges")
	}
	for _, change := range r.ObjectChanges {
		if change.Data.Created != nil {
			// some possible examples
			// * 0x2::coin::TreasuryCap<0x14c12b454ac6996024342312769e00bb98c70ad2f3546a40f62516c83aa0f0d4::testcoin::TESTCOIN>
			// * 0x14c12b454ac6996024342312769e00bb98c70ad2f3546a40f62516c83aa0f0d4::anchor::Anchor
			resource, err := sui.NewResourceType(change.Data.Created.ObjectType)
			if err != nil {
				return nil, "", fmt.Errorf("invalid resource target string")
			}
			if resource.Module == module && resource.ObjectName == objName {
				return &change.Data.Created.ObjectId, change.Data.Created.ObjectType, nil
			}
			for ; resource != nil; resource = resource.SubType {
				if resource.Module == module && resource.ObjectName == objName {
					return &change.Data.Created.ObjectId, change.Data.Created.ObjectType, nil
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

type TransactionFilter struct {
	// Query by checkpoint.
	Checkpoint *sui.BigInt `json:"Checkpoint,omitempty"`
	// Query by move function.
	MoveFunction *TransactionFilterMoveFunction `json:"MoveFunction,omitempty"`
	// Query by input object.
	InputObject *sui.ObjectId `json:"InputObject,omitempty"`
	// Query by changed object, including created, mutated and unwrapped objects.
	ChangedObject *sui.ObjectId `json:"ChangedObject,omitempty"`
	// Query by sender address.
	FromAddress *sui.Address `json:"FromAddress,omitempty"`
	// Query by recipient address.
	ToAddress *sui.Address `json:"ToAddress,omitempty"`
	// Query by sender and recipient address.
	FromAndToAddress *TransactionFilterFromAndToAddress `json:"FromAndToAddress,omitempty"`
	// Query txs that have a given address as sender or recipient.
	FromOrToAddress *sui.Address `json:"FromOrToAddress,omitempty"`
	// Query by transaction kind
	TransactionKind *string `json:"TransactionKind,omitempty"`
	// Query transactions of any given kind in the input.
	TransactionKindIn []string `json:"TransactionKindIn,omitempty"`
}

type TransactionFilterMoveFunction struct {
	Package  *sui.ObjectId   `json:"package"`
	Module   *sui.Identifier `json:"module,omitempty"`
	Function *sui.Identifier `json:"function,omitempty"`
}

type TransactionFilterFromAndToAddress struct {
	From *sui.Address `json:"from"`
	To   *sui.Address `json:"to"`
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

type TransactionBlockResponseQuery struct {
	Filter  *TransactionFilter                  `json:"filter,omitempty"`
	Options *SuiTransactionBlockResponseOptions `json:"options,omitempty"`
}

type TransactionBlocksPage = Page[SuiTransactionBlockResponse, sui.TransactionDigest]

type DryRunTransactionBlockResponse struct {
	Effects        WrapperTaggedJson[SuiTransactionBlockEffects] `json:"effects"`
	Events         []Event                                       `json:"events"`
	ObjectChanges  []WrapperTaggedJson[ObjectChange]             `json:"objectChanges"`
	BalanceChanges []BalanceChange                               `json:"balanceChanges"`
	Input          WrapperTaggedJson[SuiTransactionBlockData]    `json:"input"`
}
