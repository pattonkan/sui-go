package conn

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/Khan/genqlient/graphql"

	"github.com/pattonkan/sui-go/sui"
	"github.com/pattonkan/sui-go/suiclient/suigraphql"
)

type GraphQLClient struct {
	url    string
	client graphql.Client
}

func NewGraphQLClient(url string) *GraphQLClient {
	return &GraphQLClient{
		url: strings.TrimRight(url, "/"),
		client: graphql.NewClient(
			url,
			&http.Client{
				Timeout: 30 * time.Second,
			},
		),
	}
}

// GetGraphQLClient returns the underlying GraphQL client for custom GraphQL queries.
func (c *GraphQLClient) GetGraphQLClient() graphql.Client {
	return c.client
}

// Build and send a custom GraphQL query, and returned the response from the server as raw bytes directly.
func (c *GraphQLClient) Query(query string, variables map[string]interface{}) ([]byte, error) {
	reqBody := map[string]interface{}{
		"query":     query,
		"variables": variables,
	}
	reqBytes, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request body: %w", err)
	}

	req, err := http.NewRequest("POST", c.url, bytes.NewBuffer(reqBytes))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	rawBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	return rawBytes, nil
}

func (c *GraphQLClient) DevInspectTransactionBlock(
	ctx context.Context,
	txBytes string,
	txMeta suigraphql.TransactionMetadata,
	showBalanceChanges *bool,
	showEffects *bool,
	showRawEffects *bool,
	showEvents *bool,
	showInput *bool,
	showObjectChanges *bool,
	showRawInput *bool,
) (*suigraphql.DevInspectTransactionBlockResponse, error) {
	return suigraphql.DevInspectTransactionBlock(
		ctx,
		c.client,
		txBytes,
		txMeta,
		showBalanceChanges,
		showEffects,
		showRawEffects,
		showEvents,
		showInput,
		showObjectChanges,
		showRawInput,
	)
}

func (c *GraphQLClient) DryRunTransactionBlock(
	ctx context.Context,
	txBytes string,
	showBalanceChanges *bool,
	showEffects *bool,
	showRawEffects *bool,
	showEvents *bool,
	showInput *bool,
	showObjectChanges *bool,
	showRawInput *bool,
) (*suigraphql.DryRunTransactionBlockResponse, error) {
	return suigraphql.DryRunTransactionBlock(
		ctx,
		c.client,
		txBytes,
		showBalanceChanges,
		showEffects,
		showRawEffects,
		showEvents,
		showInput,
		showObjectChanges,
		showRawInput,
	)
}

func (c *GraphQLClient) ExecuteTransactionBlock(
	ctx context.Context,
	txBytes string,
	signatures []string,
	showBalanceChanges *bool,
	showEffects *bool,
	showRawEffects *bool,
	showEvents *bool,
	showInput *bool,
	showObjectChanges *bool,
	showRawInput *bool,
) (*suigraphql.ExecuteTransactionBlockResponse, error) {
	return suigraphql.ExecuteTransactionBlock(
		ctx,
		c.client,
		txBytes,
		signatures,
		showBalanceChanges,
		showEffects,
		showRawEffects,
		showEvents,
		showInput,
		showObjectChanges,
		showRawInput,
	)
}

func (c *GraphQLClient) GetAllBalances(
	ctx context.Context,
	owner sui.Address,
	limit *int,
	cursor *string,
) (*suigraphql.GetAllBalancesResponse, error) {
	return suigraphql.GetAllBalances(
		ctx,
		c.client,
		owner,
		limit,
		cursor,
	)
}

func (c *GraphQLClient) GetBalance(
	ctx context.Context,
	owner sui.Address,
	fetchCoinType *string,
) (*suigraphql.GetBalanceResponse, error) {
	return suigraphql.GetBalance(
		ctx,
		c.client,
		owner,
		fetchCoinType,
	)
}

func (c *GraphQLClient) GetChainIdentifier(
	ctx context.Context,
) (*suigraphql.GetChainIdentifierResponse, error) {
	return suigraphql.GetChainIdentifier(
		ctx,
		c.client,
	)
}

func (c *GraphQLClient) GetCheckpoint(
	ctx context.Context,
	id *suigraphql.CheckpointId,
) (*suigraphql.GetCheckpointResponse, error) {
	return suigraphql.GetCheckpoint(
		ctx,
		c.client,
		id,
	)
}

func (c *GraphQLClient) GetCheckpoints(
	ctx context.Context,
	first *int,
	before *string,
	last *int,
	after *string,
) (*suigraphql.GetCheckpointsResponse, error) {
	return suigraphql.GetCheckpoints(
		ctx,
		c.client,
		first,
		before,
		last,
		after,
	)
}

