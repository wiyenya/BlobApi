package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
)

type ListResponse struct {
	Data []BlobData `json:"data"`
}

func (h *BlobHandler) GetBlobList(w http.ResponseWriter, r *http.Request) {

	blobs, err := h.Model.GetBlobList()
	if err != nil {

		Log(r).WithError(err).Error("error getting blob:")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	if blobs == nil {

		Log(r).WithError(err).Error("No blob found")
		ape.RenderErr(w, problems.NotFound())
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
		Log(r).WithError(err).Error("Cannot encode response")
		ape.RenderErr(w, problems.InternalError())
	}

}
