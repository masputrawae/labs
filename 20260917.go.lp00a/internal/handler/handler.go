package handler

import (
	"database/sql"
	"errors"
	"log"
	"net/http"
	"todo/internal/model"
	"todo/internal/repo"
)

type Handler interface {
	Create(w http.ResponseWriter, r *http.Request)
	Update(w http.ResponseWriter, r *http.Request)
	GetAll(w http.ResponseWriter, r *http.Request)
	Delete(w http.ResponseWriter, r *http.Request)
	DeleteAll(w http.ResponseWriter, r *http.Request)
}

type handler struct {
	repo repo.Repo
}

func New(repo repo.Repo) Handler {
	return &handler{repo: repo}
}

func (h *handler) Create(w http.ResponseWriter, r *http.Request) {
	req, ok := Bind[model.TodoCreateRequest](w, r)
	if !ok {
		return
	}

	if req.Task == "" {
		JSON(w, http.StatusBadRequest, H{"error": "task can not empty"})
		return
	}

	result, err := h.repo.Create(r.Context(), req.Task)
	if err != nil {
		log.Println("failed add todo to database", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	JSON(w, http.StatusOK, H{"data": result})
}

func (h *handler) Update(w http.ResponseWriter, r *http.Request) {
	req, ok := Bind[model.TodoUpdateRequest](w, r)
	if !ok {
		return
	}

	if req.IsDone == nil && req.Task == nil {
		JSON(w, http.StatusBadRequest, H{"error": "at least one field has been changed"})
		return
	}

	todoID, ok := GetParamInt64(w, r)
	if !ok {
		JSON(w, http.StatusBadRequest, H{"error": "id must be an integer"})
		return
	}

	result, err := h.repo.Update(r.Context(), todoID, req)
	if err != nil {
		log.Println("failed update (database):", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	JSON(w, http.StatusOK, H{"data": result})
}

func (h *handler) GetAll(w http.ResponseWriter, r *http.Request) {
	result, err := h.repo.GetAll(r.Context())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	JSON(w, http.StatusOK, H{"data": result})
}

func (h *handler) Delete(w http.ResponseWriter, r *http.Request) {
	todoID, ok := GetParamInt64(w, r)
	if !ok {
		return
	}

	err := h.repo.Delete(r.Context(), todoID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *handler) DeleteAll(w http.ResponseWriter, r *http.Request) {
	err := h.repo.DeleteAll(r.Context())
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
