package requests

import (
	"encoding/base32"
	"encoding/json"
	"fmt"
	"hash"
	"net/http"

	"github.com/go-chi/chi"
	validation "github.com/go-ozzo/ozzo-validation"
	"gitlab.com/distributed_lab/logan/v3/errors"
	// "gitlab.com/tokend/api/internal/api/resources"
	// "gitlab.com/tokend/api/internal/types"
	// "gitlab.com/tokend/go/hash"
)

func CreateBlobRequest(r *http.Request) (resources.BlobRequest, error) {
	var request resources.BlobRequestResponse

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return request.Data, errors.Wrap(err, "failed to unmarshal")
	}

	// deciding on owner address, relation ID has precedence
	// TODO(v4) remove url param
	if request.Data.Relationships.Owner.Data == nil {
		request.Data.Relationships.Owner.Data = &resources.Key{ID: chi.URLParam(r, "address")}
	}

	return request.Data, ValidateCreateBlobRequest(request.Data)
}

func ValidateCreateBlobRequest(r resources.BlobRequest) error {
	return validation.Errors{
		"/data/type":                        validation.Validate(&r.Type, validation.Required),
		"/data/attributes/value":            validation.Validate(&r.Attributes.Value, validation.Required),
		"/data/relationships/owner/data/id": validation.Validate(&r.Relationships.Owner.Data.ID, validation.Required),
	}.Filter()
}

func Blob(r resources.BlobRequest, allowedRole uint64) (*types.Blob, error) {
	blob, err := types.GetBlobType(string(r.Type))
	if err != nil {
		return nil, errors.Wrap(err, "failed to create blob")
	}
	msg := fmt.Sprintf("%s%d%s", r.Relationships.Owner.Data.ID, blob, r.Attributes.Value)
	hash := hash.Hash([]byte(msg))
	owner := types.Address(r.Relationships.Owner.Data.ID)

	return &types.Blob{
		ID:                base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(hash[:]),
		Type:              blob,
		Value:             r.Attributes.Value,
		CreatorSignerRole: &allowedRole,
		OwnerAddress:      &owner,
	}, nil
}
