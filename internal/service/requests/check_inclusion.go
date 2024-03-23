package requests

import (
	"encoding/json"
	"net/http"

	"github.com/pkg/errors"
)

type GenerateProofRequestData struct {
	PemBlock     string `json:"pem_block"`
	MerkleTreeId string `json:"merkle_tree_id"`
}

type GenerateProofRequest struct {
	Data GenerateProofRequestData `json:"data"`
}

func NewGenerateProofRequest(r *http.Request) (GenerateProofRequest, error) {
	var request GenerateProofRequest

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return request, errors.Wrap(err, "failed to unmarshal")
	}

	return request, nil
}
