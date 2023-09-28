package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

var testApp App

func TestMain(m *testing.M) {
	err := testApp.Initialize(DbUser, DbPassword, "test")
	if err != nil {
		panic(err)
	}

	createTable()
	m.Run()
}

func createTable() {
	query := `CREATE TABLE IF NOT EXISTS products (
			  id 		INT			 NOT NULL AUTO_INCREMENT,
			  name		VARCHAR(255) NOT NULL UNIQUE,
			  quantity	INT,
			  price 	FLOAT(10, 2),
			  PRIMARY KEY (id),
			  CHECK (name!='' AND quantity>=0 AND price>=0)
			  );`

	_, err := testApp.DB.Exec(query)
	if err != nil {
		panic(err)
	}
}

func clearTable() {
	query := "DELETE FROM `products`;"
	_, err := testApp.DB.Exec(query)
	if err != nil {
		panic(err)
	}

	query = "ALTER TABLE `products` AUTO_INCREMENT = 1;"
	_, err = testApp.DB.Exec(query)
	if err != nil {
		panic(err)
	}
}

func addProduct(name string, quantity int, price float64) {
	query := fmt.Sprintf("INSERT INTO `products`(name, quantity, price) VALUES('%v', %v, %v);", name, quantity, price)
	_, err := testApp.DB.Exec(query)
	if err != nil {
		panic(err)
	}
}

func checkStatusCode(t *testing.T, expectedCode, actualCode int) {
	if expectedCode != actualCode {
		t.Errorf("Expected status %v, Received %v \n", expectedCode, actualCode)
	}
}

func TestGetProduct(t *testing.T) {
	clearTable()
	addProduct("table", 80, 1000)
	request, err := http.NewRequest("GET", "/products/1", nil)
	if err != nil {
		panic(err)
	}
	response := sendRequest(request)
	checkStatusCode(t, http.StatusOK, response.Code)
}

func sendRequest(request *http.Request) *httptest.ResponseRecorder {
	recorder := httptest.NewRecorder()
	testApp.Router.ServeHTTP(recorder, request)
	return recorder
}

func TestCreateProduct(t *testing.T) {
	clearTable()
	body := []byte(`{"name":"pencil", "quantity":12, "price":1234}`)
	request, err := http.NewRequest("POST", "/products", bytes.NewBuffer(body))
	if err != nil {
		panic(err)
	}
	response := sendRequest(request)
	checkStatusCode(t, http.StatusCreated, response.Code)
}

func TestDeleteProduct(t *testing.T) {
	clearTable()
	addProduct("table", 10, 100)

	request, err := http.NewRequest("DELETE", "/products/1", nil)
	if err != nil {
		panic(err)
	}
	response := sendRequest(request)
	checkStatusCode(t, http.StatusOK, response.Code)

	request, err = http.NewRequest("GET", "/products/1", nil)
	if err != nil {
		panic(err)
	}
	response = sendRequest(request)
	checkStatusCode(t, http.StatusNotFound, response.Code)
}

func TestUpdateProduct(t *testing.T) {
	clearTable()
	addProduct("table", 10, 100)

	body := []byte(`{"name":"pencil", "quantity":12, "price":100}`)
	request, err := http.NewRequest("PUT", "/products/1", bytes.NewBuffer(body))
	if err != nil {
		panic(err)
	}
	response := sendRequest(request)
	checkStatusCode(t, http.StatusOK, response.Code)

	var p map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &p)

	if p["name"] != "pencil" {
		t.Errorf("Expected name %v, Received %v \n", "pencil", p["name"])
	}
}
