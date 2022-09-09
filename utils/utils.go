package utils

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"swoyo/models"
)

const (
	ENCODED_URLS = "encoded_urls"
)

func handlerError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func HandleHttpErrors(w http.ResponseWriter, customErrMsg string, statusCode int, err error) (isError bool) {
	if err != nil {
		errorMsg := err.Error()
		if customErrMsg != "" {
			errorMsg = customErrMsg
		}
		response := models.ErrorResponse{
			Message: errorMsg,
			Status:  statusCode,
		}

		message, err := json.Marshal(response)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return true
		}

		w.WriteHeader(response.Status)
		_, err = w.Write(message)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return true
		}
		return true
	}
	return
}

func ValidateUrl(urlString string) error {
	u, err := url.Parse(urlString)
	if err != nil || u.Scheme == "" || u.Host == "" {
		return err
	}

	_, err = url.ParseRequestURI(urlString)
	if err != nil {
		return err
	}
	return err
}

func SetCorsHeaders(f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Headers:", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET,POST,PUT,DELETE,PATCH,OPTIONS")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		f(w, r)
	}
}
