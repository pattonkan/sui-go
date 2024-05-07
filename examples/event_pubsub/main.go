package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/howjmay/sui-go/sui"
	"github.com/howjmay/sui-go/sui/conn"
	"github.com/howjmay/sui-go/sui_signer"
	"github.com/howjmay/sui-go/sui_types"

	"github.com/howjmay/sui-go/examples/event_pubsub/lib"
)

func main() {
	done := make(chan os.Signal, 1)
	signal.Notify(done, syscall.SIGINT, syscall.SIGTERM)

	api := sui.NewSuiClient(conn.TestnetEndpointUrl)
	sender, err := sui_signer.NewSignerWithMnemonic(sui_signer.TEST_MNEMONIC)
	if err != nil {
		log.Panic(err)
	}
	digest, err := sui.RequestFundFromFaucet(sender.Address.String(), conn.TestnetFaucetUrl)
	if err != nil {
		log.Panic(err)
	}
	log.Println("digest: ", digest)

	packageID, err := sui_types.PackageIDFromHex("")
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
			subscriber.SubscribeEvent(context.Background(), packageID)
		}
	}()

	<-done
}
