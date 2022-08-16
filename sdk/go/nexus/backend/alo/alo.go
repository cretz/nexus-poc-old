package alo

import (
	"context"
	"errors"
	"io"
	"net/http"

	"github.com/cretz/nexus-poc/sdk/go/nexus/backend/backendpb"
)

// TODO(cretz): Custom HTTP error type
var ErrAlreadyExists = errors.New("already exists")
var ErrNotFound = errors.New("not found")
var ErrNotImplemented = errors.New("not implemented")

type ALOStartRequest struct {
	*backendpb.AloInfo
	Body io.Reader
}

type ALOHandler interface {
	Start(ctx context.Context, req *ALOStartRequest) error
	GetInfo(ctx context.Context, id string) (*backendpb.AloInfo, error)
	GetResult(ctx context.Context, id string) ([]byte, error)
	Cancel(ctx context.Context, id string) error
}

type HTTPOptions struct {
	Handler ALOHandler
}

func NewHTTPHandler(options HTTPOptions) (http.Handler, error) {
	panic("TODO")
}
