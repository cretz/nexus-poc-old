package main

import (
	"context"
	"fmt"
	"log"

	"github.com/cretz/nexus-poc/sdk/go/nexus/backend/backendpb"
	"github.com/cretz/nexus-poc/sdk/go/nexus/frontend"
)

func main() {
	ctx := context.TODO()

	// Create Nexus client. Note, users could just use the normal Go HTTP client
	// if they wanted.
	client, err := frontend.Dial(ctx, frontend.ClientOptions{Endpoint: "http://localhost:8080"})
	if err != nil {
		log.Fatal(err)
	}

	// Start greeting ALO
	alo, err := client.StartALO(ctx, frontend.StartALOOptions{
		Name: "greeting",
		ID:   "my-id",
		Arg:  "Nexus",
	})
	if err != nil {
		log.Fatal(err)
	}

	// Load info and confirm it is running
	if info, err := alo.LoadInfo(ctx); err != nil {
		log.Fatal(err)
	} else if info.Status != backendpb.AloInfo_RUNNING {
		log.Fatal(fmt.Errorf("expected running, got status: %v", info.Status))
	}

	// Send a query and confirm it is the default greeting
	var resultStr string
	res, err := client.CallSync(ctx, frontend.CallSyncOptions{Name: "greeting/get", ALOID: "my-id"})
	if err != nil {
		log.Fatal(err)
	} else if err = res.Get(&resultStr); err != nil {
		log.Fatal(err)
	} else if resultStr != "Hello, Nexus!" {
		log.Fatal(fmt.Errorf("expected normal greeting, got: %q", resultStr))
	}

	// Send a signal to update the prefix, then another signal to finish the
	// workflow
	_, err = client.CallSync(ctx, frontend.CallSyncOptions{
		Name:  "greeting/update-prefix",
		ALOID: "my-id",
		Arg:   "Howdy",
	})
	if err != nil {
		log.Fatal(err)
	}
	_, err = client.CallSync(ctx, frontend.CallSyncOptions{Name: "greeting/finish", ALOID: "my-id"})
	if err != nil {
		log.Fatal(err)
	}

	// Now confirm the result
	if res, err := alo.WaitForResult(ctx); err != nil {
		log.Fatal(err)
	} else if err = res.Get(&resultStr); err != nil {
		log.Fatal(err)
	} else if resultStr != "Howdy, Nexus!" {
		log.Fatal(fmt.Errorf("expected normal greeting, got: %q", resultStr))
	}
}
