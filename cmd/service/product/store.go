package product

import (
	"database/sql"
	"fmt"
	"strings"

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

func (s *Store) UpdateProduct(productUpdate types.ProductAndInventoryUpdate) error {
	// Set up transactions
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}

	// Prepare the base SQL update query for products table

	updateQuery := `UPDATE products SET `
	var setClauses []string
	var setValues []interface{}

	if productUpdate.Name != nil {
		setClauses = append(setClauses, "name = ?")
		setValues = append(setValues, *productUpdate.Name)
	}

	if productUpdate.Description != nil {
		setClauses = append(setClauses, "description = ?")
		setValues = append(setValues, *productUpdate.Description)
	}

	if productUpdate.Image != nil {
		setClauses = append(setClauses, "image_url = ?")
		setValues = append(setValues, *productUpdate.Image)
	}

	if productUpdate.IsActive != nil {
		setClauses = append(setClauses, "isActive = ?")
		setValues = append(setValues, *productUpdate.IsActive)
	}

	if productUpdate.Price != nil {
		setClauses = append(setClauses, "price = ?")
		setValues = append(setValues, *productUpdate.Price)
	}

	if len(setClauses) > 0 {
		updateQuery += strings.Join(setClauses, ", ")
		updateQuery += " WHERE id = ?"
		setValues = append(setValues, productUpdate.ProductID)

		// Execute the update query
		_, err := tx.Exec(updateQuery, setValues...)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	// Build query for the inventory table
	inventoryUpdateQuery := "UPDATE product_inventory SET "
	var inventorySetClauses []string
	var inventoryValues []interface{}

	if productUpdate.Quantity != nil {
		inventorySetClauses = append(inventorySetClauses, "quantity_available = ?")
		inventoryValues = append(inventoryValues, *productUpdate.Quantity)
	}

	if productUpdate.Stock != nil {
		inventorySetClauses = append(inventorySetClauses, "stock = ?")
		inventoryValues = append(inventoryValues, *productUpdate.Stock)
	}
	if len(inventorySetClauses) > 0 {
		inventoryUpdateQuery += strings.Join(inventorySetClauses, ", ")
		inventoryUpdateQuery += " WHERE product_id = ?"
		inventoryValues = append(inventoryValues, productUpdate.ProductID)

		// Execute the inventory update query
		_, err := tx.Exec(inventoryUpdateQuery, inventoryValues...)
		if err != nil {
			tx.Rollback()
			return err
		}
	}
	tx.Commit()
	return nil
}
