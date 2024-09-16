package utils

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"
)

type HTTPError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func RespondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Failed to marshal JSON response: %v", err)
		RespondWithError(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	if len(response) > 0 {
		w.Write(response)
	} else {
		w.Write([]byte("{}"))
	}
}

func RespondWithError(w http.ResponseWriter, code int, message string) {
	RespondWithJSON(w, code, HTTPError{Code: code, Message: message})
}

func ParseQueryParamInt(r *http.Request, key string, defaultValue int) (int, error) {
	values := r.URL.Query()
	if valStr := values.Get(key); valStr != "" {
		val, err := strconv.Atoi(valStr)
		if err != nil {
			return 0, err
		}
		return val, nil
	}
	return defaultValue, nil
}

var validate = validator.New()

func ValidateStruct(s interface{}) error {
	return validate.Struct(s)
}

func GetPaginationParams(r *http.Request) (limit, offset int, err error) {
	limit, err = ParseQueryParamInt(r, "limit", 5)
	if err != nil {
		return limit, offset, err
	}
	offset, err = ParseQueryParamInt(r, "offset", 0)
	return limit, offset, err
}
