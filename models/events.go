package models

import "github.com/howjmay/sui-go/sui_types"

type EventId struct {
	TxDigest sui_types.TransactionDigest `json:"txDigest"`
	EventSeq SafeSuiBigInt[uint64]       `json:"eventSeq"`
}

// refer struct `SuiEvent` in `crates/sui-json-rpc-types/src/sui_event.rs`
type SuiEvent struct {
	Id EventId `json:"id"`
	// Move package where this event was emitted.
	PackageId *sui_types.ObjectID `json:"packageId"`
	// Move module where this event was emitted.
	TransactionModule sui_types.Identifier `json:"transactionModule"`
	// Sender's Sui sui_types.address.
	Sender *sui_types.SuiAddress `json:"sender"`
	// Move event type.
	Type *sui_types.StructTag `json:"type"`
	// Parsed json value of the event
	ParsedJson interface{} `json:"parsedJson,omitempty"`
	// Base 58 encoded bcs bytes of the move event
	Bcs         sui_types.Base58       `json:"bcs"`
	TimestampMs *SafeSuiBigInt[uint64] `json:"timestampMs,omitempty"`
}

type EventPage = Page[SuiEvent, EventId]

type EventFilter struct {
	/// Query by sender address
	Sender *sui_types.SuiAddress `json:"Sender,omitempty"`
	/// Return events emitted by the given transaction
	///digest of the transaction, as base-64 encoded string
	Transaction *sui_types.TransactionDigest `json:"Transaction,omitempty"`
	/// Return events emitted in a specified Package.
	Package *sui_types.ObjectID `json:"Package,omitempty"`
	/// Return events emitted in a specified Move module.
	/// If the event is defined in Module A but emitted in a tx with Module B,
	/// query `MoveModule` by module B returns the event.
	/// Query `MoveEventModule` by module A returns the event too.
	MoveModule *EventFilterMoveModule `json:"MoveModule,omitempty"`
	/// Return events with the given Move event struct name (struct tag).
	/// For example, if the event is defined in `0xabcd::MyModule`, and named
	/// `Foo`, then the struct tag is `0xabcd::MyModule::Foo`.
	MoveEventType  *sui_types.StructTag       `json:"MoveEventType,omitempty"`
	MoveEventField *EventFilterMoveEventField `json:"MoveEventField,omitempty"`
	// Return events emitted in [start_time, end_time] interval
	TimeRange *EventFilterTimeRange `json:"TimeRange,omitempty"`

	All *[]EventFilter    `json:"All,omitempty"`
	Any *[]EventFilter    `json:"Any,omitempty"`
	And *AndOrEventFilter `json:"And,omitempty"`
	Or  *AndOrEventFilter `json:"Or,omitempty"`
}

type EventFilterMoveModule struct {
	// the Move package ID
	Package *sui_types.ObjectID `json:"package"`
	// the module name
	Module sui_types.Identifier `json:"module"`
}

type EventFilterMoveEventField struct {
	Path string `json:"path"`
	// FIXME may need to be enum
	Value interface{} `json:"value"`
}

type EventFilterTimeRange struct {
	// left endpoint of time interval, milliseconds since epoch, inclusive
	StartTime SafeSuiBigInt[uint64] `json:"startTime"`
	// right endpoint of time interval, milliseconds since epoch, exclusive
	EndTime SafeSuiBigInt[uint64] `json:"endTime"`
}

type AndOrEventFilter struct {
	Filter1 *EventFilter
	Filter2 *EventFilter
}
