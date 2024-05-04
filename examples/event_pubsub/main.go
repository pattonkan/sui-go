package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/howjmay/go-sui-sdk/sui"
	"github.com/howjmay/go-sui-sdk/sui/conn"
	"github.com/howjmay/go-sui-sdk/sui_types"

	"github.com/howjmay/go-sui-sdk/examples/event_pubsub/lib"
)

func main() {
	done := make(chan os.Signal, 1)
	signal.Notify(done, syscall.SIGINT, syscall.SIGTERM)

	api := sui.NewSuiClient(conn.Dial(conn.TestnetEndpointUrl))
	sender, err := sui_types.NewAccountWithMnemonic(sui_types.TEST_MNEMONIC)
	if err != nil {
		log.Panic(err)
	}
	digest, err := sui.RequestFundFromFaucet(sender.Address, conn.TestnetFaucetUrl)
	if err != nil {
		log.Panic(err)
	}
	log.Println("digest: ", digest)

	packageID, err := sui_types.NewPackageIDFromHex("")
	if err != nil {
		log.Panic(err)
	}

	log.Println("sender: ", sender.Address)
	publisher := lib.NewPublisher(api, sender)
	subscriber := lib.NewSubscriber(api)

	go func() {
		for {
			publisher.PublishEvents(context.Background(), packageID)
		}
	}()

	go func() {
		for {
			subscriber.SubscribeEvents(context.Background(), packageID)
		}
	}()

	<-done
}