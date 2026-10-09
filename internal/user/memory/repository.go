package memory

import (
	"context"
	"sync"

	"github.com/Max-shiri-90/OrderPulse/internal/user"
)

type Repository struct {
	mu     sync.RWMutex
	users  []user.User
	nextID int64
}

func NewRepository() *Repository {
	return &Repository{
		users:  make([]user.User, 0),
		nextID: 1,
	}
}

func (r *Repository) Create(
	ctx context.Context,
	u user.User,
) (user.User, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	u.ID = r.nextID
	r.nextID++

	r.users = append(r.users, u)

	return u, nil
}

func (r *Repository) GetByID(
	ctx context.Context,
	id int64,
) (user.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, u := range r.users {
		if u.ID == id {
			return u, nil
		}
	}

	return user.User{}, user.ErrNotFound
}

func (r *Repository) GetByEmail(
	ctx context.Context,
	email string,
) (user.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, u := range r.users {
		if u.Email == email {
			return u, nil
		}
	}

	return user.User{}, user.ErrNotFound
}
