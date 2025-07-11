package conn_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/pattonkan/sui-go/sui"
	"github.com/pattonkan/sui-go/suiclient/conn"
	"github.com/stretchr/testify/require"
)

func TestGraphQL(t *testing.T) {
	client := conn.NewGraphQLClient(conn.MainnetGraphQLEndpointUrl)

	addr, err := sui.AddressFromHex("0x7a89979774c55814f41fc1e3354e2ba38d3d62096d469d86b3132e947de1e8da")
	require.NoError(t, err)
	resp, err := client.GetAllBalances(context.TODO(), *addr, nil, nil)
	require.NoError(t, err)

	fmt.Println("All Balances:", resp.Address.Balances)
}
