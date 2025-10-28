package suiptb

import "github.com/pattonkan/sui-go/sui"

type CallArg struct {
	Pure   *[]byte    `json:"pure,omitempty"`
	Object *ObjectArg `json:"object,omitempty"`
}

func (c CallArg) IsBcsEnum() {}

type ObjectArg struct {
	ImmOrOwnedObject *sui.ObjectRef   `json:"imm_or_owned_object,omitempty"`
	SharedObject     *SharedObjectArg `json:"shared_object,omitempty"`
	Receiving        *sui.ObjectRef   `json:"receiving,omitempty"`
}

type SharedObjectArg struct {
	Id                   *sui.ObjectId
	InitialSharedVersion sui.SequenceNumber
	Mutable              bool
}

func (o ObjectArg) IsBcsEnum() {}

func (o ObjectArg) id() *sui.ObjectId {
	switch {
	case o.ImmOrOwnedObject != nil:
		return o.ImmOrOwnedObject.ObjectId
	case o.SharedObject != nil:
		return o.SharedObject.Id
	case o.Receiving != nil:
		return o.Receiving.ObjectId
	default:
		return &sui.ObjectId{}
	}
}

func ObjectArgFromSharedObjectRef(ref *sui.SharedObjectRef, mutable bool) ObjectArg {
	return ObjectArg{
		SharedObject: &SharedObjectArg{
			Id:                   ref.Id,
			InitialSharedVersion: ref.InitialSharedVersion,
			Mutable:              mutable,
		},
	}
}
