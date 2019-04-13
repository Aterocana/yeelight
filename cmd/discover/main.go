package main

import (
	"fmt"
	"log"
	"os"
	"time"
	"yeelight"
)

func main() {
	logger := log.New(os.Stdout, "", log.Ltime)
	discovery := yeelight.NewDiscoveryService()

	errc := discovery.GetErrors()
	devc := discovery.GetDiscoveredDevices()

	go func(c <-chan error) {
		for err := range c {
			fmt.Println(err)
		}
	}(errc)

	go func(c <-chan *yeelight.YeeLight) {
		for d := range c {
			logger.Printf("%v\n", d)
		}
	}(devc)

	if err := discovery.Open(); err != nil {
		log.Fatalf("%+v\n", err)
	}

	if err := discovery.DiscoveryRequest(); err != nil {
		log.Fatalf("%+v\n", err)
	}

	<-time.After(30 * time.Second)
}
