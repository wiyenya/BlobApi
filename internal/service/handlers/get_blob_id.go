package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	requests "BlobApi/internal/service/requests"

	resources "BlobApi/resources"

	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
)

func GetBlobID(w http.ResponseWriter, r *http.Request) {

	id, err := requests.DecodeGetBlobRequest(r)
	if err != nil || id < 1 {
		Log(r).WithError(err).Error("Invalid ID")
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	connector := HorizonConnector(r)

	//Retrieve record by ID
	blob, err := connector.Get(id)
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

	ape.Render(w, &resp)
}
