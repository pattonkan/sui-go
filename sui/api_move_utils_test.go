package sui_test

import (
	"context"
	"testing"

	"github.com/howjmay/sui-go/models"
	"github.com/howjmay/sui-go/sui"
	"github.com/howjmay/sui-go/sui/conn"
	"github.com/howjmay/sui-go/sui_types"
	"github.com/howjmay/sui-go/sui_types/serialization"
	"github.com/stretchr/testify/require"
)

func TestGetMoveFunctionArgTypes(t *testing.T) {
	client := sui.NewSuiClient(conn.MainnetEndpointUrl)

	t.Run("address value", func(t *testing.T) {
		res, err := client.GetMoveFunctionArgTypes(
			context.Background(),
			sui_types.MustPackageIDFromHex("0x000000000000000000000000000000000000000000000000000000000000dee9"),
			"clob_v2",
			"usr_open_orders_for_address",
		)
		require.NoError(t, err)
		require.Equal(t, []models.MoveFunctionArgType{models.MoveFunctionArgTypeImmutableReference, models.MoveFunctionArgTypePure}, res)
	})

	t.Run("objects", func(t *testing.T) {
		res, err := client.GetMoveFunctionArgTypes(
			context.Background(),
			sui_types.MustPackageIDFromHex("0x000000000000000000000000000000000000000000000000000000000000dee9"),
			"clob_v2",
			"withdraw_fees",
		)
		require.NoError(t, err)
		require.Equal(t, []models.MoveFunctionArgType{models.MoveFunctionArgTypeImmutableReference, models.MoveFunctionArgTypeMutableReference, models.MoveFunctionArgTypeMutableReference}, res)
	})
}

func TestGetNormalizedMoveFunction(t *testing.T) {
	client := sui.NewSuiClient(conn.MainnetEndpointUrl)

	t.Run("address value", func(t *testing.T) {
		res, err := client.GetNormalizedMoveFunction(
			context.Background(),
			sui_types.MustPackageIDFromHex("0x000000000000000000000000000000000000000000000000000000000000dee9"),
			"clob_v2",
			"usr_open_orders_for_address",
		)
		require.NoError(t, err)
		expected := &models.MoveNormalizedFunction{
			Visibility:     models.MoveVisibilityPublic,
			IsEntry:        false,
			TypeParameters: []models.MoveAbilitySet{{Abilities: []models.MoveAbility{}}, {Abilities: []models.MoveAbility{}}},
			Parameters: []models.MoveNormalizedType{
				{Reference: &models.MoveNormalizedType{
					Struct: &models.MoveNormalizedTypeStructType{
						Address: sui_types.MustSuiAddressFromHex("0xdee9"),
						Module:  sui_types.Identifier("clob_v2"),
						Name:    sui_types.Identifier("Pool"),
						TypeArguments: []models.MoveNormalizedType{
							{TypeParameter: models.NewMoveTypeParameterIndex(0)},
							{TypeParameter: models.NewMoveTypeParameterIndex(1)},
						},
					},
				}},
				{Address: &serialization.EmptyEnum{}},
			},
			Return: []models.MoveNormalizedType{
				{Reference: &models.MoveNormalizedType{
					Struct: &models.MoveNormalizedTypeStructType{
						Address: sui_types.MustSuiAddressFromHex("0x2"),
						Module:  sui_types.Identifier("linked_table"),
						Name:    sui_types.Identifier("LinkedTable"),
						TypeArguments: []models.MoveNormalizedType{
							{U64: &serialization.EmptyEnum{}},
							{U64: &serialization.EmptyEnum{}},
						},
					},
				}},
			},
		}
		require.Equal(t, expected, res)
	})

	t.Run("objects", func(t *testing.T) {
		res, err := client.GetNormalizedMoveFunction(
			context.Background(),
			sui_types.MustPackageIDFromHex("0x000000000000000000000000000000000000000000000000000000000000dee9"),
			"clob_v2",
			"withdraw_fees",
		)
		require.NoError(t, err)
		expected := &models.MoveNormalizedFunction{
			Visibility:     models.MoveVisibilityPublic,
			IsEntry:        false,
			TypeParameters: []models.MoveAbilitySet{{Abilities: []models.MoveAbility{}}, {Abilities: []models.MoveAbility{}}},
			Parameters: []models.MoveNormalizedType{
				{Reference: &models.MoveNormalizedType{
					Struct: &models.MoveNormalizedTypeStructType{
						Address:       sui_types.MustSuiAddressFromHex("0xdee9"),
						Module:        sui_types.Identifier("clob_v2"),
						Name:          sui_types.Identifier("PoolOwnerCap"),
						TypeArguments: []models.MoveNormalizedType{},
					},
				}},
				{MutableReference: &models.MoveNormalizedType{
					Struct: &models.MoveNormalizedTypeStructType{
						Address: sui_types.MustSuiAddressFromHex("0xdee9"),
						Module:  sui_types.Identifier("clob_v2"),
						Name:    sui_types.Identifier("Pool"),
						TypeArguments: []models.MoveNormalizedType{
							{TypeParameter: models.NewMoveTypeParameterIndex(0)},
							{TypeParameter: models.NewMoveTypeParameterIndex(1)},
						},
					},
				}},
				{MutableReference: &models.MoveNormalizedType{
					Struct: &models.MoveNormalizedTypeStructType{
						Address:       sui_types.MustSuiAddressFromHex("0x2"),
						Module:        sui_types.Identifier("tx_context"),
						Name:          sui_types.Identifier("TxContext"),
						TypeArguments: []models.MoveNormalizedType{},
					},
				}},
			},
			Return: []models.MoveNormalizedType{
				{Struct: &models.MoveNormalizedTypeStructType{
					Address: sui_types.MustSuiAddressFromHex("0x2"),
					Module:  sui_types.Identifier("coin"),
					Name:    sui_types.Identifier("Coin"),
					TypeArguments: []models.MoveNormalizedType{
						{TypeParameter: models.NewMoveTypeParameterIndex(1)},
					},
				}},
			},
		}
		require.Equal(t, expected, res)
	})
}

