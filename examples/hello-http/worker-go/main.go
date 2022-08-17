package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"

	"github.com/cretz/nexus-poc/sdk/go/nexus/backend/backendpb"
	"github.com/cretz/nexus-poc/sdk/go/nexus/backend/client"
	"github.com/cretz/nexus-poc/sdk/go/nexus/backend/worker"
)

func HandleGreeting2(context.Context, *backendpb.CallRequest) (*backendpb.CallResponse, error) {

	panic("TODO")
}

func HandleGreeting(resp http.ResponseWriter, req *http.Request) {
	// Check method and content type
	if req.Method != "GET" {
		http.Error(resp, "must be GET", http.StatusMethodNotAllowed)
		return
	} else if req.Header.Get("Content-Type") != "application/json" {
		http.Error(resp, "must be application/json content type", http.StatusBadRequest)
		return
	}

	// Read body
	reqBody := struct {
		Name string `json:"name"`
	}{}
	b, err := io.ReadAll(req.Body)
	if closeErr := req.Body.Close(); closeErr != nil && err == nil {
		err = closeErr
	}
	if err == nil {
		err = json.Unmarshal(b, &reqBody)
	}
	if err != nil {
		http.Error(resp, fmt.Sprintf("invalid request: %v", err), http.StatusBadRequest)
		return
	}

	// Respond with greeting JSON body
	b, _ = json.Marshal(map[string]string{"greeting": fmt.Sprintf("Hello, %v!", reqBody.Name)})
	resp.Header().Add("Content-Type", "application/json")
	resp.Write(b)
}

func main() {
	ctx := context.TODO()

	// Create client
	client, err := client.Dial(ctx, client.Options{Target: "nexus-backend.example.com"})
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	// Start worker
	mux := http.NewServeMux()
	mux.HandleFunc("/greeting", HandleGreeting)
	worker, err := worker.New(worker.Options{Client: client, Service: "my-service", HTTPHandler: mux})
	if err != nil {
		log.Fatal(err)
	} else if err = worker.Start(ctx); err != nil {
		log.Fatal(err)
	}
	defer worker.Stop()

	// Wait for completion
	log.Print("Worker started, ctrl+c to stop")
	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, os.Interrupt)
	<-signalCh
}
