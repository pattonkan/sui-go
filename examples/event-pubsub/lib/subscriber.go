package lib

import (
	"context"
	"fmt"
	"log"

	"github.com/pattonkan/sui-go/sui"
	"github.com/pattonkan/sui-go/suiclient"
)

type Subscriber struct {
	client *suiclient.ClientImpl
	// *account.Account
}

func NewSubscriber(client *suiclient.ClientImpl) *Subscriber {
	return &Subscriber{client: client}
}

func (s *Subscriber) SubscribeEvent(ctx context.Context, packageId *sui.PackageId) {
	resultCh := make(chan suiclient.Event)
	err := s.client.SubscribeEvent(context.Background(), &suiclient.EventFilter{Package: packageId}, resultCh)
	if err != nil {
		log.Fatal(err)
	}

	for result := range resultCh {
		fmt.Println("result: ", result)
	}
}
