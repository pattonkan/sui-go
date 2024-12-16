package suiptb

import "github.com/pattonkan/sui-go/sui"

type CallArg struct {
	Pure   *[]byte
	Object *ObjectArg
}

func (c CallArg) IsBcsEnum() {}

type ObjectArg struct {
	ImmOrOwnedObject *sui.ObjectRef
	SharedObject     *SharedObjectArg
	Receiving        *sui.ObjectRef
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
