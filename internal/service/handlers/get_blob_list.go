package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func (h *BlobHandler) GetBlobList(w http.ResponseWriter, r *http.Request) {

	blobs, err := h.Model.GetBlobList()
	if err != nil {
		http.Error(w, fmt.Sprintf("error getting blobs: %v", err), http.StatusInternalServerError)
		return
	}

	if blobs == nil {
		http.Error(w, "No blob found", http.StatusNotFound)
		return
	}

	// Wrap Blob in AttributeData and Response structures
	resp := Response{
		Data: AttributeData{
			Attributes: blobs,
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
