package frontend

import (
	"context"
	"net/http"

	"github.com/cretz/nexus-poc/sdk/go/nexus/backend/backendpb"
)

type Client struct {
}

type ClientOptions struct {
	Endpoint string
}

func Dial(ctx context.Context, options ClientOptions) (*Client, error) {
	panic("TODO")
}

type StartALOOptions struct {
	Name string
	ID   string
	Arg  interface{}
}

type CallSyncOptions struct {
	Name string
	// TODO(cretz): Optional, and delivered as metadata to the call maybe?
	ALOID string
	Arg   interface{}
}

type CallOptions struct{}

func (*Client) StartALO(ctx context.Context, options StartALOOptions) (*ALOHandle, error) {
	panic("TODO")
}

func (*Client) CallSync(ctx context.Context, options CallSyncOptions) (*SyncResult, error) {
	panic("TODO")
}

func (*Client) CallHTTP(
	ctx context.Context,
	req *http.Request,
	options CallOptions,
) (*http.Response, error) {
	panic("TODO")
}

func (*Client) Call(
	ctx context.Context,
	req *backendpb.CallRequest,
	options CallOptions,
) (*backendpb.CallResponse, error) {
	panic("TODO")
}

type ALOHandle struct {
	StartInfo *backendpb.AloInfo
	// TODO(cretz): This could embed workflow.Future (and maybe be an interface)
	// if/when we have the ability to _push_ the result
}

func (*ALOHandle) LoadInfo(ctx context.Context) (*backendpb.AloInfo, error) {
	panic("TODO")
}

func (*ALOHandle) WaitForResult(ctx context.Context) (*ALOResult, error) {
	panic("TODO")
}

type ALOResult struct {
	// TODO(cretz): The raw result, metadata, etc
}

func (*ALOResult) Get(valuePtr interface{}) error {
	panic("TODO")
}

type SyncResult struct {
	// TODO(cretz): The raw result, metadata, etc
}

func (*SyncResult) Get(valuePtr interface{}) error {
	panic("TODO")
}