func (c *GraphQLClient) GetCoinMetadata(
	ctx context.Context,
	coinType string,
) (*suigraphql.GetCoinMetadataResponse, error) {
	return suigraphql.GetCoinMetadata(
		ctx,
		c.client,
		coinType,
	)
}

func (c *GraphQLClient) GetCoins(
	ctx context.Context,
	owner sui.Address,
	first *int,
	cursor *string,
	fetchCoinType *string,
) (*suigraphql.GetCoinsResponse, error) {
	return suigraphql.GetCoins(
		ctx,
		c.client,
		owner,
		first,
		cursor,
		fetchCoinType,
	)
}

func (c *GraphQLClient) GetCommitteeInfo(
	ctx context.Context,
	epochId *uint64,
	after *string,
) (*suigraphql.GetCommitteeInfoResponse, error) {
	return suigraphql.GetCommitteeInfo(
		ctx,
		c.client,
		epochId,
		after,
	)
}

func (c *GraphQLClient) GetCurrentEpoch(
	ctx context.Context,
) (*suigraphql.GetCurrentEpochResponse, error) {
	return suigraphql.GetCurrentEpoch(
		ctx,
		c.client,
	)
}

func (c *GraphQLClient) GetDynamicFieldObject(
	ctx context.Context,
	parentId sui.Address,
	name suigraphql.DynamicFieldName,
) (*suigraphql.GetDynamicFieldObjectResponse, error) {
	return suigraphql.GetDynamicFieldObject(
		ctx,
		c.client,
		parentId,
		name,
	)
}

func (c *GraphQLClient) GetDynamicFields(
	ctx context.Context,
	parentId sui.Address,
	first *int,
	cursor *string,
) (*suigraphql.GetDynamicFieldsResponse, error) {
	return suigraphql.GetDynamicFields(
		ctx,
		c.client,
		parentId,
		first,
		cursor,
	)
}

func (c *GraphQLClient) GetLatestCheckpointSequenceNumber(
	ctx context.Context,
) (*suigraphql.GetLatestCheckpointSequenceNumberResponse, error) {
	return suigraphql.GetLatestCheckpointSequenceNumber(
		ctx,
		c.client,
	)
}

func (c *GraphQLClient) GetLatestSuiSystemState(
	ctx context.Context,
) (*suigraphql.GetLatestSuiSystemStateResponse, error) {
	return suigraphql.GetLatestSuiSystemState(
		ctx,
		c.client,
	)
}

func (c *GraphQLClient) GetMoveFunctionArgTypes(
	ctx context.Context,
	packageId sui.Address,
	module string,
	function string,
) (*suigraphql.GetMoveFunctionArgTypesResponse, error) {
	return suigraphql.GetMoveFunctionArgTypes(
		ctx,
		c.client,
		packageId,
		module,
		function,
	)
}

func (c *GraphQLClient) GetNormalizedMoveFunction(
	ctx context.Context,
	packageId sui.Address,
	module string,
	function string,
) (*suigraphql.GetNormalizedMoveFunctionResponse, error) {
	return suigraphql.GetNormalizedMoveFunction(
		ctx,
		c.client,
		packageId,
		module,
		function,
	)
}

func (c *GraphQLClient) GetNormalizedMoveModule(
	ctx context.Context,
	packageId sui.Address,
	module string,
) (*suigraphql.GetNormalizedMoveModuleResponse, error) {
	return suigraphql.GetNormalizedMoveModule(
		ctx,
		c.client,
		packageId,
		module,
	)
}

func (c *GraphQLClient) GetNormalizedMoveModulesByPackage(
	ctx context.Context,
	packageId sui.Address,
	cursor *string,
) (*suigraphql.GetNormalizedMoveModulesByPackageResponse, error) {
	return suigraphql.GetNormalizedMoveModulesByPackage(
		ctx,
		c.client,
		packageId,
		cursor,
	)
}

func (c *GraphQLClient) GetNormalizedMoveStruct(
	ctx context.Context,
	packageId sui.Address,
	module string,
	moveStruct string,
) (*suigraphql.GetNormalizedMoveStructResponse, error) {
	return suigraphql.GetNormalizedMoveStruct(
		ctx,
		c.client,
		packageId,
		module,
		moveStruct,
	)
}

func (c *GraphQLClient) GetObject(
	ctx context.Context,
	id sui.Address,
	showBcs *bool,
	showOwner *bool,
	showPreviousTransaction *bool,
	showContent *bool,
	showDisplay *bool,
	showType *bool,
	showStorageRebate *bool,
) (*suigraphql.GetObjectResponse, error) {
	return suigraphql.GetObject(
		ctx,
		c.client,
		id,
		showBcs,
		showOwner,
		showPreviousTransaction,
		showContent,
		showDisplay,
		showType,
		showStorageRebate,
	)
}

