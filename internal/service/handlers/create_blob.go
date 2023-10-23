package handlers

import (
	"fmt"
	"net/http"

	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"

	postgres "BlobApi/internal/data/postgres"
	requests "BlobApi/internal/service/requests"
	resources "BlobApi/resources"
)

type BlobHandler struct {
	Model *postgres.BlobModel
}

func NewBlobHandler(m *postgres.BlobModel) *BlobHandler {
	return &BlobHandler{Model: m}
}

func (h *BlobHandler) CreateBlob(w http.ResponseWriter, r *http.Request) {
	// Decoding the request body

	req, err := requests.DecodeCreateBlobRequest(r)
	if err != nil {
		Log(r).WithError(err).Error("BadRequest")
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	id := req.Relationships.UserId

	// Inserting a blob
	blobId, err := h.Model.Insert(id, req.Attributes.Value)
	if err != nil {
		ape.RenderErr(w, problems.InternalError())
		return
	}

	// Getting a blob to return the created resource
	blob, err := h.Model.Get(blobId)
	if err != nil {

		Log(r).WithError(err).Error("error getting blob:")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	// Wrap Blob in AttributeData and Response structures

	response := resources.BlobResponse{
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

	w.WriteHeader(http.StatusCreated)
	ape.Render(w, &response)
}
