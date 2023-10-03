package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	requests "BlobApi/internal/service/requests"
)

type Response struct {
	Data AttributeData `json:"data"`
}

type AttributeData struct {
	Attributes interface{} `json:"attributes"`
}

func (h *BlobHandler) GetBlobID(w http.ResponseWriter, r *http.Request) {

	id, err := requests.DecodeGetBlobRequest(r)
	if err != nil || id < 1 {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	// Retrieve record by ID
	blob, err := h.Model.Get(id)
	if err != nil {
		http.Error(w, fmt.Sprintf("error getting blob: %v", err), http.StatusInternalServerError)
		return
	}

	if blob == nil {
		http.Error(w, "No blob found", http.StatusNotFound)
		return
	}

	// Wrap Blob in AttributeData and Response structures
	resp := Response{
		Data: AttributeData{
			Attributes: blob,
		},
	}

	// Content-Type header for the response
	w.Header().Set("Content-Type", "application/json")

	// Encoding and sending the response
	err_res := json.NewEncoder(w).Encode(resp)
	if err_res != nil {
		http.Error(w, "Cannot encode response", http.StatusInternalServerError)
	}

}