func TestGetNormalizedMoveModule(t *testing.T) {
	client := sui.NewSuiClient(conn.MainnetEndpointUrl)

	res, err := client.GetNormalizedMoveModule(
		context.Background(),
		sui_types.MustPackageIDFromHex("0x06864a6f921804860930db6ddbe2e16acdf8504495ea7481637a1c8b9a8fe54b"),
		"cetus",
	)
	require.NoError(t, err)
	expected := &models.MoveNormalizedModule{
		FileFormatVersion: 6,
		Address:           sui_types.MustSuiAddressFromHex("0x6864a6f921804860930db6ddbe2e16acdf8504495ea7481637a1c8b9a8fe54b"),
		Name:              "cetus",
		Friends:           []*models.MoveModuleID{},
		Structs: map[sui_types.Identifier]*models.MoveNormalizedStruct{
			"CETUS": {
				Abilities:      models.MoveAbilitySet{Abilities: []models.MoveAbility{models.MoveAbilityDrop}},
				TypeParameters: []*models.MoveStructTypeParameter{},
				Fields:         []*models.MoveNormalizedField{{Name: sui_types.Identifier("dummy_field"), Type: &models.MoveNormalizedType{Bool: &serialization.EmptyEnum{}}}},
			},
			"InitEvent": {
				Abilities:      models.MoveAbilitySet{Abilities: []models.MoveAbility{models.MoveAbilityCopy, models.MoveAbilityDrop}},
				TypeParameters: []*models.MoveStructTypeParameter{},
				Fields: []*models.MoveNormalizedField{
					{
						Name: sui_types.Identifier("cap_id"),
						Type: &models.MoveNormalizedType{Struct: &models.MoveNormalizedTypeStructType{
							Address:       sui_types.MustSuiAddressFromHex("0x2"),
							Module:        sui_types.Identifier("object"),
							Name:          sui_types.Identifier("ID"),
							TypeArguments: []models.MoveNormalizedType{},
						}},
					},
					{
						Name: sui_types.Identifier("metadata_id"),
						Type: &models.MoveNormalizedType{Struct: &models.MoveNormalizedTypeStructType{
							Address:       sui_types.MustSuiAddressFromHex("0x2"),
							Module:        sui_types.Identifier("object"),
							Name:          sui_types.Identifier("ID"),
							TypeArguments: []models.MoveNormalizedType{},
						}},
					},
				},
			},
		},
		ExposedFunctions: map[sui_types.Identifier]*models.MoveNormalizedFunction{},
	}
	require.Equal(t, expected, res)
}

