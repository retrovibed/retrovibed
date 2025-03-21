package main

import (
	"context"
	"fmt"
	"log"
	"math"

	"github.com/hashicorp/mdns"
	"github.com/retrovibed/retrovibed/internal/netx"
)

func StorageProfit(users int, factor float64) float64 {
	return (100 * (factor/(math.Log2(float64(users+1))) + (1 - factor)))
}

func Derp(users int, factor float64) {
	costper := StorageProfit(users, factor)
	log.Println(users, ":", costper, costper*float64(1000))
}

func main() {

	netx.HostIP("0.0.0.0")
	ctx, done := context.WithCancel(context.Background())
	// Make a channel for results and start listening
	entriesCh := make(chan *mdns.ServiceEntry, 4)
	q := mdns.DefaultParams("_shallows._udp")
	q.Entries = entriesCh

	go func() {
		defer close(entriesCh)
		defer done()
		// Start the lookup
		if err := mdns.QueryContext(ctx, q); err != nil {
			log.Fatalf("failed to look up service: %v", err)
		}
	}()

	for entry := range entriesCh {
		fmt.Printf("Got new entry: %v\n", entry)
	}

}
