# Nexus

Nexus is a system for implementing and invoking remote calls. Implementers can handle calls via task handlers. Callers
use traditional protocols to make clients and calls.

This POC is currently under development.

## Components

This repository has the following components:

* [api/](api) - Nexus type definitions and internal gRPC API
* [nexus-cli/](nexus-cli) - The Nexus CLI (In Go, binary called `nexus`, uses the Go SDK)
* [sdk/](sdk) - Language-specific SDKs

Note, this POC only deals in HTTP. gRPC and others may come later.

## ALOs over HTTP

ALOs are modelled on top of the protocol, not within it. For HTTP:

### Starting an ALO

`POST /my-service/<path>/<identifier>` that, on success, returns a `Content-Type` of `application/x-nexus-alo` with a
response body that is JSON formatted `nexus.backend.v1.AloInfo` object and the status code will be 201.

The request body can be in anything and the headers are turned into ALO metadata.

#### URL result callbacks

If the request header `Nexus-Result-Callback` is present with a URL (or multiple headers of such), those URLs are
invoked with the a request with `Content-Type` of `application/x-nexus-alo` and a body of `nexus.backend.v1.AloInfo`.
Recipients may then ask for the ALO result. TODO(cretz): Should we instead send the result/error in the callback?

### Getting ALO info

`GET /my-service/<path>/<identifier>` will, on success, return a `Content-Type` of `application/x-nexus-alo` and the
response body will be a JSON formatted `nexus.backend.v1.AloInfo` object and the status code will be 200.

### Getting an ALO result

`GET /my-service/<path>/<identifier>/result` will block and return a result when available (or error). The response body
can be anything that represents the success (200 status code) or failure.

### Cancelling an ALO

`POST /my-service/<path>/<identifier>/cancel` will cancel the ALO. There is no response on success, but the status code
is 202.

## Walkthrough: Creating a low-level hello world HTTP service

This walkthrough assumes a Nexus backend is running at https://nexus-backend.example.com

### Registering a service

To create a Nexus service, you must first register it with the backend. This can be done individually via:

    nexus service register --service-name my-service --service-description "My service" --service-http --backend nexus-backend.example.com

But this can also be done with a services file. For example, given the `services.yaml` below:

```yaml
services:
  - name: my-service
    description: My service
    http: true
```

Running:

    nexus service register --file services.yaml --backend nexus-backend.example.com

Will register all services in the file. To _replace_ all existing services (i.e. deleting ones not in the file), you can
use `--replace`.

Now that the service is registered with the backend, let's create a worker.

### Implementing a worker

To handle service, a worker must be created to handle the calls. See [this file](examples/hello-http/worker-go/main.go)
for a simple worker in Go. Running that file will start a worker to handle all calls to the `my-service` service.

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

## Walkthrough: Creating a low-level service with an ALO response

This relies on knowledge from the previous walkthrough.

### Registering a service

`services.yaml`:

```yaml
services:
  - name: my-service
    description: My service
    http: true
```

Run:

    nexus service register --file services.yaml --backend nexus-backend.example.com

### Implementing a worker

To handle service, a worker must be created to handle the ALO calls. See
[this file](examples/hello-alo/worker-go/main.go) for a simple ALO worker in Go. Running that file will start a worker
to handle all `greeting` calls to the `my-service` service as ALOs.

### Running a frontend

`config.yaml`:

```yaml
frontend:
  bindings:
    - http: 0.0.0.0:8080
  endpoints:
    - http: default
      backend: nexus-backend.example.com
      service: my-service
```

Run:

    nexus frontend start --config config.yaml

Frontend HTTP server now running

### Invoking

Invocation is straightforward HTTP. Can simply cURL:

    curl -X POST http://localhost:8080/my-service/greeting/my-alo-id \
        -H 'Content-Type: application/json' \
        -d '{"name":"Nexus"}'

This will return a 201 with a `Content-Type` of `application/x-nexus-alo` and a body of:

    {"id":"my-alo-id","status":"RUNNING"}

To check the status:

    curl http://localhost:8080/my-service/greeting/my-alo-id

This may return a 200 with a `Content-Type` of `application/x-nexus-alo` and a body of:

    {"id":"my-alo-id","status":"COMPLETED"}

To get the result:

    curl http://localhost:8080/my-service/greeting/my-alo-id/result

This may return a 200 with a body of:

    {"greeting":"Hello, Nexus!"}

## Walkthrough: Creating a high-level HTTP service

TODO

## Walkthrough: Creating and calling Temporal-backed ALOs

TODO

## Walkthrough: Creating a high-level HTTP in Java

TODO