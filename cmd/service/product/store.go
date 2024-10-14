package product

import (
	"database/sql"
	"fmt"

	"github.com/Nasa28/CommerceCore/types"
)

type Store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{db: db}
}

func scanProductsIntoRows(rows *sql.Rows) (*types.ProductAndInventory, error) {
	product := new(types.ProductAndInventory)
	err := rows.Scan(
		&product.ID,
		&product.Name,
		&product.Description,
		&product.Price,
		&product.ImageURL,
		&product.QuantityAvailable,
		&product.Stock,
		&product.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return product, nil
}
func (s *Store) GetProductByID(id int) (*types.ProductAndInventory, error) {
	query := `
		SELECT p.id, p.name, p.description, p.price, p.image_url, 
		       pi.quantity_available, pi.stock,p.createdAt
		FROM products p
		JOIN product_inventory pi ON p.id = pi.product_id
		WHERE p.id = ?`

	row, err := s.db.Query(query, id)
	if err != nil {
		return nil, err
	}
	productWithInventory := new(types.ProductAndInventory)

	for row.Next() {
		productWithInventory, err = scanProductsIntoRows(row)
		if err != nil {
			return nil, err
		}
	}
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("product not found")
		}
		return nil, err
	}

	return productWithInventory, nil
}

func (s *Store) CreateProduct(product types.CreateProductPayload) error {

	// Start a new transaction
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	result, err := tx.Exec("INSERT INTO products (name,description,price,image_url) VALUES(?,?,?,?)", product.Name, product.Description, product.Price, product.Image)

	if err != nil {
		tx.Rollback()
		return err
	}
	// Get the last inserted product ID
	lastInsertID, err := result.LastInsertId()
	if err != nil {
		tx.Rollback()
		return err
	}

	_, err = tx.Exec("INSERT INTO product_inventory (product_id,quantity_available,stock) VALUES (?,?,?)", lastInsertID, product.Quantity, product.Stock)

	if err != nil {
		tx.Rollback()

		return err
	}
	tx.Commit()

	return nil

}
