package sui

import (
	"context"

	"github.com/howjmay/sui-go/models"
	"github.com/howjmay/sui-go/sui_types"
)

func (s *ImplSuiAPI) GetMoveFunctionArgTypes(
	ctx context.Context,
	packageID *sui_types.PackageID,
	module sui_types.Identifier,
	function sui_types.Identifier,
) ([]models.MoveFunctionArgType, error) {
	var resp []models.MoveFunctionArgType
	return resp, s.http.CallContext(ctx, &resp, getMoveFunctionArgTypes, packageID, module, function)
}

func (s *ImplSuiAPI) GetNormalizedMoveFunction(
	ctx context.Context,
	packageID *sui_types.PackageID,
	module sui_types.Identifier,
	function sui_types.Identifier,
) (*models.MoveNormalizedFunction, error) {
	var resp models.MoveNormalizedFunction
	return &resp, s.http.CallContext(ctx, &resp, getNormalizedMoveFunction, packageID, module, function)
}

func (s *ImplSuiAPI) GetNormalizedMoveModule(
	ctx context.Context,
	packageID *sui_types.PackageID,
	module sui_types.Identifier,
) (*models.MoveNormalizedModule, error) {
	var resp models.MoveNormalizedModule
	return &resp, s.http.CallContext(ctx, &resp, getNormalizedMoveModule, packageID, module)
}

func (s *ImplSuiAPI) GetNormalizedMoveModulesByPackage(
	ctx context.Context,
	packageID *sui_types.PackageID,
) (map[sui_types.Identifier]*models.MoveNormalizedModule, error) {
	var resp map[sui_types.Identifier]*models.MoveNormalizedModule
	return resp, s.http.CallContext(ctx, &resp, getNormalizedMoveModulesByPackage, packageID)
}

func (s *ImplSuiAPI) GetNormalizedMoveStruct(
	ctx context.Context,
	packageID *sui_types.PackageID,
	module sui_types.Identifier,
	object sui_types.Identifier,
) (*models.MoveNormalizedStruct, error) {
	var resp models.MoveNormalizedStruct
	return &resp, s.http.CallContext(ctx, &resp, getNormalizedMoveStruct, packageID, module, object)
}
