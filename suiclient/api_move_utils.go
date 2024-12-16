package suiclient

import (
	"context"

	"github.com/pattonkan/sui-go/sui"
)

func (s *ClientImpl) GetMoveFunctionArgTypes(
	ctx context.Context,
	packageId *sui.PackageId,
	module sui.Identifier,
	function sui.Identifier,
) ([]sui.MoveFunctionArgType, error) {
	var resp []sui.MoveFunctionArgType
	return resp, s.http.CallContext(ctx, &resp, getMoveFunctionArgTypes, packageId, module, function)
}

func (s *ClientImpl) GetNormalizedMoveFunction(
	ctx context.Context,
	packageId *sui.PackageId,
	module sui.Identifier,
	function sui.Identifier,
) (*sui.MoveNormalizedFunction, error) {
	var resp sui.MoveNormalizedFunction
	return &resp, s.http.CallContext(ctx, &resp, getNormalizedMoveFunction, packageId, module, function)
}

func (s *ClientImpl) GetNormalizedMoveModule(
	ctx context.Context,
	packageId *sui.PackageId,
	module sui.Identifier,
) (*sui.MoveNormalizedModule, error) {
	var resp sui.MoveNormalizedModule
	return &resp, s.http.CallContext(ctx, &resp, getNormalizedMoveModule, packageId, module)
}

func (s *ClientImpl) GetNormalizedMoveModulesByPackage(
	ctx context.Context,
	packageId *sui.PackageId,
) (map[sui.Identifier]*sui.MoveNormalizedModule, error) {
	var resp map[sui.Identifier]*sui.MoveNormalizedModule
	return resp, s.http.CallContext(ctx, &resp, getNormalizedMoveModulesByPackage, packageId)
}

func (s *ClientImpl) GetNormalizedMoveStruct(
	ctx context.Context,
	packageId *sui.PackageId,
	module sui.Identifier,
	object sui.Identifier,
) (*sui.MoveNormalizedStruct, error) {
	var resp sui.MoveNormalizedStruct
	return &resp, s.http.CallContext(ctx, &resp, getNormalizedMoveStruct, packageId, module, object)
}
