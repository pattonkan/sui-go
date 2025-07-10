# GraphQL

The schema and query files we used are all from [`MystenLabs/ts-sdks`](https://github.com/MystenLabs/ts-sdks) and [`MystenLabs/sui-rust-sdk`](https://github.com/MystenLabs/sui-rust-sdk). Unfortunately, we need to modify them for adapt the limitation of Golang and `genqlient`. The following is going to explain how to generate the GraphQL interface from the schema and queries, and what are the changes we need to do.

## How To Generate GraphQL Client

For updating GraphQL's client, we need to run `genqlient` to generate the GraphQL interface.

To generate `genqlient` interface, go to `suiclient/suigraphql`, and run

```bash
$ genqlient
```

## Optional Arguments

Genqlient doesn't automatically make the optional arguments in GraphQL `query` and `mutation` into pointers. It is necessary to use pointers for optional arguments; otherwise, the server may encounter errors when processing incoming requests. To successfully generate optional arguments as pointers you should add `# @genqlient(pointer: true)` above the targeting argument. See the following example.

```graphql
query GetAllBalances(
	$owner: SuiAddress!,
	# @genqlient(pointer: true)
	$limit: Int,
	# @genqlient(pointer: true)
	$cursor: String,
) {
	address(address: $owner) {
		balances(first: $limit, after: $cursor) {
			pageInfo {
				hasNextPage
				endCursor
			}
			nodes {
				coinType {
					repr
				}
				coinObjectCount
				totalBalance
			}
		}
	}
}
```

## Function Exposure

The function names of each `query` and `mutation` in the generated file are exactly the same as their declaration in `queries/*.graphql`. To expose them, we need to ensure the leading character is capital case in `*.graphql`, so it can meet the syntax of Golang. 

## Duplicated Field Names in Queries and Schema

Duplicated field names are not allowed in `genqlient`, because Golang's limitation. To adapt the limitation, we need to add alias for the duplicated fields. See the following example, where `contents` are the fields that have a duplicated name.

### Before

```graphql
fragment RPC_MOVE_OBJECT_FIELDS on MoveObject {
	objectId: address
	bcs @include(if: $showBcs)
	contents @include(if: $showType) {
		type {
			repr
		}
	}
	hasPublicTransfer @include(if: $showContent)
	contents @include(if: $showContent) {
		data
		type {
			repr
			layout
			signature
		}
	}

	hasPublicTransfer @include(if: $showBcs)
	contents @include(if: $showBcs) {
		bcs
		type {
			repr
		}
	}

	owner @include(if: $showOwner) {
		...RPC_OBJECT_OWNER_FIELDS
	}
	previousTransactionBlock @include(if: $showPreviousTransaction) {
		digest
	}

	storageRebate @include(if: $showStorageRebate)
	digest
	version
	display @include(if: $showDisplay) {
		key
		value
		error
	}
}
```

### After

We add the alias `contents_type` and `contents_content`.

```graphql
fragment RPC_MOVE_OBJECT_FIELDS on MoveObject {
	objectId: address
	bcs @include(if: $showBcs)
	contents_type: contents @include(if: $showType) {
		type {
			repr
		}
	}
	hasPublicTransfer @include(if: $showContent)
	contents_content: contents @include(if: $showContent) {
		data
		type {
			repr
			layout
			signature
		}
	}

	hasPublicTransfer @include(if: $showBcs)
	contents @include(if: $showBcs) {
		bcs
		type {
			repr
		}
	}

	owner @include(if: $showOwner) {
		...RPC_OBJECT_OWNER_FIELDS
	}
	previousTransactionBlock @include(if: $showPreviousTransaction) {
		digest
	}

	storageRebate @include(if: $showStorageRebate)
	digest
	version
	display @include(if: $showDisplay) {
		key
		value
		error
	}
}
```
