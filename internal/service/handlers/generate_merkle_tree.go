package handlers

import (
	"bytes"
	"encoding/pem"
	"errors"
	"net/http"

	"github.com/debabky/pem-inclusion-prover-svc/internal/service/requests"
	"github.com/debabky/pem-inclusion-prover-svc/resources"

	"github.com/wealdtech/go-merkletree/v2"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
)

const PEM_BLOCK_TYPE = "CERTIFICATE"

func GenerateMerkleTree(w http.ResponseWriter, r *http.Request) {
	request, err := requests.NewCreateMerkleTreeRequest(r)
	if err != nil {
		Log(r).WithError(err).Error("Failed to parse the request")
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	pemBlocks, err := parsePemBlocks(request.Data.PemBlocks)
	if err != nil {
		Log(r).WithError(err).Error("Failed to parse the PEM blocks")
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	data := make([][]byte, 0)
	for _, block := range pemBlocks {
		data = append(data, block.Bytes)
	}

	tree, err := merkletree.NewTree(merkletree.WithData(data))
	if err != nil {
		Log(r).WithError(err).Error("Failed to construct a Merkle tree")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	marshaledTree, err := tree.MarshalJSON()
	if err != nil {
		Log(r).WithError(err).Error("Failed to marhsal the Merkle tree")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	cid, err := Ipfs(r).Add(bytes.NewReader(marshaledTree))
	if err != nil {
		Log(r).WithError(err).Error("Failed to save the Merkle tree to IPFS")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	ape.Render(w, resources.MerkleRootHash{
		Key: resources.Key{
			ID:   cid,
			Type: resources.MERKLE_ROOT_HASH,
		},
		Attributes: resources.MerkleRootHashAttributes{
			Hash: tree.String(),
		},
	})
}

func parsePemBlocks(rawPemBlocks []string) ([]*pem.Block, error) {
	pemBlocks := make([]*pem.Block, len(rawPemBlocks))

	for i, rawPemBlock := range rawPemBlocks {
		pemBlock, _ := pem.Decode([]byte(rawPemBlock))
		if pemBlock == nil || pemBlock.Type != PEM_BLOCK_TYPE {
			return nil, errors.New("failed to decode a pem block")
		}
		pemBlocks[i] = pemBlock
	}

	return pemBlocks, nil
}
