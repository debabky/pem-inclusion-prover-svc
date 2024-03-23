package requests

import (
	"encoding/json"
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/pkg/errors"
)

type CreateMerkleTreeRequestData struct {
	PemBlocks []string `json:"pem_blocks"`
}

type CreateMerkleTreeRequest struct {
	Data CreateMerkleTreeRequestData `json:"data"`
}

func NewCreateMerkleTreeRequest(r *http.Request) (CreateMerkleTreeRequest, error) {
	var request CreateMerkleTreeRequest

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return request, errors.Wrap(err, "failed to unmarshal")
	}

	return request, validateCreateMerkleTreeRequest(&request)
}

func validateCreateMerkleTreeRequest(r *CreateMerkleTreeRequest) error {
	return validation.Errors{
		"pem_blocks": validation.Validate(&r.Data.PemBlocks, validation.Required),
	}.Filter()
}
