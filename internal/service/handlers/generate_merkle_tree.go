package handlers

import (
	"bytes"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"net/http"

	"github.com/debabky/pem-inclusion-prover-svc/internal/service/requests"
	"github.com/debabky/pem-inclusion-prover-svc/resources"

	"github.com/wealdtech/go-merkletree/v2"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
)

func GenerateMerkleTree(w http.ResponseWriter, r *http.Request) {
	request, err := requests.NewCreateMerkleTreeRequest(r)
	if err != nil {
		Log(r).WithError(err).Error("Failed to parse the request")
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	pemBlocks, err := parsePemBlocks(request.Data.PemBlocks)
	if err != nil {
		Log(r).WithError(err).Error("Failed to parse PEM blocks")
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	data := make([][]byte, 0)
	for _, block := range pemBlocks {
		cert, err := x509.ParseCertificate(block.Bytes)
		if err != nil {
			Log(r).WithError(err).Error("Failed to convert a PEM block to a certificate")
			ape.RenderErr(w, problems.BadRequest(err)...)
			return
		}
		data = append(data, cert.Raw)
	}

	tree, err := merkletree.NewTree(merkletree.WithData(data))
	if err != nil {
		Log(r).WithError(err).Error("Failed to construct a Merkle tree")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	marshaledTree, err := tree.MarshalJSON()
	if err != nil {
		Log(r).WithError(err).Error("Failed to marhsale a Merkle tree")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	cid, err := Ipfs(r).Add(bytes.NewReader(marshaledTree))
	if err != nil {
		Log(r).WithError(err).Error("Failed to save the tree to IPFS")
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

func parsePemBlocks(pemBlocks []string) ([]*pem.Block, error) {
	blocks := make([]*pem.Block, len(pemBlocks))

	for i, pem_block := range pemBlocks {
		pem, _ := pem.Decode([]byte(pem_block))
		if pem == nil {
			return nil, errors.New("failed to decode a pem block")
		}
		blocks[i] = pem
	}

	return blocks, nil
}
