package sui

type Data struct {
	Move    *MoveObject
	Package *MovePackage
}

func (d Data) IsBcsEnum() {}

type MoveObject struct {
	Type              MoveObjectType
	HasPublicTransfer bool
	Version           SequenceNumber
	Contents          []uint8
}

type Owner struct {
	AddressOwner *Address `json:"AddressOwner"`
	ObjectOwner  *Address `json:"ObjectOwner"`
	Shared       *struct {
		InitialSharedVersion SequenceNumber `json:"initial_shared_version"`
	} `json:"Shared,omitempty"`
	Immutable *EmptyEnum `json:"Immutable,omitempty"`
}

func (o Owner) IsBcsEnum() {}

func (o Owner) Tag() string {
	return ""
}

func (o Owner) Content() string {
	return ""
}
