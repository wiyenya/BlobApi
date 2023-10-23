package handlers

import (
	"fmt"
	"net/http"

	requests "BlobApi/internal/service/requests"

	resources "BlobApi/resources"

	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
)

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
	resp := resources.BlobResponse{
		Data: resources.Blob{
			Key: resources.Key{
				ID:           fmt.Sprint(blob.Index),
				ResourceType: "Blob",
			},
			Attributes: resources.BlobAttributes{
				Obj: blob.Data,
			},
			Relationships: &resources.BlobRelationships{
				UserId: *blob.User_id,
			},
		},
	}

	ape.Render(w, &resp)
}
