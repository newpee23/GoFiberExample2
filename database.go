package main

import (
	"database/sql"
	"fmt"
)

func createProduct(product *Product) error {
	var id int
	err := dbConnect.QueryRow(`INSERT INTO products(name, price) VALUES($1, $2) RETURNING id;`, product.Name, product.Price).Scan(&id)

	if err != nil {
		return err
	}

	fmt.Printf("New product ID is %d\n", id)
	return nil
}

func getProductById(id int) (Product, error) {
	var p Product
	row := dbConnect.QueryRow(`SELECT id, name, price FROM products WHERE id = $1;`, id)
	err := row.Scan(&p.ID, &p.Name, &p.Price)
	if err != nil {
		return Product{}, err
	}
	return p, nil
}

func getProducts() ([]Product, error) {
	rows, err := dbConnect.Query("SELECT id, name, price FROM products")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []Product
	for rows.Next() {
		var p Product
		err := rows.Scan(&p.ID, &p.Name, &p.Price)
		if err != nil {
			return nil, err
		}
		products = append(products, p)
	}

	// Check for errors from iterating over rows
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return products, nil
}

func updateProduct(id int, product *Product) (Product, error) {
	var p Product
	row := dbConnect.QueryRow(`UPDATE products SET name = $1, price = $2 WHERE id = $3 RETURNING id , name , price;`, product.Name, product.Price, id)

	err := row.Scan(&p.ID, &p.Name, &p.Price)
	if err != nil {
		return Product{}, err
	}
	return p, err
}

func deleteProduct(id int) error {
	_, err := dbConnect.Exec(`DELETE FROM products WHERE id = $1;`, id)
	return err
}

func getProductsAndSuppliers() ([]ProductWithSupplier, error) {
	// SQL JOIN query
	query := `
		SELECT
			p.id AS product_id,
			p.name AS product_name,
			p.price,
			s.name AS supplier_name
		FROM
			products p
		LEFT JOIN suppliers s
			ON p.supplier_id = s.id;`

	rows, err := dbConnect.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []ProductWithSupplier
	for rows.Next() {
		var p ProductWithSupplier
		var supplierName sql.NullString
		err := rows.Scan(&p.ProductID, &p.ProductName, &p.Price, &supplierName)
		if err != nil {
			return nil, err
		}
		p.SupplierName = formatSupplierName(supplierName)
		products = append(products, p)
	}

	// Check if there are no rows returned
	if len(products) == 0 {
		return nil, fmt.Errorf("no products with suppliers found")
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return products, nil
}

// Function to format SupplierName
func formatSupplierName(s sql.NullString) string {
	if s.Valid {
		return s.String
	}
	return "-"
}

func addProductAndSupplier(data *ProductWithSupplier) error {
	// Start a transaction
	tx, err := dbConnect.Begin()
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			// Rollback the transaction in case of an error
			tx.Rollback()
			return
		}
	}()

	// Insert into the supplier table and retrieve the inserted ID
	var supplierID int64
	err = tx.QueryRow("INSERT INTO suppliers (name) VALUES ($1) RETURNING id", data.SupplierName).Scan(&supplierID)
	if err != nil {
		return err
	}

	// Insert into the product table
	_, err = tx.Exec("INSERT INTO products (name, price, supplier_id) VALUES ($1, $2, $3)", data.ProductName, data.Price, supplierID)
	if err != nil {
		return err
	}

	// Commit the transaction
	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}
