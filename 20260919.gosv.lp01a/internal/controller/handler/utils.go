package handler

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/go-playground/validator/v10"
)

type H map[string]any

func JSON(w http.ResponseWriter, status int, obj any) {
	raw, err := json.Marshal(obj)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("error marshal json", err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write([]byte(raw))
}

func Bind(w http.ResponseWriter, r *http.Request, v *validator.Validate, obj any) bool {
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&obj); err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusUnprocessableEntity)
		return false
	}

	err := v.StructCtx(r.Context(), obj)
	if err == nil {
		return true
	}

	errs, ok := errors.AsType[validator.ValidationErrors](err)
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		return false
	}

	var results = make(H, len(errs))
	for i := range errs {
		f := errs[i].Field()
		switch errs[i].Tag() {
		case "required":
			results[f] = "REQUIRED"
		case "min":
			results[f] = "TOO_SHORT"
		case "max":
			results[f] = "TOO_LONG"
		case "email":
			results[f] = "INVALID_EMAIL"
		case "eqfield":
			results[f] = "EQ_FIELD"
		case "startswith":
			results[f] = "STARTS_WITH"
		case "oneof":
			results[f] = "ONE_OF"
		}
	}

	JSON(w, http.StatusBadRequest, H{"error": results})
	return false
}
