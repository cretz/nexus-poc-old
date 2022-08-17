package temporalnexus

import (
	"net/http"

	"github.com/cretz/nexus-poc/sdk/go/nexus/backend/backendpb"
	"go.temporal.io/sdk/workflow"
)

type WorkflowClientOptions struct {
	Endpoint string
}

type WorkflowClient struct {
}

func NewWorkflowClient(options WorkflowClientOptions) *WorkflowClient {
	panic("TODO")
}

type WorkflowStartALOOptions struct {
	Name string
	ID   string
	Arg  interface{}
}

type WorkflowCallSyncOptions struct {
	Name string
	// TODO(cretz): Optional, and delivered as metadata to the call maybe?
	ALOID string
	Arg   interface{}
}

type WorkflowCallOptions struct{}

func (*WorkflowClient) StartALO(ctx workflow.Context, options WorkflowStartALOOptions) (*WorkflowALOHandle, error) {
	panic("TODO")
}

func (*WorkflowClient) CallSync(ctx workflow.Context, options WorkflowCallSyncOptions) (*WorkflowSyncResult, error) {
	panic("TODO")
}

func (*WorkflowClient) CallHTTP(
	ctx workflow.Context,
	req *http.Request,
	options WorkflowCallOptions,
) (*http.Response, error) {
	panic("TODO")
}

func (*WorkflowClient) Call(
	ctx workflow.Context,
	req *backendpb.CallRequest,
	options WorkflowCallOptions,
) (*backendpb.CallResponse, error) {
	panic("TODO")
}

type WorkflowALOHandle struct {
	StartInfo *backendpb.AloInfo
	// TODO(cretz): This could embed workflow.Future (and maybe be an interface)
	// if/when we have the ability to _push_ the result
}

func (*WorkflowALOHandle) LoadInfo(ctx workflow.Context) (*backendpb.AloInfo, error) {
	panic("TODO")
}

func (*WorkflowALOHandle) WaitForResult(ctx workflow.Context) (*WorkflowALOResult, error) {
	panic("TODO")
}

type WorkflowALOResult struct {
	// TODO(cretz): The raw result, metadata, etc
}

func (*WorkflowALOResult) Get(valuePtr interface{}) error {
	panic("TODO")
}

type WorkflowSyncResult struct {
	// TODO(cretz): The raw result, metadata, etc
}

func (*WorkflowSyncResult) Get(valuePtr interface{}) error {
	panic("TODO")
}