func TestGetNormalizedMoveModulesByPackage(t *testing.T) {
	client := sui.NewSuiClient(conn.MainnetEndpointUrl)

	res, err := client.GetNormalizedMoveModulesByPackage(
		context.Background(),
		sui_types.MustPackageIDFromHex("0x06864a6f921804860930db6ddbe2e16acdf8504495ea7481637a1c8b9a8fe54b"),
	)
	require.NoError(t, err)
	expected := map[sui_types.Identifier]*models.MoveNormalizedModule{
		"cetus": {
			FileFormatVersion: 6,
			Address:           sui_types.MustSuiAddressFromHex("0x6864a6f921804860930db6ddbe2e16acdf8504495ea7481637a1c8b9a8fe54b"),
			Name:              "cetus",
			Friends:           []*models.MoveModuleID{},
			Structs: map[sui_types.Identifier]*models.MoveNormalizedStruct{
				"CETUS": {
					Abilities:      models.MoveAbilitySet{Abilities: []models.MoveAbility{models.MoveAbilityDrop}},
					TypeParameters: []*models.MoveStructTypeParameter{},
					Fields:         []*models.MoveNormalizedField{{Name: sui_types.Identifier("dummy_field"), Type: &models.MoveNormalizedType{Bool: &serialization.EmptyEnum{}}}},
				},
				"InitEvent": {
					Abilities:      models.MoveAbilitySet{Abilities: []models.MoveAbility{models.MoveAbilityCopy, models.MoveAbilityDrop}},
					TypeParameters: []*models.MoveStructTypeParameter{},
					Fields: []*models.MoveNormalizedField{
						{
							Name: sui_types.Identifier("cap_id"),
							Type: &models.MoveNormalizedType{Struct: &models.MoveNormalizedTypeStructType{
								Address:       sui_types.MustSuiAddressFromHex("0x2"),
								Module:        sui_types.Identifier("object"),
								Name:          sui_types.Identifier("ID"),
								TypeArguments: []models.MoveNormalizedType{},
							}},
						},
						{
							Name: sui_types.Identifier("metadata_id"),
							Type: &models.MoveNormalizedType{Struct: &models.MoveNormalizedTypeStructType{
								Address:       sui_types.MustSuiAddressFromHex("0x2"),
								Module:        sui_types.Identifier("object"),
								Name:          sui_types.Identifier("ID"),
								TypeArguments: []models.MoveNormalizedType{},
							}},
						},
					},
				},
			},
			ExposedFunctions: map[sui_types.Identifier]*models.MoveNormalizedFunction{},
		},
	}
	require.Equal(t, expected, res)
}

func TestGetNormalizedMoveStruct(t *testing.T) {
	client := sui.NewSuiClient(conn.MainnetEndpointUrl)

	res, err := client.GetNormalizedMoveStruct(
		context.Background(),
		sui_types.MustPackageIDFromHex("0x06864a6f921804860930db6ddbe2e16acdf8504495ea7481637a1c8b9a8fe54b"),
		sui_types.Identifier("cetus"),
		sui_types.Identifier("CETUS"),
	)
	require.NoError(t, err)
	expected := &models.MoveNormalizedStruct{
		Abilities:      models.MoveAbilitySet{Abilities: []models.MoveAbility{models.MoveAbilityDrop}},
		TypeParameters: []*models.MoveStructTypeParameter{},
		Fields:         []*models.MoveNormalizedField{{Name: sui_types.Identifier("dummy_field"), Type: &models.MoveNormalizedType{Bool: &serialization.EmptyEnum{}}}},
	}
	require.Equal(t, expected, res)
}
