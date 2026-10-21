package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	apex "github.com/ljwtorch/apexlegendsapi-go"
)

func main() {
	client, err := apex.NewClient(os.Getenv("APEX_API_KEY"))
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
