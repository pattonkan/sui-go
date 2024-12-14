package sui

type MovePackage struct {
	Id              ObjectId
	Version         SequenceNumber
	ModuleMap       map[string][]uint8
	TypeOriginTable []TypeOrigin
	LinkageTable    map[ObjectId]UpgradeInfo
}

type TypeOrigin struct {
	Module     string   `json:"moduleName"`
	StructName string   `json:"structName"`
	Package    ObjectId `json:"package"`
}

type UpgradeInfo struct {
	UpgradedId      ObjectId
	UpgradedVersion SequenceNumber
}
