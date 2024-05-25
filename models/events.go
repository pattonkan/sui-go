package models

import "github.com/howjmay/sui-go/sui_types"

type EventId struct {
	TxDigest sui_types.TransactionDigest `json:"txDigest"`
	EventSeq *BigInt                     `json:"eventSeq"`
}

type SuiEvent struct {
	Id EventId `json:"id"`
	// Move package where this event was emitted.
	PackageId sui_types.ObjectID `json:"packageId"`
	// Move module where this event was emitted.
	TransactionModule string `json:"transactionModule"`
	// Sender's Sui sui_types.address.
	Sender sui_types.SuiAddress `json:"sender"`
	// Move event type.
	Type string `json:"type"`
	// Parsed json value of the event
	ParsedJson interface{} `json:"parsedJson,omitempty"`
	// Base 58 encoded bcs bytes of the move event
	Bcs         sui_types.Base58 `json:"bcs"`
	TimestampMs *BigInt          `json:"timestampMs,omitempty"`
}

type EventPage = Page[SuiEvent, EventId]

type EventFilter struct {
	// Query by sender sui_types.address.
	Sender *sui_types.SuiAddress `json:"Sender,omitempty"`
	// Return events emitted by the given transaction.
	Transaction *sui_types.TransactionDigest `json:"Transaction,omitempty"`
	//digest of the transaction, as base-64 encoded string

	// Return events emitted in a specified Package.
	Package *sui_types.ObjectID `json:"Package,omitempty"`
	// Return events emitted in a specified Move module.
	MoveModule *EventFilterMoveModule `json:"MoveModule,omitempty"`
	// Return events with the given move event struct name
	MoveEventType  *string                    `json:"MoveEventType,omitempty"`
	MoveEventField *EventFilterMoveEventField `json:"MoveEventField,omitempty"`
	// Return events emitted in [start_time, end_time] interval
	TimeRange *EventFilterTimeRange `json:"TimeRange,omitempty"`

	All *[]EventFilter `json:"All,omitempty"`
	Any *[]EventFilter `json:"Any,omitempty"`
	// FIXME enable 'And' and 'Or' filter
	//And *struct {
	//	*EventFilter
	//	*EventFilter
	//} `json:"And,omitempty"`
	//Or *struct {
	//	EventFilter
	//	EventFilter
	//} `json:"Or,omitempty"`
}

type EventFilterMoveModule struct {
	// the Move package ID
	Package sui_types.ObjectID `json:"package"`
	// the module name
	Module string `json:"module"`
}

type EventFilterMoveEventField struct {
	Path  string      `json:"path"`
	Value interface{} `json:"value"`
}

type EventFilterTimeRange struct {
	// left endpoint of time interval, milliseconds since epoch, inclusive
	StartTime *BigInt `json:"startTime"`
	// right endpoint of time interval, milliseconds since epoch, exclusive
	EndTime *BigInt `json:"endTime"`
}
