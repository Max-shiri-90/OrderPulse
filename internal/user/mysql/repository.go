package mysql

import (
	"context"
	"database/sql"

	"github.com/Max-shiri-90/OrderPulse/internal/user"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) Create(
	ctx context.Context,
	u user.User,
) (user.User, error) {
	result, err := r.db.ExecContext(
		ctx,
		`
		INSERT INTO users (email, password_hash)
		VALUES (?, ?)
		`,
		u.Email,
		u.PasswordHash,
	)
	if err != nil {
		return user.User{}, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return user.User{}, err
	}

	u.ID = id

	return u, nil
}

func (r *Repository) GetByID(
	ctx context.Context,
	id int64,
) (user.User, error) {
	var u user.User

	err := r.db.QueryRowContext(
		ctx,
		`
		SELECT id, email, password_hash, created_at
		FROM users
		WHERE id = ?
		`,
		id,
	).Scan(
		&u.ID,
		&u.Email,
		&u.PasswordHash,
		&u.CreatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return user.User{}, user.ErrNotFound
		}

		return user.User{}, err
	}

	return u, nil
}

func (r *Repository) GetByEmail(
	ctx context.Context,
	email string,
) (user.User, error) {
	var u user.User

	err := r.db.QueryRowContext(
		ctx,
		`
		SELECT id, email, password_hash, created_at
		FROM users
		WHERE email = ?
		`,
		email,
	).Scan(
		&u.ID,
		&u.Email,
		&u.PasswordHash,
		&u.CreatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return user.User{}, user.ErrNotFound
		}

		return user.User{}, err
	}

	return u, nil
}
