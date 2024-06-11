package models

import (
	"strconv"
	"strings"

	"github.com/howjmay/sui-go/sui_types"
	"github.com/howjmay/sui-go/sui_types/serialization"
)

type SuiObjectRef struct {
	/** Base64 string representing the object digest */
	Digest sui_types.TransactionDigest `json:"digest"`
	/** Hex code as string representing the object id */
	ObjectID string `json:"objectId"` // TODO refactor datatype to ObjectID
	/** Object version */
	Version sui_types.SequenceNumber `json:"version"`
}

type SuiGasData struct {
	Payment []SuiObjectRef `json:"payment"`
	/** Gas Object's owner */
	Owner  string                `json:"owner"`
	Price  SafeSuiBigInt[uint64] `json:"price"`
	Budget SafeSuiBigInt[uint64] `json:"budget"`
}

type SuiParsedData struct {
	MoveObject *SuiParsedMoveObject `json:"moveObject,omitempty"`
	Package    *SuiMovePackage      `json:"package,omitempty"`
}

func (p SuiParsedData) Tag() string {
	return "dataType"
}

func (p SuiParsedData) Content() string {
	return ""
}

type SuiMovePackage struct {
	Disassembled map[string]interface{} `json:"disassembled"`
}

type SuiParsedMoveObject struct {
	Type              string `json:"type"`
	HasPublicTransfer bool   `json:"hasPublicTransfer"`
	// maybe use reflect to solve this
	Fields any `json:"fields"`
}

type SuiRawData struct {
	MoveObject *SuiRawMoveObject  `json:"moveObject,omitempty"`
	Package    *SuiRawMovePackage `json:"package,omitempty"`
}

func (r SuiRawData) Tag() string {
	return "dataType"
}

func (r SuiRawData) Content() string {
	return ""
}

type SuiRawMoveObject struct {
	Type              string                   `json:"type"`
	HasPublicTransfer bool                     `json:"hasPublicTransfer"`
	Version           sui_types.SequenceNumber `json:"version"`
	BcsBytes          sui_types.Base64Data     `json:"bcsBytes"`
}

type SuiRawMovePackage struct {
	Id              sui_types.ObjectID              `json:"id"`
	Version         sui_types.SequenceNumber        `json:"version"`
	ModuleMap       map[string]sui_types.Base64Data `json:"moduleMap"`
	TypeOriginTable []TypeOrigin                    `json:"typeOriginTable"`
	LinkageTable    map[string]UpgradeInfo
}

type UpgradeInfo struct {
	UpgradedId      sui_types.ObjectID
	UpgradedVersion sui_types.SequenceNumber
}

type TypeOrigin struct {
	ModuleName string             `json:"moduleName"`
	StructName string             `json:"structName"`
	Package    sui_types.ObjectID `json:"package"`
}

type SuiObjectData struct {
	ObjectID *sui_types.ObjectID                     `json:"objectId"`
	Version  SafeSuiBigInt[sui_types.SequenceNumber] `json:"version"`
	Digest   *sui_types.ObjectDigest                 `json:"digest"`
	/**
	 * Type of the object, default to be undefined unless SuiObjectDataOptions.showType is set to true
	 */
	Type *string `json:"type,omitempty"`
	/**
	 * Move object content or package content, default to be undefined unless SuiObjectDataOptions.showContent is set to true
	 */
	Content *serialization.TagJson[SuiParsedData] `json:"content,omitempty"`
	/**
	 * Move object content or package content in BCS bytes, default to be undefined unless SuiObjectDataOptions.showBcs is set to true
	 */
	Bcs *serialization.TagJson[SuiRawData] `json:"bcs,omitempty"`
	/**
	 * The owner of this object. Default to be undefined unless SuiObjectDataOptions.showOwner is set to true
	 */
	Owner *ObjectOwner `json:"owner,omitempty"`
	/**
	 * The digest of the transaction that created or last mutated this object.
	 * Default to be undefined unless SuiObjectDataOptions.showPreviousTransaction is set to true
	 */
	PreviousTransaction *sui_types.TransactionDigest `json:"previousTransaction,omitempty"`
	/**
	 * The amount of SUI we would rebate if this object gets deleted.
	 * This number is re-calculated each time the object is mutated based on
	 * the present storage gas price.
	 * Default to be undefined unless SuiObjectDataOptions.showStorageRebate is set to true
	 */
	StorageRebate *SafeSuiBigInt[uint64] `json:"storageRebate,omitempty"`
	/**
	 * Display metadata for this object, default to be undefined unless SuiObjectDataOptions.showDisplay is set to true
	 * This can also be None if the struct type does not have Display defined
	 * See more details in https://forums.sui.io/t/nft-object-display-proposal/4872
	 */
	Display interface{} `json:"display,omitempty"`
}

