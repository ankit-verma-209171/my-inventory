package main

import (
	"database/sql"
	"errors"
	"fmt"
)

type product struct {
	Id       int     `json:"id"`
	Name     string  `json:"name"`
	Quantity int     `json:"quantity"`
	Price    float64 `json:"price"`
}

func (p *product) getProductById(db *sql.DB) error {
	query := fmt.Sprintf("SELECT `id`, `name`, `quantity`, `price` FROM `products` WHERE `id`=%v;", p.Id)
	row := db.QueryRow(query)
	return row.Scan(&p.Id, &p.Name, &p.Quantity, &p.Price)
}

func (p *product) updateProductById(db *sql.DB) error {
	query := fmt.Sprintf("UPDATE `products` SET `name`='%v', `quantity`=%v, `price`=%v WHERE `id`=%v;", p.Name, p.Quantity, p.Price, p.Id)
	_, err := db.Exec(query)
	return err
}

func (p *product) deleteProductById(db *sql.DB) error {
	query := fmt.Sprintf("DELETE FROM `products` WHERE `id`=%v;", p.Id)
	_, err := db.Exec(query)
	return err
}

func (p *product) createProduct(db *sql.DB) error {
	query := fmt.Sprintf("INSERT INTO `products`(name, quantity, price) VALUES('%v', %v, %v);", p.Name, p.Quantity, p.Price)
	result, err := db.Exec(query)

	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("Product not created")
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	p.Id = int(id)
	return nil
}

func getProducts(db *sql.DB) ([]product, error) {
	query := "SELECT `id`, `name`, `quantity`, `price` FROM `products`;"
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}

	products := []product{}
	for rows.Next() {
		var p product
		err := rows.Scan(&p.Id, &p.Name, &p.Quantity, &p.Price)

		if err != nil {
			return nil, err
		}
		products = append(products, p)
	}
	return products, nil
}
