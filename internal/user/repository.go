package user

import "context"

type Repository interface {
	Create(ctx context.Context, u User) (User, error)
	GetByID(ctx context.Context, id int64) (User, error)
	GetByEmail(ctx context.Context, email string) (User, error)
}
