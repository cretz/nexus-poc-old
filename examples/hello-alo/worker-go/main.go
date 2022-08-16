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
	"sync"

	"github.com/cretz/nexus-poc/sdk/go/nexus/backend/alo"
	"github.com/cretz/nexus-poc/sdk/go/nexus/backend/backendpb"
	"github.com/cretz/nexus-poc/sdk/go/nexus/backend/client"
	"github.com/cretz/nexus-poc/sdk/go/nexus/backend/worker"
	"google.golang.org/protobuf/proto"
)

// Simple in-memory handler
type myALOHandler struct {
	alos     map[string]*myALO
	alosLock sync.RWMutex
}

type myALO struct {
	*backendpb.AloInfo
	result []byte
	done   chan struct{}
	lock   sync.RWMutex
}

func newMyALOHandler() *myALOHandler {
	return &myALOHandler{alos: map[string]*myALO{}}
}

func (m *myALOHandler) Start(ctx context.Context, req *alo.ALOStartRequest) error {
	m.alosLock.Lock()
	defer m.alosLock.Unlock()
	if m.alos[req.Id] != nil {
		return alo.ErrAlreadyExists
	}
	// Parse body
	reqBody := struct {
		Name string `json:"name"`
	}{}
	if b, err := io.ReadAll(req.Body); err != nil {
		return err
	} else if err = json.Unmarshal(b, &reqBody); err != nil {
		return err
	}
	// Run "long-running" thing in background
	myAlo := &myALO{AloInfo: req.AloInfo, done: make(chan struct{})}
	m.alos[req.Id] = myAlo
	go func() {
		myAlo.lock.Lock()
		defer myAlo.lock.Unlock()
		myAlo.result, _ = json.Marshal(map[string]string{"greeting": fmt.Sprintf("Hello, %v!", reqBody.Name)})
		close(myAlo.done)
		myAlo.Status = backendpb.AloInfo_COMPLETED
	}()
	return nil
}

func (m *myALOHandler) GetInfo(ctx context.Context, id string) (*backendpb.AloInfo, error) {
	m.alosLock.RLock()
	myAlo := m.alos[id]
	m.alosLock.RUnlock()
	if myAlo == nil {
		return nil, alo.ErrNotFound
	}
	myAlo.lock.RLock()
	defer myAlo.lock.RUnlock()
	return proto.Clone(myAlo.AloInfo).(*backendpb.AloInfo), nil
}

func (m *myALOHandler) GetResult(ctx context.Context, id string) ([]byte, error) {
	m.alosLock.RLock()
	myAlo := m.alos[id]
	m.alosLock.RUnlock()
	if myAlo == nil {
		return nil, alo.ErrNotFound
	}
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case <-myAlo.done:
		myAlo.lock.RLock()
		defer myAlo.lock.RUnlock()
		// TODO(cretz): We need some way to specify response metadata
		return append(make([]byte, 0, len(myAlo.result)), myAlo.result...), nil
	}
}

func (m *myALOHandler) Cancel(ctx context.Context, id string) error {
	return alo.ErrNotImplemented
}

func main() {
	ctx := context.TODO()

	// Create client
	client, err := client.Dial(ctx, client.Options{Target: "nexus-backend.example.com"})
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	// Prepare mux
	mux := http.NewServeMux()
	greetingHandler, err := alo.NewHTTPHandler(alo.HTTPOptions{Handler: newMyALOHandler()})
	if err != nil {
		log.Fatal(err)
	}
	mux.Handle("/greeting", greetingHandler)

	// Start worker
	worker, err := worker.New(worker.Options{Client: client, HTTPHandler: mux})
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
