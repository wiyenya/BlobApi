package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

type ListResponse struct {
	Data []BlobData `json:"data"`
}

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

	// Список для хранения преобразованных данных Blob
	var responseData []BlobData
	for _, blob := range blobs {
		if blob == nil {
			fmt.Println("Warning: encountered nil blob")
			continue
		}
		if blob.UserID == nil {
			fmt.Println("Warning: UserID is nil for blob ID:", blob.ID)
			continue
		}
		responseData = append(responseData, BlobData{
			ID: strconv.Itoa(blob.ID),
			Attributes: BlobAttributes{
				Value: blob.Data,
			},
			Relationships: BlobRelationships{
				Owner: BlobOwner{
					Data: OwnerData{
						ID: strconv.Itoa(int(*blob.UserID)),
					},
				},
			},
		})
	}

	// Собираем и отправляем ответ
	resp := ListResponse{
		Data: responseData,
	}

	// Content-Type header for the response
	w.Header().Set("Content-Type", "application/json")

	// Encoding and sending the response
	err_res := json.NewEncoder(w).Encode(resp)
	if err_res != nil {
		http.Error(w, "Cannot encode response", http.StatusInternalServerError)
	}

}
