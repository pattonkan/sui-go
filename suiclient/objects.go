package suiclient

import (
	"encoding/json"
	"strconv"
	"strings"

	"github.com/pattonkan/sui-go/sui"
)

type SuiObjectRef struct {
	/** Base64 string representing the object digest */
	Digest *sui.TransactionDigest `json:"digest"`
	/** Hex code as string representing the object id */
	ObjectId *sui.ObjectId `json:"objectId"`
	/** Object version */
	Version sui.SequenceNumber `json:"version"`
}

type SuiGasData struct {
	Payment []SuiObjectRef `json:"payment"`
	/** Gas Object's owner */
	Owner  *sui.ObjectId `json:"owner"`
	Price  *sui.BigInt   `json:"price"`
	Budget *sui.BigInt   `json:"budget"`
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
	Type              string          `json:"type"`
	HasPublicTransfer bool            `json:"hasPublicTransfer"`
	Fields            json.RawMessage `json:"fields"`
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

// FIXME Replace with sui.MoveObject
type SuiRawMoveObject struct {
	Type              sui.StructTag      `json:"type"`
	HasPublicTransfer bool               `json:"hasPublicTransfer"`
	Version           sui.SequenceNumber `json:"version"`
	BcsBytes          sui.Base64Data     `json:"bcsBytes"`
}

// FIXME Replace with sui.MovePackage
type SuiRawMovePackage struct {
	Id              *sui.ObjectId             `json:"id"`
	Version         sui.SequenceNumber        `json:"version"`
	ModuleMap       map[string]sui.Base64Data `json:"moduleMap"`
	TypeOriginTable []sui.TypeOrigin          `json:"typeOriginTable"`
	LinkageTable    map[sui.ObjectId]sui.UpgradeInfo
}

type SuiObjectData struct {
	ObjectId *sui.ObjectId     `json:"objectId"`
	Version  *sui.BigInt       `json:"version"`
	Digest   *sui.ObjectDigest `json:"digest"`
	/**
	 * Type of the object, default to be undefined unless SuiObjectDataOptions.showType is set to true
	 */
	Type *sui.ObjectType `json:"type,omitempty"`
	/**
	 * Move object content or package content, default to be undefined unless SuiObjectDataOptions.showContent is set to true
	 */
	Content *WrapperTaggedJson[SuiParsedData] `json:"content,omitempty"`
	/**
	 * Move object content or package content in BCS bytes, default to be undefined unless SuiObjectDataOptions.showBcs is set to true
	 */
	Bcs *WrapperTaggedJson[SuiRawData] `json:"bcs,omitempty"`
	/**
	 * The owner of this object. Default to be undefined unless SuiObjectDataOptions.showOwner is set to true
	 */
	Owner *ObjectOwner `json:"owner,omitempty"`
	/**
	 * The digest of the transaction that created or last mutated this object.
	 * Default to be undefined unless SuiObjectDataOptions.showPreviousTransaction is set to true
	 */
	PreviousTransaction *sui.TransactionDigest `json:"previousTransaction,omitempty"`
	/**
	 * The amount of SUI we would rebate if this object gets deleted.
	 * This number is re-calculated each time the object is mutated based on
	 * the present storage gas price.
	 * Default to be undefined unless SuiObjectDataOptions.showStorageRebate is set to true
	 */
	StorageRebate *sui.BigInt `json:"storageRebate,omitempty"`
	/**
	 * Display metadata for this object, default to be undefined unless SuiObjectDataOptions.showDisplay is set to true
	 * This can also be None if the struct type does not have Display defined
	 * See more details in https://forums.sui.io/t/nft-object-display-proposal/4872
	 */
	Display interface{} `json:"display,omitempty"`
}

func (data *SuiObjectData) Ref() sui.ObjectRef {
	return sui.ObjectRef{
		ObjectId: data.ObjectId,
		Version:  data.Version.Uint64(),
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
		ObjectId sui.ObjectId `json:"object_id"`
	} `json:"notExists,omitempty"`
	Deleted *struct {
		ObjectId sui.ObjectId       `json:"object_id"`
		Version  sui.SequenceNumber `json:"version"`
		Digest   sui.ObjectDigest   `json:"digest"`
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
	Data  *SuiObjectData                             `json:"data,omitempty"`
	Error *WrapperTaggedJson[SuiObjectResponseError] `json:"error,omitempty"`
}

func (r *SuiObjectResponse) GetMoveObjectInBcs() []byte {
	return r.Data.Bcs.Data.MoveObject.BcsBytes
}

type CheckpointedObjectId struct {
	ObjectId     sui.ObjectId `json:"objectId"`
	AtCheckpoint *sui.BigInt  `json:"atCheckpoint"`
}

type ObjectsPage = Page[SuiObjectResponse, sui.ObjectId]

type SuiObjectDataFilter struct {
	MatchAll  []*SuiObjectDataFilter `json:"MatchAll,omitempty"`
	MatchAny  []*SuiObjectDataFilter `json:"MatchAny,omitempty"`
	MatchNone []*SuiObjectDataFilter `json:"MatchNone,omitempty"`
	// Query by type a specified Package.
	Package *sui.ObjectId `json:"Package,omitempty"`
	// Query by type a specified Move module.
	MoveModule *sui.MoveModule `json:"MoveModule,omitempty"`
	// Query by type
	StructType   *sui.StructTag `json:"StructType,omitempty"`
	AddressOwner *sui.Address   `json:"AddressOwner,omitempty"`
	ObjectOwner  *sui.ObjectId  `json:"ObjectOwner,omitempty"`
	ObjectId     *sui.ObjectId  `json:"ObjectId,omitempty"`
	// allow querying for multiple object ids
	ObjectIds []*sui.ObjectId `json:"ObjectIds,omitempty"`
	Version   *sui.BigInt     `json:"Version,omitempty"`
}

type SuiObjectResponseQuery struct {
	Filter  *SuiObjectDataFilter  `json:"filter,omitempty"`
	Options *SuiObjectDataOptions `json:"options,omitempty"`
}

type SuiPastObjectResponse = WrapperTaggedJson[SuiPastObject]

type SuiPastObject struct {
	// The object exists and is found with this version
	VersionFound *SuiObjectData `json:"VersionFound,omitempty"`
	// The object does not exist
	ObjectNotExists *sui.ObjectId `json:"ObjectNotExists,omitempty"`
	// The object is found to be deleted with this version
	ObjectDeleted *SuiObjectRef `json:"ObjectDeleted,omitempty"`
	// The object exists but not found with this version
	VersionNotFound *VersionNotFoundData `json:"VersionNotFound,omitempty"`
	// The asked object version is higher than the latest
	VersionTooHigh *struct {
		ObjectId      sui.ObjectId       `json:"object_id"`
		AskedVersion  sui.SequenceNumber `json:"asked_version"`
		LatestVersion sui.SequenceNumber `json:"latest_version"`
	} `json:"VersionTooHigh,omitempty"`
}

type VersionNotFoundData struct {
	ObjectId       *sui.ObjectId
	SequenceNumber sui.SequenceNumber
}

func (c *VersionNotFoundData) UnmarshalJSON(data []byte) error {
	var err error
	input := data[1 : len(data)-2]
	elts := strings.Split(string(input), ",")
	c.ObjectId, err = sui.ObjectIdFromHex(elts[0][1 : len(elts[0])-2])
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
	ObjectId *sui.ObjectId `json:"objectId"`
	Version  *sui.BigInt   `json:"version"`
}

type SuiNamePage = Page[string, sui.ObjectId]

type SuiTypeTag string
type SuiJsonValue string
