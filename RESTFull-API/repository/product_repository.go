package repository

import (
	"be-golang-chapter-24/RESTFull-API/collections"
	"database/sql"
)

type ProductRepository interface {
	GetAll() ([]collections.Product, error)
	GetByID(id int) (*collections.Product, error)
	Create(product *collections.Product) error
	Update(product *collections.Product) error
	Delete(id int) error
}

type productRepository struct {
	db *sql.DB
}

func NewProductRepository(db *sql.DB) ProductRepository {
	return &productRepository{db}
}

func (r *productRepository) GetAll() ([]collections.Product, error) {
	rows, err := r.db.Query("SELECT id, name, price FROM products")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	products := []collections.Product{}
	for rows.Next() {
		var product collections.Product
		if err := rows.Scan(&product.ID, &product.Name, &product.Price); err != nil {
			return nil, err
		}
		products = append(products, product)
	}
	return products, nil
}

func (r *productRepository) GetByID(id int) (*collections.Product, error) {
	var product collections.Product
	err := r.db.QueryRow("SELECT id, name, price FROM products WHERE id=$1", id).Scan(&product.ID, &product.Name, &product.Price)
	if err != nil {
		return nil, err
	}
	return &product, nil
}

func (r *productRepository) Create(product *collections.Product) error {
	_, err := r.db.Exec("INSERT INTO products (name, price) VALUES ($1, $2)", product.Name, product.Price)
	return err
}

func (r *productRepository) Update(product *collections.Product) error {
	_, err := r.db.Exec("UPDATE products SET name=$1, price=$2 WHERE id=$3", product.Name, product.Price, product.ID)
	return err
}

func (r *productRepository) Delete(id int) error {
	_, err := r.db.Exec("DELETE FROM products WHERE id=$1", id)
	return err
}
