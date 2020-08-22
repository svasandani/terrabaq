package util

import (
	"log"
	"net/http"
)

// CheckHTTPError - helper function to check error and write response
func CheckHTTPError(msg string, err error, w http.ResponseWriter) bool {
	if err == nil {
		return true
	} else {
		log.Println(msg)
		log.Println(err)

		w.WriteHeader(http.StatusInternalServerError)

		w.Write(nil)

		return false
	}
}
