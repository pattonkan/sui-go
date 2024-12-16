package suiclient

import (
	"encoding/json"

	"github.com/pattonkan/sui-go/sui"
)

type EventId struct {
	TxDigest sui.TransactionDigest `json:"txDigest"`
	EventSeq *sui.BigInt           `json:"eventSeq"`
}

// refer struct `Event` in `crates/sui-json-rpc-types/src/sui_event.rs`
type Event struct {
	Id EventId `json:"id"`
	// Move package where this event was emitted.
	PackageId *sui.ObjectId `json:"packageId"`
	// Move module where this event was emitted.
	TransactionModule sui.Identifier `json:"transactionModule"`
	// Sender's Sui address.
	Sender *sui.Address `json:"sender"`
	// Move event type.
	Type *sui.StructTag `json:"type"`
	// Parsed json value of the event
	ParsedJson interface{} `json:"parsedJson,omitempty"`
	// Base 58 encoded bcs bytes of the move event
	Bcs         sui.Base58  `json:"bcs"`
	TimestampMs *sui.BigInt `json:"timestampMs,omitempty"`
}

type EventPage = Page[Event, EventId]

type EventFilter struct {
	// Query by sender address
	Sender *sui.Address `json:"Sender,omitempty"`
	// Return events emitted by the given transaction
	//digest of the transaction, as base-64 encoded string
	Transaction *sui.TransactionDigest `json:"Transaction,omitempty"`
	// Return events emitted in a specified Package.
	Package *sui.ObjectId `json:"Package,omitempty"`
	// Return events emitted in a specified Move module.
	// If the event is defined in Module A but emitted in a tx with Module B,
	// query `MoveModule` by module B returns the event.
	// Query `MoveEventModule` by module A returns the event too.
	MoveModule *EventFilterMoveModule `json:"MoveModule,omitempty"`
	// Return events with the given Move event struct name (struct tag).
	// For example, if the event is defined in `0xabcd::MyModule`, and named
	// `Foo`, then the struct tag is `0xabcd::MyModule::Foo`.
	MoveEventType *sui.StructTag `json:"MoveEventType,omitempty"`
	// Return the Event objects which contain the given field/value
	MoveEventField *EventFilterMoveEventField `json:"MoveEventField,omitempty"`
	// Return events emitted in [start_time, end_time] interval
	TimeRange *EventFilterTimeRange `json:"TimeRange,omitempty"`

	All *[]EventFilter    `json:"All,omitempty"`
	Any *[]EventFilter    `json:"Any,omitempty"`
	And *AndOrEventFilter `json:"And,omitempty"`
	Or  *AndOrEventFilter `json:"Or,omitempty"`
}

type EventFilterMoveModule struct {
	// the Move package Id
	Package *sui.ObjectId `json:"package"`
	// the module name
	Module sui.Identifier `json:"module"`
}

type EventFilterMoveEventField struct {
	Path string `json:"path"`
	// FIXME may need to be enum
	Value interface{} `json:"value"`
}

type EventFilterTimeRange struct {
	// left endpoint of time interval, milliseconds since epoch, inclusive
	StartTime *sui.BigInt `json:"startTime"`
	// right endpoint of time interval, milliseconds since epoch, exclusive
	EndTime *sui.BigInt `json:"endTime"`
}

type AndOrEventFilter struct {
	Filter1 *EventFilter
	Filter2 *EventFilter
}

func (f AndOrEventFilter) MarshalJSON() ([]byte, error) {
	return json.Marshal([2]interface{}{f.Filter1, f.Filter2})
}
