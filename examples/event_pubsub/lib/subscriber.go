package lib

import (
	"context"
	"fmt"
	"log"

	"github.com/howjmay/go-sui-sdk/sui"
	"github.com/howjmay/go-sui-sdk/sui_types"
	"github.com/howjmay/go-sui-sdk/types"
)

type Subscriber struct {
	client *sui.ImplSuiAPI
	// *account.Account
}

func NewSubscriber(client *sui.ImplSuiAPI) *Subscriber {
	return &Subscriber{client: client}
}

func (s *Subscriber) SubscribeEvent(ctx context.Context, packageID *sui_types.PackageID) {
	resultCh := make(chan types.SuiEvent)
	err := s.client.SubscribeEvent(context.Background(), types.EventFilter{Package: packageID}, resultCh)
	if err != nil {
		log.Fatal(err)
	}

	for result := range resultCh {
		fmt.Println("result: ", result)
	}
}
