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
	_, err = h.Model.Insert(id, req.Attributes.Value)
	if err != nil {
		ape.RenderErr(w, problems.InternalError())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
}
