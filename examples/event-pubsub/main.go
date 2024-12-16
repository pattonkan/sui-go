package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/pattonkan/sui-go/examples/event-pubsub/lib"
	"github.com/pattonkan/sui-go/sui"
	"github.com/pattonkan/sui-go/suiclient"
	"github.com/pattonkan/sui-go/suiclient/conn"
	"github.com/pattonkan/sui-go/suisigner"
)

func main() {
	done := make(chan os.Signal, 1)
	signal.Notify(done, syscall.SIGINT, syscall.SIGTERM)

	api := suiclient.NewClient(conn.TestnetEndpointUrl)
	sender, err := suisigner.NewSignerWithMnemonic(suisigner.TEST_MNEMONIC, suisigner.KeySchemeFlagDefault)
	if err != nil {
		log.Panic(err)
	}
	err = suiclient.RequestFundFromFaucet(sender.Address, conn.TestnetFaucetUrl)
	if err != nil {
		log.Panic(err)
	}

	packageId, err := sui.PackageIdFromHex("")
	if err != nil {
		log.Panic(err)
	}

	log.Println("sender: ", sender.Address)
	publisher := lib.NewPublisher(api, sender)
	subscriber := lib.NewSubscriber(api)

	go func() {
		for {
			publisher.PublishEvents(context.Background(), packageId)
		}
	}()

	go func() {
		for {
			subscriber.SubscribeEvent(context.Background(), packageId)
		}
	}()

	<-done
}
