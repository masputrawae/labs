package handler

import (
	"database/sql"
	"errors"
	"log"
	"net/http"
	"strconv"
	"todo/internal/domain/model"
	"todo/internal/usecase"

	"github.com/go-playground/validator/v10"
)

type Todo interface {
	Create(w http.ResponseWriter, r *http.Request)
	Update(w http.ResponseWriter, r *http.Request)
	Delete(w http.ResponseWriter, r *http.Request)
	DeleteAll(w http.ResponseWriter, r *http.Request)
	GetAll(w http.ResponseWriter, r *http.Request)
	GetStatuses(w http.ResponseWriter, r *http.Request)
	GetPriorities(w http.ResponseWriter, r *http.Request)
}

type todo struct {
	usecaseTodo usecase.Todo
	validate    *validator.Validate
}

func NewTodo(
	usecaseTodo usecase.Todo,
	validate *validator.Validate,
) Todo {
	return &todo{usecaseTodo: usecaseTodo, validate: validate}
}

func (t *todo) Create(w http.ResponseWriter, r *http.Request) {
	userID, _ := r.Context().Value("userID").(string)

	req := model.TodoCreateRequest{
		UserID: userID,
	}

	if !Bind(w, r, t.validate, &req) {
		return
	}

	log.Println(req)

	res, err := t.usecaseTodo.Create(r.Context(), req)
	if err != nil {
		log.Println("error create todo:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	JSON(w, http.StatusCreated, H{"data": res})
}

func (t *todo) Update(w http.ResponseWriter, r *http.Request) {
	userID, _ := r.Context().Value("userID").(string)
	todoID, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	req := model.TodoUpdateRequest{
		UserID: userID,
		ID:     todoID,
	}

	if !Bind(w, r, t.validate, &req) {
		return
	}

	res, err := t.usecaseTodo.Update(r.Context(), req)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		log.Println("error update todo:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	JSON(w, http.StatusOK, H{"data": res})
}

func (t *todo) Delete(w http.ResponseWriter, r *http.Request) {
	userID, _ := r.Context().Value("userID").(string)
	todoID, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := t.usecaseTodo.Delete(r.Context(), userID, todoID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		log.Println("error delete todo:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (t *todo) DeleteAll(w http.ResponseWriter, r *http.Request) {
	userID, _ := r.Context().Value("userID").(string)
	if err := t.usecaseTodo.DeleteAll(r.Context(), userID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		log.Println("error delete all todos:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (t *todo) GetAll(w http.ResponseWriter, r *http.Request) {
	userID, _ := r.Context().Value("userID").(string)
	res, err := t.usecaseTodo.GetAll(r.Context(), userID)
	if err != nil {
		log.Println("error get all todos:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	JSON(w, http.StatusOK, H{"data": res})
}

func (t *todo) GetStatuses(w http.ResponseWriter, r *http.Request) {
	res, err := t.usecaseTodo.GetStatuses(r.Context())
	if err != nil {
		log.Println("error get all statuses:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	JSON(w, http.StatusOK, H{"data": res})
}

func (t *todo) GetPriorities(w http.ResponseWriter, r *http.Request) {
	res, err := t.usecaseTodo.GetPriorities(r.Context())
	if err != nil {
		log.Println("error get all priorities:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	JSON(w, http.StatusOK, H{"data": res})
}
