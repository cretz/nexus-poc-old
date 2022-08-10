# Nexus

Nexus is a system for implementing and invoking remote calls. Implementers can handle calls via task handlers. Callers
use traditional protocols to make clients and calls.

This POC is currently under development.

## Components

This repository has the following components:

* [api/](api) - Nexus type definitions and internal gRPC API
* [nexus-cli/](nexus-cli) - The Nexus CLI (In Go, binary called `nexus`, uses the Go SDK)
* [sdk/](sdk) - Language-specific SDKs

## Walkthrough: Creating a low-level hello world HTTP service

This walkthrough assumes a Nexus backend is running at https://nexus-backend.example.com

### Registering a service

To create a Nexus service, you must first register it with the backend. This can be done individually via:

    nexus service register --service-name my-service --service-description "My service" --backend nexus-backend.example.com

But this can also be done with a services file. For example, given the `services.yaml` below:

```yaml
services:
  - name: my-service
    description: My service
```

Running:

    nexus service register --file services.yaml --backend nexus-backend.example.com

Will register all services in the file. To _replace_ all existing services (i.e. deleting ones not in the file), you can
use `--replace`.

Now that the service is registered with the backend, let's create a worker.

### Implementing a worker

To handle service, a worker must be created to handle the calls. This can be done in Go like so:

```go
package main

import (
  "context"
  "encoding/json"
  "fmt"
  "log"
  "net/http"

  "github.com/cretz/nexus-poc/sdk/go/nexus/backend/client"
  "github.com/cretz/nexus-poc/sdk/go/nexus/backend/worker"
)

func HandleMyServiceRequest(ctx context.Context, req *worker.CallRequest) *worker.CallResponse {
  // Validate HTTP
  if req.HTTP == nil {
    return worker.ErrInvalidProtocol
  } else if req.HTTP.RelativePath != "/greeting" {
    // Note the relative path has the service stripped
    return worker.NewHTTPError(http.StatusNotFound)
  } else if req.HTTP.Request.Method != "GET" {
    return worker.NewHTTPError(http.MethodNotAllowed)
  } else if req.HTTP.Request.Header.Get("Content-Type") != "application/json" {
    return worker.NewHTTPErrorf(http.StatusBadRequest, "Invalid input")
  }
  
  // Read body
  reqBody := struct{ Name string `json:"name"` }{}
  b, err := io.ReadAll(req.HTTP.Request.Body)
  if closeErr := req.HTTP.Request.Body.Close(); closeErr != nil && err == nil {
    err = closeErr
  }
  if err == nil {
    err = json.Unmarshal(b, &reqBody)
  }
  if err != nil {
    return worker.NewHTTPErrorf(http.StatusBadRequest, "Invalid input: %w", err)
  }

  // Respond with greeting JSON body
  b, _ := json.Marshal(map[string]string{"greeting": fmt.Sprintf("Hello, %v!", reqBody.Name)})
  return worker.NewHTTPOKResponse(b, "Content-Type", "application/json")
}

func main() {
  ctx := context.TODO()

  // Create client
  client, err := client.Dial(ctx, client.Options{Target:"nexus-backend.example.com"})
  if err != nil {
    log.Fatal(err)
  }
  defer client.Close()

  // Start worker
  worker, err := worker.New(worker.Options{
    ServiceHandlers: map[string]worker.NexusHandler{
      "my-service": HandleMyServiceRequest,
    },
  })
  if err != nil {
    log.Fatal(err)
  } else if err = worker.Start(ctx); err != nil {
    log.Fatal(err)
  }
  defer worker.Stop()

  // Wait for completion
  log.Print("Worker started, ctrl+c to stop")
  signalCh := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
  <-signalCh
}
```

Running this will start a worker to handle all calls to the `my-service` service.

### Running a frontend

To actually host this at HTTP, we have to start a frontend server. This can be done via:

    nexus frontend start

But the frontend will be serving nothing but a localhost management server. This management server could be used to
dynamically set the config, or the frontend can be started with a config. Given the following `config.yaml` file:

```yaml
frontend:
  bindings:
    # A name can be given to these to reference in endpoint, but if only a
    # single one per protocol is present, it is assumed to be named "default"
    - http: 0.0.0.0:8080
  endpoints:
    - http: default
      backend: nexus-backend.example.com
      service: my-service
```

This can of course have multiple endpoints, even on the same IP/port but with different paths. Now starting with:

    nexus frontend start --config config.yaml

Will start an HTTP server that can accept service calls.

### Invoking

Invocation is straightforward HTTP. Can simply cURL:

    curl -X POST http://localhost:8080/my-service/greeting \
        -H 'Content-Type: application/json' \
        -d '{"name":"Nexus"}'

This will return:

    {"greeting":"Hello, Nexus!"}

## Walkthrough: Creating a hello world gRPC service

TODO

## Walkthrough: Creating a low-level service with an ALO response

TODO

## Walkthrough: Creating a high-level HTTP and gRPC service

TODO

## Walkthrough: Creating and calling Temporal-backed ALOs

TODO

## Walkthrough: Creating a high-level HTTP and gRPC service in Java

TODO