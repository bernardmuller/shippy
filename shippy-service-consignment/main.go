package main

import (
	c "github.com/bernardmuller/shippy/shippy-service-consignment/rpc/consignment"
	"net/http"
	"sync"

	"context"
)

const (
	port = ":50051"
)

type repository interface {
	Create(*c.Consignment) (*c.Consignment, error)
	GetAll() []*c.Consignment
}

type Repository struct {
	mu           sync.RWMutex
	consignments []*c.Consignment
}

func (repo *Repository) Create(consignment *c.Consignment) (*c.Consignment, error) {
	repo.mu.Lock()
	updated := append(repo.consignments, consignment)
	repo.consignments = updated
	repo.mu.Unlock()
	return consignment, nil
}

func (repo *Repository) GetAll() []*c.Consignment {
	return repo.consignments
}

type Service struct {
	repo repository
}

func (s *Service) CreateConsignment(ctx context.Context, req *c.Consignment) (*c.Response, error) {
	consignment, err := s.repo.Create(req)
	if err != nil {
		return nil, err
	}

	return &c.Response{Created: true, Consignment: consignment}, nil
}

func (s *Service) GetConsignments(ctx context.Context, req *c.GetRequest) (*c.Response, error) {
	consignments := s.repo.GetAll()
	return &c.Response{Consignments: consignments}, nil
}

func main() {
	server := &Service{}

	service := c.NewConsignmentServiceServer(server, nil)

	err := http.ListenAndServe(port, service)
	if err != nil {
		panic(err)
	}
}
