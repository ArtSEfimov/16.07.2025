package response

import (
	"encoding/json"
	"net/http"
)

func JsonResponse(w http.ResponseWriter, data interface{}, statusCode int) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if statusCode == http.StatusNoContent {
		return
	}

	encodeErr := json.NewEncoder(w).Encode(data)
	if encodeErr != nil {
		panic(encodeErr)
	}
}
