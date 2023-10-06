package handlers

import (
	"net/http"
	"strconv"

	requests "BlobApi/internal/service/requests"

	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
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

		Log(r).WithError(err).Error("Invalid ID")
		ape.RenderErr(w, problems.BadRequest(err)...)

		return
	}

	// Retrieve record by ID
	blob, err := h.Model.Get(id)
	if err != nil {

		Log(r).WithError(err).Error("error getting blob:")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	if blob == nil {

		Log(r).WithError(err).Error("No blob found")
		ape.RenderErr(w, problems.NotFound())
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

	ape.Render(w, &resp)
}
