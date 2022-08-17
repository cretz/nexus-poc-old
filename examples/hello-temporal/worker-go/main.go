package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/cretz/nexus-poc/sdk/go/nexus/backend/client"
	"github.com/cretz/nexus-poc/sdk/go/nexus/backend/worker"
	"github.com/cretz/nexus-poc/sdk/temporal-go/temporalnexus"
	temporalclient "go.temporal.io/sdk/client"
	temporalworker "go.temporal.io/sdk/worker"
	"go.temporal.io/sdk/workflow"
)

func SayHelloWorkflow(ctx workflow.Context, name string) (string, error) {
	greeting := fmt.Sprintf("Hello, %v!", name)
	err := workflow.SetQueryHandler(ctx, "get-greeting", func() (string, error) { return greeting, nil })
	if err != nil {
		return "", err
	}
	// Wait for signals to update the prefix or to finish the workflow
	sel := workflow.NewSelector(ctx)
	done := false
	sel.AddReceive(workflow.GetSignalChannel(ctx, "update-prefix"), func(c workflow.ReceiveChannel, more bool) {
		var prefix string
		c.Receive(ctx, &prefix)
		greeting = fmt.Sprintf("%v, %v!", prefix, name)
	})
	sel.AddReceive(workflow.GetSignalChannel(ctx, "finish-workflow"), func(c workflow.ReceiveChannel, more bool) {
		c.Receive(ctx, nil)
		done = true
	})
	sel.AddReceive(ctx.Done(), func(c workflow.ReceiveChannel, more bool) {
		err = ctx.Err()
	})
	for !done && err == nil {
		sel.Select(ctx)
	}
	return greeting, err
}

func main() {
	ctx := context.TODO()

	// Create client
	client, err := client.Dial(ctx, client.Options{Target: "nexus-backend.example.com"})
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	// Build worker with Temporal pieces
	temporalClient, err := temporalclient.Dial(temporalclient.Options{HostPort: "nexus-backend.example.com:7233"})
	if err != nil {
		log.Fatal(err)
	}
	defer temporalClient.Close()
	builder, err := temporalnexus.NewWorkerBuilder(temporalnexus.WorkerBuilderOptions{Client: temporalClient})
	if err != nil {
		log.Fatal(err)
	}
	// Add workflow
	err = builder.AddWorkflowALOHandler("greeting", temporalnexus.WorkflowALOHandlerOptions{
		Workflow: SayHelloWorkflow,
		Options:  temporalclient.StartWorkflowOptions{TaskQueue: "my-task-queue"},
	})
	if err != nil {
		log.Fatal(err)
	}
	// Add queries and signals
	// TODO(cretz): Would rather "greeting/{id}/get" as a "method", but can't due
	// to no hierarchies allowed
	err = builder.AddQuerySyncHandler("greeting/get", temporalnexus.QuerySyncHandlerOptions{
		Query: "get-greeting",
	})
	if err != nil {
		log.Fatal(err)
	}
	err = builder.AddSignalSyncHandler("greeting/update-prefix", temporalnexus.SignalSyncHandlerOptions{
		Signal: "update-prefix",
	})
	if err != nil {
		log.Fatal(err)
	}
	err = builder.AddSignalSyncHandler("greeting/finish", temporalnexus.SignalSyncHandlerOptions{
		Signal: "finish-workflow",
	})
	if err != nil {
		log.Fatal(err)
	}

	// Start Temporal worker
	temporalWorker := temporalworker.New(temporalClient, "my-task-queue", temporalworker.Options{})
	temporalWorker.RegisterWorkflow(SayHelloWorkflow)
	if err := temporalWorker.Start(); err != nil {
		log.Fatal(err)
	}
	defer temporalWorker.Stop()

	// Start Nexus worker
	worker, err := worker.New(worker.Options{Client: client, Service: "my-service", HTTPHandler: builder.BuildHTTPHandler()})
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
