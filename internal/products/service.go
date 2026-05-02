package products

import (
	"context"
	repo "api_project/internal/adapters/postgresql/sqlc"
)

type Service interface {
	// Define service methods here
	ListProducts(ctx context.Context) ([]repo.Product, error) //тут ми визначаємо метод ListProducts який приймає контекст і повертає список продуктів і помилку, якщо виникає помилка, ми повертаємо її клієнту, якщо ж ні то ми повертаємо список продуктів клієнту
}


type svc struct{
	// repository
	repo repo.Querier
}

func NewService(repo repo.Querier) Service {
	return &svc{
		repo: repo,
	}
}

func (s *svc) ListProducts(ctx context.Context) ([]repo.Product, error) {
	// Implementation for listing products	
	return s.repo.ListProducts(ctx)
	
	// тут ми викликаємо метод ListProducts який знаходиться в нашому репозиторії і передаємо йому контекст, якщо виникає помилка, ми повертаємо її клієнту
}
