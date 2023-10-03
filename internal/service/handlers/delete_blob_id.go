package handlers

import (
	"fmt"
	"net/http"

	requests "BlobApi/internal/service/requests"
)

func (h *BlobHandler) DeleteBlob(w http.ResponseWriter, r *http.Request) {

	id, err := requests.DecodeDeleteBlobRequest(r)
	if err != nil || id < 1 {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	error := h.Model.Delete(id)
	if error != nil {
		http.Error(w, fmt.Sprintf("error deleting blob: %v", error), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}
