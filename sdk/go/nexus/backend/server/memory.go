package server

import (
	"context"
	"sort"
	"sync"

	"github.com/cretz/nexus-poc/sdk/go/nexus/backend/backendpb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type service struct {
	*backendpb.Service
	reqCh     chan *backendpb.CallRequest
	deletedCh chan struct{}
}

type inMemoryServer struct {
	backendpb.UnimplementedBackendServiceServer
	services     map[string]*service
	servicesLock sync.RWMutex
	// Keyed by request ID
	callResponses     map[string]chan<- *backendpb.CallResponse
	callResponsesLock sync.RWMutex
}

func NewInMemoryServer() backendpb.BackendServiceServer {
	return &inMemoryServer{
		services:      map[string]*service{},
		callResponses: map[string]chan<- *backendpb.CallResponse{},
	}
}

func (i *inMemoryServer) UpdateServices(
	ctx context.Context,
	req *backendpb.UpdateServicesRequest,
) (*backendpb.UpdateServicesResponse, error) {
	// In-memory is cheap enough to lock this whole thing
	i.servicesLock.Lock()
	defer i.servicesLock.Unlock()

	// Remove services cannot be present if replace is present
	if len(req.RemoveServices) > 0 && req.Replace {
		return nil, status.Error(codes.InvalidArgument, "cannot have remove services and replace")
	}

	// Validate services
	seenNames := map[string]bool{}
	for _, service := range req.Services {
		if service.Name == "" {
			return nil, status.Error(codes.InvalidArgument, "all services must have names")
		} else if seenNames[service.Name] {
			return nil, status.Errorf(codes.InvalidArgument, "service name %q duplicated", service.Name)
		} else if !service.Http {
			// TODO(cretz):
			return nil, status.Error(codes.InvalidArgument, "all services must support HTTP at this time")
		}
		seenNames[service.Name] = true
	}

	// Validate remove services
	for _, service := range req.RemoveServices {
		if i.services[service] == nil {
			return nil, status.Errorf(codes.NotFound, "service %q not found for removal", service)
		} else if seenNames[service] {
			return nil, status.Errorf(codes.InvalidArgument, "duplicate service %q for removal", service)
		}
		seenNames[service] = true
	}

	// Apply updates now that there are no more errors that can occur
	if req.Replace {
		// Close/delete all services that are not in the request
		for name, service := range i.services {
			beingUpdated := false
			for _, updateService := range req.Services {
				if updateService.Name == name {
					beingUpdated = true
					break
				}
			}
			if !beingUpdated {
				close(service.deletedCh)
				delete(i.services, name)
			}
		}
	}
	// Close/delete specific services
	for _, serviceName := range req.RemoveServices {
		close(i.services[serviceName].deletedCh)
		delete(i.services, serviceName)
	}
	for _, serviceProto := range req.Services {
		// If the service exists, just update the proto
		if svc := i.services[serviceProto.Name]; svc != nil {
			svc.Service = serviceProto
		} else {
			// Create the service
			i.services[serviceProto.Name] = &service{
				Service:   serviceProto,
				reqCh:     make(chan *backendpb.CallRequest),
				deletedCh: make(chan struct{}),
			}
		}
	}
	return &backendpb.UpdateServicesResponse{Services: i.getServicesUnlocked(ctx)}, nil
}

func (i *inMemoryServer) GetServices(
	ctx context.Context,
	req *backendpb.GetServicesRequest,
) (*backendpb.GetServicesResponse, error) {
	i.servicesLock.RLock()
	defer i.servicesLock.RUnlock()
	return &backendpb.GetServicesResponse{Services: i.getServicesUnlocked(ctx)}, nil
}

func (i *inMemoryServer) getServicesUnlocked(ctx context.Context) []*backendpb.Service {
	ret := make([]*backendpb.Service, 0, len(i.services))
	for _, service := range i.services {
		ret = append(ret, service.Service)
	}
	sort.Slice(ret, func(i, j int) bool { return ret[i].Name < ret[j].Name })
	return ret
}

func (i *inMemoryServer) Call(
	ctx context.Context,
	req *backendpb.CallRequest,
) (*backendpb.CallResponse, error) {
	// Validate
	if req.RequestId == "" {
		return nil, status.Error(codes.InvalidArgument, "missing request ID")
	} else if req.GetHttp() == nil {
		// TODO(cretz):
		return nil, status.Error(codes.InvalidArgument, "all calls must be HTTP currently")
	}

	// Get channel to send request to
	i.servicesLock.RLock()
	service := i.services[req.Service]
	i.servicesLock.RUnlock()
	if service == nil {
		return nil, status.Error(codes.NotFound, "service not found")
	}

	// Add a channel for response and remove on complete
	respCh := make(chan *backendpb.CallResponse, 1)
	i.callResponsesLock.Lock()
	i.callResponses[req.RequestId] = respCh
	i.callResponsesLock.Unlock()
	defer func() {
		i.callResponsesLock.Lock()
		delete(i.callResponses, req.RequestId)
		i.callResponsesLock.Unlock()
	}()

	// Attempt send until context done
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case service.reqCh <- req:
	}

	// Wait for response until context done
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case resp := <-respCh:
		return resp, nil
	}
}

func (i *inMemoryServer) StreamTasks(srv backendpb.BackendService_StreamTasksServer) error {
	panic("TODO")
}