func (data *SuiObjectData) Ref() sui_types.ObjectRef {
	return sui_types.ObjectRef{
		ObjectID: data.ObjectID,
		Version:  data.Version.data,
		Digest:   data.Digest,
	}
}

type SuiObjectDataOptions struct {
	/* Whether to fetch the object type, default to be false */
	ShowType bool `json:"showType,omitempty"`
	/* Whether to fetch the object content, default to be false */
	ShowContent bool `json:"showContent,omitempty"`
	/* Whether to fetch the object content in BCS bytes, default to be false */
	ShowBcs bool `json:"showBcs,omitempty"`
	/* Whether to fetch the object owner, default to be false */
	ShowOwner bool `json:"showOwner,omitempty"`
	/* Whether to fetch the previous transaction digest, default to be false */
	ShowPreviousTransaction bool `json:"showPreviousTransaction,omitempty"`
	/* Whether to fetch the storage rebate, default to be false */
	ShowStorageRebate bool `json:"showStorageRebate,omitempty"`
	/* Whether to fetch the display metadata, default to be false */
	ShowDisplay bool `json:"showDisplay,omitempty"`
}

type SuiObjectResponseError struct {
	NotExists *struct {
		ObjectID sui_types.ObjectID `json:"object_id"`
	} `json:"notExists,omitempty"`
	Deleted *struct {
		ObjectID sui_types.ObjectID       `json:"object_id"`
		Version  sui_types.SequenceNumber `json:"version"`
		Digest   sui_types.ObjectDigest   `json:"digest"`
	} `json:"deleted,omitempty"`
	UnKnown      *struct{} `json:"unKnown"`
	DisplayError *struct {
		Error string `json:"error"`
	} `json:"displayError"`
}

func (e SuiObjectResponseError) Tag() string {
	return "code"
}

func (e SuiObjectResponseError) Content() string {
	return ""
}

type SuiObjectResponse struct {
	Data  *SuiObjectData                                 `json:"data,omitempty"`
	Error *serialization.TagJson[SuiObjectResponseError] `json:"error,omitempty"`
}

type CheckpointSequenceNumber = uint64

type CheckpointedObjectID struct {
	ObjectID     sui_types.ObjectID                       `json:"objectId"`
	AtCheckpoint *SafeSuiBigInt[CheckpointSequenceNumber] `json:"atCheckpoint"`
}

type ObjectsPage = Page[SuiObjectResponse, sui_types.ObjectID]

// TODO need use Enum
type SuiObjectDataFilter struct {
	Package    *sui_types.ObjectID `json:"Package,omitempty"`
	MoveModule *MoveModule         `json:"MoveModule,omitempty"`
	StructType string              `json:"StructType,omitempty"`
}

type SuiObjectResponseQuery struct {
	Filter  *SuiObjectDataFilter  `json:"filter,omitempty"`
	Options *SuiObjectDataOptions `json:"options,omitempty"`
}

type SuiPastObjectResponse = serialization.TagJson[SuiPastObject]

// TODO need test VersionNotFound
type SuiPastObject struct {
	// The object exists and is found with this version
	VersionFound *SuiObjectData `json:"VersionFound,omitempty"`
	// The object does not exist
	ObjectNotExists *sui_types.ObjectID `json:"ObjectNotExists,omitempty"`
	// The object is found to be deleted with this version
	ObjectDeleted *SuiObjectRef `json:"ObjectDeleted,omitempty"`
	// The object exists but not found with this version
	VersionNotFound *VersionNotFoundData `json:"VersionNotFound,omitempty"`
	// The asked object version is higher than the latest
	VersionTooHigh *struct {
		ObjectID      sui_types.ObjectID       `json:"object_id"`
		AskedVersion  sui_types.SequenceNumber `json:"asked_version"`
		LatestVersion sui_types.SequenceNumber `json:"latest_version"`
	} `json:"VersionTooHigh,omitempty"`
}

type VersionNotFoundData struct {
	ObjectID       *sui_types.ObjectID
	SequenceNumber sui_types.SequenceNumber
}

func (c *VersionNotFoundData) UnmarshalJSON(data []byte) error {
	var err error
	input := data[1 : len(data)-2]
	elts := strings.Split(string(input), ",")
	c.ObjectID, err = sui_types.ObjectIDFromHex(elts[0][1 : len(elts[0])-2])
	if err != nil {
		return err
	}
	seq, err := strconv.ParseUint(elts[1], 10, 64)
	if err != nil {
		return err
	}
	c.SequenceNumber = seq
	return nil
}

func (s SuiPastObject) Tag() string {
	return "status"
}

func (s SuiPastObject) Content() string {
	return "details"
}

type SuiGetPastObjectRequest struct {
	ObjectId *sui_types.ObjectID                     `json:"objectId"`
	Version  SafeSuiBigInt[sui_types.SequenceNumber] `json:"version"`
}

type SuiNamePage = Page[string, sui_types.ObjectID]
