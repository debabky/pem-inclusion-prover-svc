package handlers

import (
	"bytes"
	"crypto/x509"
	"encoding/hex"
	"encoding/json"
	"encoding/pem"
	"errors"
	"net/http"

	"github.com/debabky/pem-inclusion-prover-svc/internal/service/requests"
	"github.com/debabky/pem-inclusion-prover-svc/resources"
	"github.com/wealdtech/go-merkletree/v2"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
)

func CheckInclusion(w http.ResponseWriter, r *http.Request) {
	request, err := requests.NewGenerateProofRequest(r)
	if err != nil {
		Log(r).WithError(err).Error("Failed to parse the request")
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	data, err := Ipfs(r).Cat(request.Data.MerkleTreeId)
	if err != nil {
		Log(r).WithError(err).Error("Failed to get the merkle tree from the ipfs")
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	buf := new(bytes.Buffer)
	buf.ReadFrom(data)
	var tree merkletree.MerkleTree
	if err := json.Unmarshal(buf.Bytes(), &tree); err != nil {
		Log(r).WithError(err).Error("Failed to unmarshal the tree")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	pemBlock, _ := pem.Decode([]byte(request.Data.PemBlock))
	if pemBlock == nil {
		Log(r).Error("Failed to parse the pem block")
		ape.RenderErr(w, problems.BadRequest(errors.New("failed to parse the pem block"))...)
		return
	}

	cert, err := x509.ParseCertificate(pemBlock.Bytes)
	if err != nil {
		Log(r).WithError(err).Error("Failed to convert a PEM block to a certificate")
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	proof, err := tree.GenerateProof(cert.Raw, 0)
	if err != nil {
		Log(r).WithError(err).Error("Pem block is not present in the tree")
		ape.RenderErr(w, problems.BadRequest(errors.New("pem block is not present in the tree"))...)
		return
	}

	hashes := make([]string, len(proof.Hashes))
	for i, hash := range proof.Hashes {
		hashes[i] = hex.EncodeToString(hash)
	}

	ape.Render(w, resources.InclusionProof{
		Key: resources.Key{
			ID:   tree.String(),
			Type: resources.INCLUSION_PROOF,
		},
		Attributes: resources.InclusionProofAttributes{
			Hashes: hashes,
			Index:  proof.Index,
		},
	})
}
