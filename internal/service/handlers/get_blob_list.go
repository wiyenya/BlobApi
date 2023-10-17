package handlers

import (
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

	// List for storing converted data Blob
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

	// Collect the response
	resp := ListResponse{
		Data: responseData,
	}

	// Content-Type header for the response
	w.Header().Set("Content-Type", "application/json")
	ape.Render(w, &resp)
}
