package repository

import (
	"context"

	"nvsu-api/model"
)

type Repository interface {
	GetUser(ctx context.Context, email string) (model.User, error)
	GetProduct(ctx context.Context, id string) (model.Product, error)
	GetOrder(ctx context.Context, id string) (model.Order, error)
	GetProducts(ctx context.Context) ([]*model.Product, error)
	CreateUser(ctx context.Context, in model.User) (model.User, error)
	CreateProduct(ctx context.Context, in model.Product) (model.Product, error)
	CreateOrder(ctx context.Context, in model.Order) (model.Order, error)
	UpdateUser(ctx context.Context, in model.User) (model.User, error)
	UpdateOrder(ctx context.Context, in model.Order) (model.Order, error)
	DeleteUser(ctx context.Context, email string) error
}