func (c *GraphQLClient) GetOwnedObjects(
	ctx context.Context,
	owner sui.Address,
	limit *int,
	cursor *string,
	showBcs *bool,
	showContent *bool,
	showDisplay *bool,
	showType *bool,
	showOwner *bool,
	showPreviousTransaction *bool,
	showStorageRebate *bool,
	filter *suigraphql.ObjectFilter,
) (*suigraphql.GetOwnedObjectsResponse, error) {
	return suigraphql.GetOwnedObjects(
		ctx,
		c.client,
		owner,
		limit,
		cursor,
		showBcs,
		showContent,
		showDisplay,
		showType,
		showOwner,
		showPreviousTransaction,
		showStorageRebate,
		filter,
	)
}

func (c *GraphQLClient) GetProtocolConfig(
	ctx context.Context,
	protocolVersion *uint64,
) (*suigraphql.GetProtocolConfigResponse, error) {
	return suigraphql.GetProtocolConfig(
		ctx,
		c.client,
		protocolVersion,
	)
}

func (c *GraphQLClient) GetReferenceGasPrice(
	ctx context.Context,
) (*suigraphql.GetReferenceGasPriceResponse, error) {
	return suigraphql.GetReferenceGasPrice(
		ctx,
		c.client,
	)
}

func (c *GraphQLClient) GetStakes(
	ctx context.Context,
	owner sui.Address,
	limit *int,
	cursor *string,
) (*suigraphql.GetStakesResponse, error) {
	return suigraphql.GetStakes(
		ctx,
		c.client,
		owner,
		limit,
		cursor,
	)
}

func (c *GraphQLClient) GetStakesByIds(
	ctx context.Context,
	ids []sui.Address,
	limit *int,
	cursor *string,
) (*suigraphql.GetStakesByIdsResponse, error) {
	return suigraphql.GetStakesByIds(
		ctx,
		c.client,
		ids,
		limit,
		cursor,
	)
}

func (c *GraphQLClient) GetTotalSupply(
	ctx context.Context,
	coinType string,
) (*suigraphql.GetTotalSupplyResponse, error) {
	return suigraphql.GetTotalSupply(
		ctx,
		c.client,
		coinType,
	)
}

func (c *GraphQLClient) GetTotalTransactionBlocks(
	ctx context.Context,
) (*suigraphql.GetTotalTransactionBlocksResponse, error) {
	return suigraphql.GetTotalTransactionBlocks(
		ctx,
		c.client,
	)
}

func (c *GraphQLClient) GetTransactionBlock(
	ctx context.Context,
	digest string,
	showBalanceChanges *bool,
	showEffects *bool,
	showRawEffects *bool,
	showEvents *bool,
	showInput *bool,
	showObjectChanges *bool,
	showRawInput *bool,
) (*suigraphql.GetTransactionBlockResponse, error) {
	return suigraphql.GetTransactionBlock(
		ctx,
		c.client,
		digest,
		showBalanceChanges,
		showEffects,
		showRawEffects,
		showEvents,
		showInput,
		showObjectChanges,
		showRawInput,
	)
}

func (c *GraphQLClient) GetTypeLayout(
	ctx context.Context,
	targetType string,
) (*suigraphql.GetTypeLayoutResponse, error) {
	return suigraphql.GetTypeLayout(
		ctx,
		c.client,
		targetType,
	)
}

func (c *GraphQLClient) GetValidatorsApy(
	ctx context.Context,
) (*suigraphql.GetValidatorsApyResponse, error) {
	return suigraphql.GetValidatorsApy(
		ctx,
		c.client,
	)
}

func (c *GraphQLClient) MultiGetObjects(
	ctx context.Context,
	ids []sui.Address,
	limit *int,
	cursor *string,
	showBcs *bool,
	showContent *bool,
	showDisplay *bool,
	showType *bool,
	showOwner *bool,
	showPreviousTransaction *bool,
	showStorageRebate *bool,
) (*suigraphql.MultiGetObjectsResponse, error) {
	return suigraphql.MultiGetObjects(
		ctx,
		c.client,
		ids,
		limit,
		cursor,
		showBcs,
		showContent,
		showDisplay,
		showType,
		showOwner,
		showPreviousTransaction,
		showStorageRebate,
	)
}

func (c *GraphQLClient) MultiGetTransactionBlocks(
	ctx context.Context,
	digests []string,
	limit *int,
	cursor *string,
	showBalanceChanges *bool,
	showEffects *bool,
	showRawEffects *bool,
	showEvents *bool,
	showInput *bool,
	showObjectChanges *bool,
	showRawInput *bool,
) (*suigraphql.MultiGetTransactionBlocksResponse, error) {
	return suigraphql.MultiGetTransactionBlocks(
		ctx,
		c.client,
		digests,
		limit,
		cursor,
		showBalanceChanges,
		showEffects,
		showRawEffects,
		showEvents,
		showInput,
		showObjectChanges,
		showRawInput,
	)
}

