package worker

import (
	"context"
	"net/http"

	"github.com/cretz/nexus-poc/sdk/go/nexus/backend/client"
)

type Options struct {
	Client      *client.Client
	Service     string
	HTTPHandler http.Handler
}

type Worker struct {
}

func New(Options) (*Worker, error) {
	panic("TODO")
}

func (*Worker) Start(context.Context) error {
	panic("TODO")
}

func (*Worker) Stop() error {
	panic("TODO")
}
