package main

import (
	"log"
	"net/http"
	"todo/internal/database"
	"todo/internal/handler"
	"todo/internal/repo"
)

func main() {
	db, err := database.InitSQLite3("./app.db")
	if err != nil {
		log.Fatal("failed connecting database", err)
	}

	defer db.Close()

	repo := repo.New(db)
	handler := handler.New(repo)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/todos", handler.GetAll)
	mux.HandleFunc("POST /api/todos", handler.Create)
	mux.HandleFunc("PATCH /api/todos/{id}", handler.Update)
	mux.HandleFunc("DELETE /api/todos/{id}", handler.Delete)
	mux.HandleFunc("DELETE /api/todos", handler.DeleteAll)

	log.Fatal(http.ListenAndServe(":8080", mux))
}
