package mysql

import (
	"context"
	"database/sql"

	"github.com/Max-shiri-90/OrderPulse/internal/product"
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
	p product.Product,
) (product.Product, error) {
	result, err := r.db.ExecContext(
		ctx,
		`
		INSERT INTO products (name, price, quantity)
		VALUES (?, ?, ?)
		`,
		p.Name,
		p.Price,
		p.Quantity,
	)
	if err != nil {
		return product.Product{}, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return product.Product{}, err
	}

	p.ID = id

	return p, nil
}

func (r *Repository) GetByID(
	ctx context.Context,
	id int64,
) (product.Product, error) {
	var p product.Product

	err := r.db.QueryRowContext(
		ctx,
		`
		SELECT id, name, price, quantity
		FROM products
		WHERE id = ?
		`,
		id,
	).Scan(
		&p.ID,
		&p.Name,
		&p.Price,
		&p.Quantity,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return product.Product{}, product.ErrNotFound
		}

		return product.Product{}, err
	}

	return p, nil
}

func (r *Repository) List(
	ctx context.Context,
) ([]product.Product, error) {
	rows, err := r.db.QueryContext(
		ctx,
		`
		SELECT id, name, price, quantity
		FROM products
		ORDER BY id
		`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	products := make([]product.Product, 0)

	for rows.Next() {
		var p product.Product

		if err := rows.Scan(
			&p.ID,
			&p.Name,
			&p.Price,
			&p.Quantity,
		); err != nil {
			return nil, err
		}

		products = append(products, p)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return products, nil
}
