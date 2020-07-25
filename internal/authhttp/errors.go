package authhttp

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"

	"github.com/Tsapen/authorization/internal/auth"
)

func processError(w http.ResponseWriter, url *url.URL, method, message string, err error) {
	http.Error(w, message, errStatus(err))
	logError(url, method, message, err)
}

func logError(url *url.URL, method, message string, err error) {
	log.Printf("%s %s: %s: %s\n", method, url, message, err)
}

func errStatus(err error) int {
	switch err.(type) {
	case *json.SyntaxError:
		return http.StatusBadRequest
	case auth.BadParametersError,
		*json.UnmarshalTypeError:
		return http.StatusUnprocessableEntity
	}

	switch err {
	case auth.ErrBadParameters:
		return http.StatusUnprocessableEntity

	case auth.ErrNotFound:
		return http.StatusNotFound

	default:

	}

	return http.StatusInternalServerError
}