func (c *GraphQLClient) PaginateCheckpointTransactionBlocks(
	ctx context.Context,
	id *suigraphql.CheckpointId,
	after *string,
) (*suigraphql.PaginateCheckpointTransactionBlocksResponse, error) {
	return suigraphql.PaginateCheckpointTransactionBlocks(
		ctx,
		c.client,
		id,
		after,
	)
}

func (c *GraphQLClient) PaginateEpochValidators(
	ctx context.Context,
	id uint64,
	after *string,
) (*suigraphql.PaginateEpochValidatorsResponse, error) {
	return suigraphql.PaginateEpochValidators(
		ctx,
		c.client,
		id,
		after,
	)
}

func (c *GraphQLClient) PaginateMoveModuleLists(
	ctx context.Context,
	packageId sui.Address,
	module string,
	hasMoreFriends bool,
	hasMoreStructs bool,
	hasMoreFunctions bool,
	hasMoreEnums bool,
	afterFriends *string,
	afterStructs *string,
	afterFunctions *string,
	afterEnums *string,
) (*suigraphql.PaginateMoveModuleListsResponse, error) {
	return suigraphql.PaginateMoveModuleLists(
		ctx,
		c.client,
		packageId,
		module,
		hasMoreFriends,
		hasMoreStructs,
		hasMoreFunctions,
		hasMoreEnums,
		afterFriends,
		afterStructs,
		afterFunctions,
		afterEnums,
	)
}

func (c *GraphQLClient) PaginateTransactionBlockLists(
	ctx context.Context,
	digest string,
	hasMoreEvents bool,
	hasMoreBalanceChanges bool,
	hasMoreObjectChanges bool,
	afterEvents string,
	afterBalanceChanges string,
	afterObjectChanges string,
) (*suigraphql.PaginateTransactionBlockListsResponse, error) {
	return suigraphql.PaginateTransactionBlockLists(
		ctx,
		c.client,
		digest,
		hasMoreEvents,
		hasMoreBalanceChanges,
		hasMoreObjectChanges,
		afterEvents,
		afterBalanceChanges,
		afterObjectChanges,
	)
}

func (c *GraphQLClient) QueryEvents(
	ctx context.Context,
	filter suigraphql.EventFilter,
	before *string,
	after *string,
	first *int,
	last *int,
) (*suigraphql.QueryEventsResponse, error) {
	return suigraphql.QueryEvents(
		ctx,
		c.client,
		filter,
		before,
		after,
		first,
		last,
	)
}

func (c *GraphQLClient) QueryTransactionBlocks(
	ctx context.Context,
	first *int,
	last *int,
	before *string,
	after *string,
	showBalanceChanges *bool,
	showEffects *bool,
	showRawEffects *bool,
	showEvents *bool,
	showInput *bool,
	showObjectChanges *bool,
	showRawInput *bool,
	filter *suigraphql.TransactionBlockFilter,
) (*suigraphql.QueryTransactionBlocksResponse, error) {
	return suigraphql.QueryTransactionBlocks(
		ctx,
		c.client,
		first,
		last,
		before,
		after,
		showBalanceChanges,
		showEffects,
		showRawEffects,
		showEvents,
		showInput,
		showObjectChanges,
		showRawInput,
		filter,
	)
}

func (c *GraphQLClient) ResolveNameServiceAddress(
	ctx context.Context,
	domain *string,
) (*suigraphql.ResolveNameServiceAddressResponse, error) {
	return suigraphql.ResolveNameServiceAddress(
		ctx,
		c.client,
		domain,
	)
}

func (c *GraphQLClient) ResolveNameServiceNames(
	ctx context.Context,
	address sui.Address,
	limit *int,
	cursor *string,
) (*suigraphql.ResolveNameServiceNamesResponse, error) {
	return suigraphql.ResolveNameServiceNames(
		ctx,
		c.client,
		address,
		limit,
		cursor,
	)
}

func (c *GraphQLClient) TryGetPastObject(
	ctx context.Context,
	id sui.Address,
	version *uint64,
	showBcs *bool,
	showOwner *bool,
	showPreviousTransaction *bool,
	showContent *bool,
	showDisplay *bool,
	showType *bool,
	showStorageRebate *bool,
) (*suigraphql.TryGetPastObjectResponse, error) {
	return suigraphql.TryGetPastObject(
		ctx,
		c.client,
		id,
		version,
		showBcs,
		showOwner,
		showPreviousTransaction,
		showContent,
		showDisplay,
		showType,
		showStorageRebate,
	)
}
