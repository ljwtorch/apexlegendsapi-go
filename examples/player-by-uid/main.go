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

	response, err := client.PlayerByUID(ctx, apex.PlayerByUIDRequest{
		UID:      "1015157856349",
		Platform: apex.PlatformPC,
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(response.Body))
}
