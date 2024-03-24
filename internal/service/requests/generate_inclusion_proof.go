package requests

import (
	"encoding/json"
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/pkg/errors"
)

type GenerateProofRequestData struct {
	PemBlock     string `json:"pem_block"`
	MerkleTreeId string `json:"merkle_tree_id"`
}

type GenerateProofRequest struct {
	Data     GenerateProofRequestData `json:"data"`
	Encoding string
}

func NewGenerateProofRequest(r *http.Request) (GenerateProofRequest, error) {
	var request GenerateProofRequest

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return request, errors.Wrap(err, "failed to unmarshal")
	}

	request.Encoding = r.URL.Query().Get("encoding")

	return request, validateGenerateProofRequest(&request)
}

func validateGenerateProofRequest(r *GenerateProofRequest) error {
	return validation.Errors{
		"pem_block":      validation.Validate(&r.Data.PemBlock, validation.Required),
		"merkle_tree_id": validation.Validate(&r.Data.MerkleTreeId, validation.Required),
		"encoding":       validation.Validate(&r.Encoding, validation.Required, validation.In("raw", "hex", "decimal")),
	}.Filter()
}
