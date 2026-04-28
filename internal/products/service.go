package products

import "context"

type Service interface {
	// Define service methods here
	ListProducts(ctx context.Context) error
}


type svc struct{
	// repository
}

func NewService() Service {
	return &svc{}
}

func (s *svc) ListProducts(ctx context.Context) error {
	// Implementation for listing products
	return nil
}