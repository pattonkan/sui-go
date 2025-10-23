package movebcs

import "github.com/pattonkan/sui-go/sui"

type MoveUpgradeCap struct {
	Id        *sui.ObjectId
	PackageID *sui.PackageId
	Version   uint64
	Policy    uint8
}
