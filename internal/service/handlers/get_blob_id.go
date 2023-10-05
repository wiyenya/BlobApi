package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	requests "BlobApi/internal/service/requests"
)

type Response struct {
	Data BlobData `json:"data"`
}

type BlobData struct {
	ID            string            `json:"id"`
	Attributes    BlobAttributes    `json:"attributes"`
	Relationships BlobRelationships `json:"relationships"`
}

type BlobAttributes struct {
	Value string `json:"value"`
}

type BlobRelationships struct {
	Owner BlobOwner `json:"owner"`
}

type BlobOwner struct {
	Data OwnerData `json:"data"`
}

type OwnerData struct {
	ID string `json:"id"`
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
		Data: BlobData{
			ID: strconv.Itoa(blob.ID), // Convert int ID to string
			Attributes: BlobAttributes{
				Value: blob.Data,
			},
			Relationships: BlobRelationships{
				Owner: BlobOwner{
					Data: OwnerData{
						ID: strconv.Itoa(int(*blob.UserID)), // Convert int UserID to string
					},
				},
			},
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
