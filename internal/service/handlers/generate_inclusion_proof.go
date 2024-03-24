package handlers

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"encoding/pem"
	"errors"
	"math/big"
	"net/http"

	"github.com/debabky/go-merkletree/v2"
	"github.com/debabky/go-merkletree/v2/poseidon"
	"github.com/debabky/pem-inclusion-prover-svc/internal/service/requests"
	"github.com/debabky/pem-inclusion-prover-svc/resources"
	"github.com/rarimo/certificate-transparency-go/x509"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
)

func GenerateInclusionProof(w http.ResponseWriter, r *http.Request) {
	request, err := requests.NewGenerateProofRequest(r)
	if err != nil {
		Log(r).WithError(err).Error("Failed to parse the request")
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	data, err := Ipfs(r).Cat(request.Data.MerkleTreeId)
	if err != nil {
		Log(r).WithError(err).Error("Failed to get the Merkle tree from the ipfs")
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	buf := new(bytes.Buffer)
	buf.ReadFrom(data)
	var tree merkletree.MerkleTree
	if err := json.Unmarshal(buf.Bytes(), &tree); err != nil {
		Log(r).WithError(err).Error("Failed to unmarshal the Merkle tree")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	bytes := []byte{}
	switch request.Encoding {
	case "raw":
		pemBlock, _ := pem.Decode([]byte(request.Data.PemBlock))
		if pemBlock == nil || pemBlock.Type != PEM_BLOCK_TYPE {
			Log(r).Error("Failed to parse the pem block")
			ape.RenderErr(w, problems.BadRequest(errors.New("failed to parse the pem block"))...)
			return
		}

		x509, err := x509.ParseCertificate(pemBlock.Bytes)
		if err != nil {
			Log(r).WithError(err).Error("Failed to parse the PEM block to a certificate")
			ape.RenderErr(w, problems.BadRequest(err)...)
			return
		}

		hash := poseidon.New()
		bytes = hash.Hash(x509.RawSubjectPublicKeyInfo[:HASH_DATA_LENGTH])
	case "hex":
		bytes, err = hex.DecodeString(request.Data.PemBlock)
		if err != nil {
			Log(r).WithError(err).Error("Failed to parse HEX PEM block")
			ape.RenderErr(w, problems.BadRequest(errors.New("failed to parse HEX PEM block"))...)
			return
		}
	case "decimal":
		bigInt, ok := big.NewInt(0).SetString(request.Data.PemBlock, 10)
		if !ok {
			Log(r).Error("Failed to parse decimal PEM block")
			ape.RenderErr(w, problems.BadRequest(errors.New("failed to parse decimal PEM block"))...)
			return
		}
		bytes = bigInt.Bytes()
	}

	proof, err := tree.GenerateProof(bytes, 0)
	if err != nil {
		Log(r).WithError(err).Error("The pem block is not present in the tree")
		ape.RenderErr(w, problems.BadRequest(errors.New("the pem block is not present in the tree"))...)
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
