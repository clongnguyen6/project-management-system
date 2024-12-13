package response

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/go-playground/validator"
)

// Response represents a standard error response for Swagger documentation
type Response struct {
	Status	string `json:"status"`
	Error	string `json:"error"`
}

const(
	StatusOK = "OK"
	StatusError = "Error"
)

// func WriteJson(w http.ResponseWriter, status int, data interface{}) error {
// 	w.Header().Set("Content-Type", "application/json")
// 	w.WriteHeader(status)

// 	return json.NewEncoder(w).Encode(data)
// }

func WriteJson(rw http.ResponseWriter, status int, data interface{}) error {
	js, err := json.Marshal(data)
	if err != nil {
		return err
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(status)
	_, err = rw.Write(js)
	if err != nil {
		return err
	}
	return nil
}

func GeneralError(err error) Response {
	return Response{
		Status: StatusError,
		Error: err.Error(),
	}
}

func ValidationError(errs validator.ValidationErrors) Response {
	var errMaps []string

	for _, err := range errs {
		switch err.ActualTag() {
		case "required":
			errMaps = append(errMaps, fmt.Sprintf("field %s is required field", err.Field()))
		default:
			errMaps = append(errMaps, fmt.Sprintf("field %s is invalid", err.Field()))
		}
	}

	return Response{
		Status: StatusError,
		Error: strings.Join(errMaps, ", "),
	}
}

func HandleCacheControl(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		headers := rw.Header()
		headers.Set("Cache-Control", "no-cache, no-store, max-age=0, must-revalidate")
		headers.Set("Pragma", "no-cache")
		headers.Set("Expires", "0")
		next.ServeHTTP(rw, req)
	})
}
