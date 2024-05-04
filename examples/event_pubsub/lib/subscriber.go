package lib

import (
	"context"

	"github.com/howjmay/go-sui-sdk/sui"
	"github.com/howjmay/go-sui-sdk/sui_types"
)

type Subscriber struct {
	client *sui.ImplSuiAPI
	// *account.Account
}

func NewSubscriber(client *sui.ImplSuiAPI) *Subscriber {
	return &Subscriber{client: client}
}

func (s *Subscriber) SubscribeEvents(ctx context.Context, packageID *sui_types.PackageID) {

}
