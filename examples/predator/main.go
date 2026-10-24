package main

import (
	"context"
	"fmt"
	"log"
	"time"

	apex "github.com/ljwtorch/apexlegendsapi-go"
)

func main() {
	client, err := apex.NewClient("85427f5704bb1e9293ca6c8e40257d3f")
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	response, err := client.Predator(ctx)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(response.Body))
}
