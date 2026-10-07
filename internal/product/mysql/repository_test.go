package mysql

import (
	"context"
	"database/sql"
	"os"
	"testing"

	_ "github.com/go-sql-driver/mysql"

	"github.com/Max-shiri-90/OrderPulse/internal/product"
)

func testDatabase(t *testing.T) *sql.DB {
	t.Helper()

	dsn := "orderpulse:orderpulsepassword@tcp(127.0.0.1:3307)/orderpulse?parseTime=true"

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		t.Fatalf("failed to open database: %v", err)
	}

	if err := db.Ping(); err != nil {
		db.Close()
		t.Fatalf("failed to connect to database: %v", err)
	}

	t.Cleanup(func() {
		db.Close()
	})

	return db
}

func TestMain(m *testing.M) {
	code := m.Run()
	os.Exit(code)
}

func TestCreate(t *testing.T) {
	db := testDatabase(t)

	repository := NewRepository(db)

	_, err := db.Exec("DELETE FROM products")
	if err != nil {
		t.Fatalf("failed to clean products table: %v", err)
	}

	p := product.Product{
		Name:     "Test Keyboard",
		Price:    5000000,
		Quantity: 10,
	}

	created, err := repository.Create(context.Background(), p)
	if err != nil {
		t.Fatalf("failed to create product: %v", err)
	}

	if created.ID == 0 {
		t.Fatal("expected created product to have an ID")
	}

	if created.Name != p.Name {
		t.Fatalf("expected name %q, got %q", p.Name, created.Name)
	}

	if created.Price != p.Price {
		t.Fatalf("expected price %d, got %d", p.Price, created.Price)
	}

	if created.Quantity != p.Quantity {
		t.Fatalf("expected quantity %d, got %d", p.Quantity, created.Quantity)
	}
}

func TestGetByID(t *testing.T) {
	db := testDatabase(t)

	repository := NewRepository(db)

	_, err := db.Exec("DELETE FROM products")
	if err != nil {
		t.Fatalf("failed to clean products table: %v", err)
	}

	created, err := repository.Create(
		context.Background(),
		product.Product{
			Name:     "Test Mouse",
			Price:    2500000,
			Quantity: 20,
		},
	)
	if err != nil {
		t.Fatalf("failed to create product: %v", err)
	}

	found, err := repository.GetByID(
		context.Background(),
		created.ID,
	)
	if err != nil {
		t.Fatalf("failed to get product: %v", err)
	}

	if found.ID != created.ID {
		t.Fatalf("expected ID %d, got %d", created.ID, found.ID)
	}

	if found.Name != "Test Mouse" {
		t.Fatalf("expected name %q, got %q", "Test Mouse", found.Name)
	}
}

func TestGetByIDNotFound(t *testing.T) {
	db := testDatabase(t)

	repository := NewRepository(db)

	_, err := db.Exec("DELETE FROM products")
	if err != nil {
		t.Fatalf("failed to clean products table: %v", err)
	}

	_, err = repository.GetByID(
		context.Background(),
		999999,
	)

	if err != product.ErrNotFound {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

func TestList(t *testing.T) {
	db := testDatabase(t)

	repository := NewRepository(db)

	_, err := db.Exec("DELETE FROM products")
	if err != nil {
		t.Fatalf("failed to clean products table: %v", err)
	}

	_, err = repository.Create(
		context.Background(),
		product.Product{
			Name:     "Keyboard",
			Price:    5000000,
			Quantity: 10,
		},
	)
	if err != nil {
		t.Fatalf("failed to create first product: %v", err)
	}

	_, err = repository.Create(
		context.Background(),
		product.Product{
			Name:     "Mouse",
			Price:    2500000,
			Quantity: 20,
		},
	)
	if err != nil {
		t.Fatalf("failed to create second product: %v", err)
	}

	products, err := repository.List(context.Background())
	if err != nil {
		t.Fatalf("failed to list products: %v", err)
	}

	if len(products) != 2 {
		t.Fatalf("expected 2 products, got %d", len(products))
	}

	if products[0].Name != "Keyboard" {
		t.Fatalf("expected first product to be Keyboard, got %q", products[0].Name)
	}

	if products[1].Name != "Mouse" {
		t.Fatalf("expected second product to be Mouse, got %q", products[1].Name)
	}
}
