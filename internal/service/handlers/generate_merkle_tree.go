package handlers

import (
	"bytes"
	"encoding/pem"
	"net/http"

	"github.com/debabky/pem-inclusion-prover-svc/internal/service/requests"
	"github.com/debabky/pem-inclusion-prover-svc/resources"
	"github.com/rarimo/certificate-transparency-go/x509"
	"github.com/wealdtech/go-merkletree/v2"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

const PEM_BLOCK_TYPE = "CERTIFICATE"

func GenerateMerkleTree(w http.ResponseWriter, r *http.Request) {
	request, err := requests.NewCreateMerkleTreeRequest(r)
	if err != nil {
		Log(r).WithError(err).Error("Failed to parse the request")
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	certificates, err := parsePemBlocks(request.Data.PemBlocks)
	if err != nil {
		Log(r).WithError(err).Error("Failed to parse the PEM blocks")
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	data := make([][]byte, 0)
	for _, certificate := range certificates {
		data = append(data, certificate.RawSubjectPublicKeyInfo)
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

func parsePemBlocks(rawPemBlocks []string) ([]*x509.Certificate, error) {
	certificates := make([]*x509.Certificate, len(rawPemBlocks))

	for i, rawPemBlock := range rawPemBlocks {

		pemBlock, _ := pem.Decode([]byte(rawPemBlock))
		if pemBlock == nil || pemBlock.Type != PEM_BLOCK_TYPE {
			return nil, errors.From(errors.New("failed to decode a pem block"), logan.F{
				"index":     i,
				"pem_block": rawPemBlock,
			})
		}

		x509, err := x509.ParseCertificate(pemBlock.Bytes)
		if err != nil {
			return nil, errors.From(err, logan.F{
				"block_number": i,
				"pem_block":    rawPemBlock,
			})
		}

		certificates[i] = x509
	}

	return certificates, nil
}
