package handlers

import (
	"BlobApi/internal/data"
	"BlobApi/internal/service/requests"
	"BlobApi/resources"
	"encoding/json"
	"fmt"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"net/http"
	"time"
)

func CreateBlob(w http.ResponseWriter, r *http.Request) {
	// Decoding the request body

	req, err := requests.DecodeCreateBlobRequest(r)
	if err != nil {
		Log(r).WithError(err).Error("BadRequest")
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	connector := HorizonConnector(r)

	id := req.Relationships.UserId

	//data to JSON (map to bytes)
	jsonData, errMarshal := json.Marshal(req.Attributes.Value)
	if errMarshal != nil {
		return
	}

	// Inserting a blob
	blobId, err := connector.Insert(id, jsonData)
	if err != nil {
		ape.RenderErr(w, problems.InternalError())
		return
	}

	//Getting a blob to return the created resource

	var blob *data.Blob
	for {
		blobFromGet, _ := connector.Get(blobId)
		if blobFromGet != nil {
			blob = blobFromGet
			break
		}
		time.Sleep(time.Second)
	}

	// Wrap Blob in AttributeData and Response structures

	BlobDataUnmarshal := make(map[string]interface{})
	errUnmarshal := json.Unmarshal(blob.Data, &BlobDataUnmarshal)
	if errUnmarshal != nil {
		return
	}

	resp := resources.BlobResponse{
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
	ape.Render(w, &resp)
}
