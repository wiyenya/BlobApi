package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"

	horizon "BlobApi/internal/data/horizon"
	//postgres "BlobApi/internal/data/postgres"
	requests "BlobApi/internal/service/requests"
	resources "BlobApi/resources"
)

type BlobHandlerHorizon struct {
	Model *horizon.HorizonModel
}

type BlobHandler struct {
	Model *horizon.HorizonModel
}

func NewBlobHandlerHorizon(m *horizon.HorizonModel) *BlobHandlerHorizon {
	return &BlobHandlerHorizon{Model: m}
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

	//data to JSON (map to bytes)
	jsonData, errMarshal := json.Marshal(req.Attributes.Value)
	if errMarshal != nil {
		return
	}

	// Inserting a blob
	blobId, err := h.Model.Insert(id, jsonData)
	if err != nil {
		ape.RenderErr(w, problems.InternalError())
		return
	}

	//Getting a blob to return the created resource
	blob, err := h.Model.Get(blobId)
	if err != nil {

		Log(r).WithError(err).Error("error getting blob:")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	//Wrap Blob in AttributeData and Response structures

	BlobDataUnmarshal := make(map[string]interface{})
	errUnmarshal := json.Unmarshal(blob.Data, &BlobDataUnmarshal)
	if errUnmarshal != nil {
		return
	}

	response := resources.BlobResponse{
		Data: resources.Blob{
			Key: resources.Key{
				ID:           fmt.Sprint(blob.Index),
				ResourceType: "Blob",
			},
			Attributes: resources.BlobAttributes{
				Obj: BlobDataUnmarshal,
			},
			Relationships: &resources.BlobRelationships{
				UserId: *blob.UserId,
			},
		},
	}

	w.WriteHeader(http.StatusCreated)
	ape.Render(w, &response)
}
