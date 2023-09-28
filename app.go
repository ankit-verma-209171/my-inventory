package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

type App struct {
	Router *mux.Router
	DB     *sql.DB
}

func (app *App) Initialize(DbUser, DbPassword, DbName string) (err error) {
	connectionString := fmt.Sprintf("%v:%v@tcp(127.0.0.1:3306)/%v", DbUser, DbPassword, DbName)
	app.DB, err = sql.Open("mysql", connectionString)

	if err != nil {
		return err
	}

	app.Router = mux.NewRouter().StrictSlash(true)
	app.handleRequests()
	return nil
}

func (app *App) Run(address string) {
	err := http.ListenAndServe(address, app.Router)
	if err != nil {
		log.Fatalln(err)
	}
}

func sendResponse(w http.ResponseWriter, statusCode int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	w.Write(response)
}

func sendError(w http.ResponseWriter, statusCode int, err string) {
	errorMessage := map[string]string{"error": err}
	sendResponse(w, statusCode, errorMessage)
}

func (app *App) getProducts(w http.ResponseWriter, r *http.Request) {
	products, err := getProducts(app.DB)
	if err != nil {
		sendError(w, http.StatusInternalServerError, err.Error())
		return
	}
	sendResponse(w, http.StatusOK, products)
}

func (app *App) getProductById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		sendError(w, http.StatusBadRequest, "Invalid product id")
		return
	}

	p := product{Id: id}
	err = p.getProductById(app.DB)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			sendError(w, http.StatusNotFound, "Product not found")
		default:
			sendError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}
	sendResponse(w, http.StatusOK, p)
}

func (app *App) createProduct(w http.ResponseWriter, r *http.Request) {
	var p product
	err := json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		sendError(w, http.StatusBadRequest, "Invalid product details")
		return
	}

	err = p.createProduct(app.DB)
	if err != nil {
		sendError(w, http.StatusBadRequest, err.Error())
		return
	}

	sendResponse(w, http.StatusCreated, p)
}

func (app *App) updateProductById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		sendError(w, http.StatusBadRequest, "Invalid product id")
		return
	}

	var p product
	err = json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		sendError(w, http.StatusBadRequest, "Invalid product details")
		return
	}

	p.Id = id
	err = p.updateProductById(app.DB)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			sendError(w, http.StatusNotFound, "Product not updated")
		default:
			sendError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}
	sendResponse(w, http.StatusOK, p)
}

func (app *App) deleteProductById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		sendError(w, http.StatusBadRequest, "Invalid product id")
		return
	}

	p := product{Id: id}
	err = p.deleteProductById(app.DB)
	if err != nil {
		sendError(w, http.StatusInternalServerError, err.Error())
		return
	}
	sendResponse(w, http.StatusOK, map[string]string{"success": "Deleted product successfully"})
}

func (app *App) handleRequests() {
	app.Router.HandleFunc("/products", app.createProduct).Methods("POST")
	app.Router.HandleFunc("/products", app.getProducts).Methods("GET")
	app.Router.HandleFunc("/products/{id}", app.getProductById).Methods("GET")
	app.Router.HandleFunc("/products/{id}", app.updateProductById).Methods("PUT")
	app.Router.HandleFunc("/products/{id}", app.deleteProductById).Methods("DELETE")
}
