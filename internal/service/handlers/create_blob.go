package handlers

import (
	"net/http"
	"strconv"

	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"

	postgres "BlobApi/internal/data/postgres"
	requests "BlobApi/internal/service/requests"
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

	// Converting string ID to int
	id, err := strconv.Atoi(req.Relationships.Owner.Data.ID)
	if err != nil {
		Log(r).WithError(err).Error("Invalid ID format:")
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	// Inserting a blob
	id, err = h.Model.Insert(id, req.Attributes.Value)
	if err != nil {
		ape.RenderErr(w, problems.InternalError())
		return
	}

	// Getting a blob to return the created resource
	blob, err := h.Model.Get(id)
	if err != nil {

		Log(r).WithError(err).Error("error getting blob:")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	// Wrap Blob in AttributeData and Response structures

	response := Response{
		Data: BlobData{
			ID: blob.Index, // Convert int ID to string
			Attributes: BlobAttributes{
				Value: blob.Data,
			},
			Relationships: BlobRelationships{
				Owner: BlobOwner{
					Data: OwnerData{
						ID: *blob.User_id, // Convert int UserID to string
					},
				},
			},
		},
	}

	w.WriteHeader(http.StatusCreated)
	ape.Render(w, &response)
}
