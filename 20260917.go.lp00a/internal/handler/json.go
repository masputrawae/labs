package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
)

type H map[string]any

func Bind[T any](w http.ResponseWriter, r *http.Request) (T, bool) {
	var obj T
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&obj); err != nil {
		log.Println("failed decode json:", err)
		w.WriteHeader(http.StatusBadRequest)
		return obj, false
	}
	return obj, true
}

func JSON(w http.ResponseWriter, status int, data any) {
	body, err := json.Marshal(data)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if _, err := w.Write(body); err != nil {
		log.Printf("failed to write response: %v", err)
	}
}

func GetParamInt64(w http.ResponseWriter, r *http.Request) (int64, bool) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		log.Println("failed parse id str to int64:", err)
		w.WriteHeader(http.StatusBadRequest)
		return -1, false
	}

	return id, true
}
