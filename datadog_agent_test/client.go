package main

import (
	"log"
	"time"

	"github.com/DataDog/datadog-go/statsd"
)

func main() {
	statsd, err := statsd.New("127.0.0.1:8125")
	if err != nil {
		log.Fatalf("Error creating statsd client: %v", err)
	}

	for true {
		statsd.Incr("example_metric.increment", []string{"environment:dev"}, 1)
		statsd.Decr("example_metric.decrement", []string{"environment:dev"}, 1)
		statsd.Count("example_metric.count", 2, []string{"environment:dev"}, 1)
		time.Sleep(10 * time.Second)

	}
}