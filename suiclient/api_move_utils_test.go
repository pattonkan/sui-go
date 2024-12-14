package suiclient_test

import (
	"context"
	"testing"

	"github.com/howjmay/sui-go/sui"
	"github.com/howjmay/sui-go/suiclient"
	"github.com/howjmay/sui-go/suiclient/conn"
	"github.com/stretchr/testify/require"
)

func TestGetMoveFunctionArgTypes(t *testing.T) {
	client := suiclient.NewClient(conn.MainnetEndpointUrl)

	t.Run("address value", func(t *testing.T) {
		res, err := client.GetMoveFunctionArgTypes(
			context.Background(),
			sui.MustPackageIdFromHex("0x000000000000000000000000000000000000000000000000000000000000dee9"),
			"clob_v2",
			"usr_open_orders_for_address",
		)
		require.NoError(t, err)
		require.Equal(t, []sui.MoveFunctionArgType{sui.MoveFunctionArgTypeImmutableReference, sui.MoveFunctionArgTypePure}, res)
	})

	t.Run("objects", func(t *testing.T) {
		res, err := client.GetMoveFunctionArgTypes(
			context.Background(),
			sui.MustPackageIdFromHex("0x000000000000000000000000000000000000000000000000000000000000dee9"),
			"clob_v2",
			"withdraw_fees",
		)
		require.NoError(t, err)
		require.Equal(t, []sui.MoveFunctionArgType{sui.MoveFunctionArgTypeImmutableReference, sui.MoveFunctionArgTypeMutableReference, sui.MoveFunctionArgTypeMutableReference}, res)
	})
}

