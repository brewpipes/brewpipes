package handler

import (
	"net/http"
	"strconv"
)

func parseInt64Param(value string) (int64, error) {
	return strconv.ParseInt(value, 10, 64)
}

func methodNotAllowed(w http.ResponseWriter) {
	http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
}
