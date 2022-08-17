package temporalnexus

import (
	"context"
	"net/http"

	"github.com/cretz/nexus-poc/sdk/go/nexus/backend/alo"
	"github.com/cretz/nexus-poc/sdk/go/nexus/backend/backendpb"
	"go.temporal.io/sdk/client"
)

type WorkerBuilder struct {
	client       client.Client
	aloHandlers  map[string]alo.Handler
	syncHandlers map[string]func(context.Context, *backendpb.CallRequest) (*backendpb.CallResponse, error)
	httpHandlers map[string]http.Handler
}

type WorkerBuilderOptions struct {
	Client client.Client
}

func NewWorkerBuilder(options WorkerBuilderOptions) (*WorkerBuilder, error) {
	panic("TODO")
}

func (w *WorkerBuilder) AddALOHandler(name string, handler alo.Handler) error {
	panic("TODO")
}

func (w *WorkerBuilder) AddWorkflowALOHandler(name string, options WorkflowALOHandlerOptions) error {
	panic("TODO")
}

func (w *WorkerBuilder) AddSyncHandler(
	name string,
	handler func(context.Context, *backendpb.CallRequest) (*backendpb.CallResponse, error),
) error {
	panic("TODO")
}

func (w *WorkerBuilder) AddQuerySyncHandler(name string, options QuerySyncHandlerOptions) error {
	panic("TODO")
}

func (w *WorkerBuilder) AddSignalSyncHandler(name string, options SignalSyncHandlerOptions) error {
	panic("TODO")
}

func (w *WorkerBuilder) AddHTTPHandler(name string, handler http.Handler) error {
	panic("TODO")
}

func (w *WorkerBuilder) BuildHTTPHandler() http.Handler {
	panic("TODO")
}

type WorkflowALOHandlerOptions struct {
	Workflow interface{}
	// Should not have ID
	Options client.StartWorkflowOptions
	// TODO(cretz): Lots more customization of search attributes, input/output
	// types, etc
}

type workflowALOHandler struct {
	WorkflowALOHandlerOptions
}

func (w *workflowALOHandler) Start(ctx context.Context, req *alo.StartRequest) error {
	panic("TODO")
}

func (w *workflowALOHandler) GetInfo(ctx context.Context, id string) (*backendpb.AloInfo, error) {
	panic("TODO")
}

func (w *workflowALOHandler) GetResult(ctx context.Context, id string) ([]byte, error) {
	panic("TODO")
}

func (w *workflowALOHandler) Cancel(ctx context.Context, id string) error {
	panic("TODO")
}

func StartALOWorkflow(
	ctx context.Context,
	client client.Client,
	req *alo.StartRequest,
	options client.StartWorkflowOptions,
	workflow interface{},
) error {
	panic("TODO")
}

func JSONArgsFromBody(req *alo.StartRequest) (interface{}, error) {
	panic("TODO")
}

type QuerySyncHandlerOptions struct {
	Query string
	// TODO(cretz): Lots more customization
}

type querySyncHandler struct {
}

type SignalSyncHandlerOptions struct {
	Signal string
	// TODO(cretz): Lots more customization
}

type signalSyncHandler struct {
}
