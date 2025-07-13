package conn_test

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/Khan/genqlient/graphql"
	"github.com/stretchr/testify/require"

	"github.com/pattonkan/sui-go/sui"
	"github.com/pattonkan/sui-go/suiclient/conn"
)

func TestGraphQL(t *testing.T) {
	client := conn.NewGraphQLClient(conn.TestnetGraphQLEndpointUrl)

	t.Run("Standard API Call", func(t *testing.T) {
		addr, err := sui.AddressFromHex("0x7a89979774c55814f41fc1e3354e2ba38d3d62096d469d86b3132e947de1e8da")
		require.NoError(t, err)
		resp, err := client.GetAllBalances(context.TODO(), *addr, nil, nil)
		require.NoError(t, err)

		fmt.Println("All Balances:", resp.Address.Balances)
	})

	t.Run("Custom query", func(t *testing.T) {
		q := `
query GetAllBalances($owner: SuiAddress!, $limit: Int, $cursor: String) {
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
}`
		b, err := client.Query(q, map[string]interface{}{
			"owner": "0xe25afa59deccfec819aaa67bf14f049982d2a1ca87c49c8614da5ea2dc438f72",
		})
		require.NoError(t, err)
		fmt.Println("raw bytes:", string(b))

		var resp graphql.Response
		err = json.Unmarshal(b, &resp)
		require.NoError(t, err)
		fmt.Println("unmarshalled:", resp)
	})
}