func TestGetNormalizedMoveFunction(t *testing.T) {
	client := suiclient.NewClient(conn.MainnetEndpointUrl)

	t.Run("address value", func(t *testing.T) {
		res, err := client.GetNormalizedMoveFunction(
			context.Background(),
			sui.MustPackageIdFromHex("0x000000000000000000000000000000000000000000000000000000000000dee9"),
			"clob_v2",
			"usr_open_orders_for_address",
		)
		require.NoError(t, err)
		expected := &sui.MoveNormalizedFunction{
			Visibility:     sui.MoveVisibilityPublic,
			IsEntry:        false,
			TypeParameters: []sui.MoveAbilitySet{{Abilities: []sui.MoveAbility{}}, {Abilities: []sui.MoveAbility{}}},
			Parameters: []sui.MoveNormalizedType{
				{Reference: &sui.MoveNormalizedType{
					Struct: &sui.MoveNormalizedTypeStructType{
						Address: sui.MustAddressFromHex("0xdee9"),
						Module:  sui.Identifier("clob_v2"),
						Name:    sui.Identifier("Pool"),
						TypeArguments: []sui.MoveNormalizedType{
							{TypeParameter: sui.NewMoveTypeParameterIndex(0)},
							{TypeParameter: sui.NewMoveTypeParameterIndex(1)},
						},
					},
				}},
				{Address: &sui.EmptyEnum{}},
			},
			Return: []sui.MoveNormalizedType{
				{Reference: &sui.MoveNormalizedType{
					Struct: &sui.MoveNormalizedTypeStructType{
						Address: sui.MustAddressFromHex("0x2"),
						Module:  sui.Identifier("linked_table"),
						Name:    sui.Identifier("LinkedTable"),
						TypeArguments: []sui.MoveNormalizedType{
							{U64: &sui.EmptyEnum{}},
							{U64: &sui.EmptyEnum{}},
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
			sui.MustPackageIdFromHex("0x000000000000000000000000000000000000000000000000000000000000dee9"),
			"clob_v2",
			"withdraw_fees",
		)
		require.NoError(t, err)
		expected := &sui.MoveNormalizedFunction{
			Visibility:     sui.MoveVisibilityPublic,
			IsEntry:        false,
			TypeParameters: []sui.MoveAbilitySet{{Abilities: []sui.MoveAbility{}}, {Abilities: []sui.MoveAbility{}}},
			Parameters: []sui.MoveNormalizedType{
				{Reference: &sui.MoveNormalizedType{
					Struct: &sui.MoveNormalizedTypeStructType{
						Address:       sui.MustAddressFromHex("0xdee9"),
						Module:        sui.Identifier("clob_v2"),
						Name:          sui.Identifier("PoolOwnerCap"),
						TypeArguments: []sui.MoveNormalizedType{},
					},
				}},
				{MutableReference: &sui.MoveNormalizedType{
					Struct: &sui.MoveNormalizedTypeStructType{
						Address: sui.MustAddressFromHex("0xdee9"),
						Module:  sui.Identifier("clob_v2"),
						Name:    sui.Identifier("Pool"),
						TypeArguments: []sui.MoveNormalizedType{
							{TypeParameter: sui.NewMoveTypeParameterIndex(0)},
							{TypeParameter: sui.NewMoveTypeParameterIndex(1)},
						},
					},
				}},
				{MutableReference: &sui.MoveNormalizedType{
					Struct: &sui.MoveNormalizedTypeStructType{
						Address:       sui.MustAddressFromHex("0x2"),
						Module:        sui.Identifier("tx_context"),
						Name:          sui.Identifier("TxContext"),
						TypeArguments: []sui.MoveNormalizedType{},
					},
				}},
			},
			Return: []sui.MoveNormalizedType{
				{Struct: &sui.MoveNormalizedTypeStructType{
					Address: sui.MustAddressFromHex("0x2"),
					Module:  sui.Identifier("coin"),
					Name:    sui.Identifier("Coin"),
					TypeArguments: []sui.MoveNormalizedType{
						{TypeParameter: sui.NewMoveTypeParameterIndex(1)},
					},
				}},
			},
		}
		require.Equal(t, expected, res)
	})
}

func TestGetNormalizedMoveModule(t *testing.T) {
	client := suiclient.NewClient(conn.MainnetEndpointUrl)

	res, err := client.GetNormalizedMoveModule(
		context.Background(),
		sui.MustPackageIdFromHex("0x06864a6f921804860930db6ddbe2e16acdf8504495ea7481637a1c8b9a8fe54b"),
		"cetus",
	)
	require.NoError(t, err)
	expected := &sui.MoveNormalizedModule{
		FileFormatVersion: 6,
		Address:           sui.MustAddressFromHex("0x6864a6f921804860930db6ddbe2e16acdf8504495ea7481637a1c8b9a8fe54b"),
		Name:              "cetus",
		Friends:           []*sui.MoveModuleId{},
		Structs: map[sui.Identifier]*sui.MoveNormalizedStruct{
			"CETUS": {
				Abilities:      sui.MoveAbilitySet{Abilities: []sui.MoveAbility{sui.MoveAbilityDrop}},
				TypeParameters: []*sui.MoveStructTypeParameter{},
				Fields:         []*sui.MoveNormalizedField{{Name: sui.Identifier("dummy_field"), Type: &sui.MoveNormalizedType{Bool: &sui.EmptyEnum{}}}},
			},
			"InitEvent": {
				Abilities:      sui.MoveAbilitySet{Abilities: []sui.MoveAbility{sui.MoveAbilityCopy, sui.MoveAbilityDrop}},
				TypeParameters: []*sui.MoveStructTypeParameter{},
				Fields: []*sui.MoveNormalizedField{
					{
						Name: sui.Identifier("cap_id"),
						Type: &sui.MoveNormalizedType{Struct: &sui.MoveNormalizedTypeStructType{
							Address:       sui.MustAddressFromHex("0x2"),
							Module:        sui.Identifier("object"),
							Name:          sui.Identifier("ID"),
							TypeArguments: []sui.MoveNormalizedType{},
						}},
					},
					{
						Name: sui.Identifier("metadata_id"),
						Type: &sui.MoveNormalizedType{Struct: &sui.MoveNormalizedTypeStructType{
							Address:       sui.MustAddressFromHex("0x2"),
							Module:        sui.Identifier("object"),
							Name:          sui.Identifier("ID"),
							TypeArguments: []sui.MoveNormalizedType{},
						}},
					},
				},
			},
		},
		ExposedFunctions: map[sui.Identifier]*sui.MoveNormalizedFunction{},
	}
	require.Equal(t, expected, res)
}

func TestGetNormalizedMoveModulesByPackage(t *testing.T) {
	client := suiclient.NewClient(conn.MainnetEndpointUrl)

	res, err := client.GetNormalizedMoveModulesByPackage(
		context.Background(),
		sui.MustPackageIdFromHex("0x06864a6f921804860930db6ddbe2e16acdf8504495ea7481637a1c8b9a8fe54b"),
	)
	require.NoError(t, err)
	expected := map[sui.Identifier]*sui.MoveNormalizedModule{
		"cetus": {
			FileFormatVersion: 6,
			Address:           sui.MustAddressFromHex("0x6864a6f921804860930db6ddbe2e16acdf8504495ea7481637a1c8b9a8fe54b"),
			Name:              "cetus",
			Friends:           []*sui.MoveModuleId{},
			Structs: map[sui.Identifier]*sui.MoveNormalizedStruct{
				"CETUS": {
					Abilities:      sui.MoveAbilitySet{Abilities: []sui.MoveAbility{sui.MoveAbilityDrop}},
					TypeParameters: []*sui.MoveStructTypeParameter{},
					Fields:         []*sui.MoveNormalizedField{{Name: sui.Identifier("dummy_field"), Type: &sui.MoveNormalizedType{Bool: &sui.EmptyEnum{}}}},
				},
				"InitEvent": {
					Abilities:      sui.MoveAbilitySet{Abilities: []sui.MoveAbility{sui.MoveAbilityCopy, sui.MoveAbilityDrop}},
					TypeParameters: []*sui.MoveStructTypeParameter{},
					Fields: []*sui.MoveNormalizedField{
						{
							Name: sui.Identifier("cap_id"),
							Type: &sui.MoveNormalizedType{Struct: &sui.MoveNormalizedTypeStructType{
								Address:       sui.MustAddressFromHex("0x2"),
								Module:        sui.Identifier("object"),
								Name:          sui.Identifier("ID"),
								TypeArguments: []sui.MoveNormalizedType{},
							}},
						},
						{
							Name: sui.Identifier("metadata_id"),
							Type: &sui.MoveNormalizedType{Struct: &sui.MoveNormalizedTypeStructType{
								Address:       sui.MustAddressFromHex("0x2"),
								Module:        sui.Identifier("object"),
								Name:          sui.Identifier("ID"),
								TypeArguments: []sui.MoveNormalizedType{},
							}},
						},
					},
				},
			},
			ExposedFunctions: map[sui.Identifier]*sui.MoveNormalizedFunction{},
		},
	}
	require.Equal(t, expected, res)
}

func TestGetNormalizedMoveStruct(t *testing.T) {
	client := suiclient.NewClient(conn.MainnetEndpointUrl)

	res, err := client.GetNormalizedMoveStruct(
		context.Background(),
		sui.MustPackageIdFromHex("0x06864a6f921804860930db6ddbe2e16acdf8504495ea7481637a1c8b9a8fe54b"),
		sui.Identifier("cetus"),
		sui.Identifier("CETUS"),
	)
	require.NoError(t, err)
	expected := &sui.MoveNormalizedStruct{
		Abilities:      sui.MoveAbilitySet{Abilities: []sui.MoveAbility{sui.MoveAbilityDrop}},
		TypeParameters: []*sui.MoveStructTypeParameter{},
		Fields:         []*sui.MoveNormalizedField{{Name: sui.Identifier("dummy_field"), Type: &sui.MoveNormalizedType{Bool: &sui.EmptyEnum{}}}},
	}
	require.Equal(t, expected, res)
}
