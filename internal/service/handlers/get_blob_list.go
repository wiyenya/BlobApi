package handlers

import (
	resources "BlobApi/resources"
	"encoding/json"
	"fmt"
	"net/http"

	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
)

func GetBlobList(w http.ResponseWriter, r *http.Request) {

	connector := HorizonConnector(r)
	blobs, err := connector.GetBlobList()
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
	var responseData []resources.Blob
	for _, blob := range blobs {

		BlobDataUnmarshal := make(map[string]interface{})
		errUnmarshal := json.Unmarshal(blob.Data, &BlobDataUnmarshal)
		if errUnmarshal != nil {
			return
		}

		if blob == nil {
			fmt.Println("Warning: encountered nil blob")
			continue
		}
		if blob.UserId == nil {
			fmt.Println("Warning: UserID is nil for blob ID:", blob.Index)
			continue
		}
		responseData = append(responseData, resources.Blob{
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
		})
	}

	// Collect the response
	resp := resources.BlobListResponse{
		Data: responseData,
	}

	ape.Render(w, &resp)
}
