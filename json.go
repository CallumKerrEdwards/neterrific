package neterrific

import (
	"encoding/json"
	"net/http"
)

func SendHTTPJSONError(w http.ResponseWriter, status int, err error) {
	SendJSON(w, status, Payload{
		"status": status,
		"error":  err.Error(),
	})
}

func SendJSON(w http.ResponseWriter, status int, p any) {
	w.WriteHeader(status)

	err := json.NewEncoder(w).Encode(p)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}
