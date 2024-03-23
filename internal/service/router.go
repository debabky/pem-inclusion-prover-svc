package service

import (
	"github.com/debabky/pem-inclusion-prover-svc/internal/service/handlers"
	"github.com/go-chi/chi"
	"gitlab.com/distributed_lab/ape"
)

func (s *service) router() chi.Router {
	r := chi.NewRouter()

	r.Use(
		ape.RecoverMiddleware(s.log),
		ape.LoganMiddleware(s.log),
		ape.CtxMiddleware(
			handlers.CtxLog(s.log),
			handlers.CtxIpfs(s.ipfs),
		),
	)
	r.Route("/integrations/pem-inclusion-prover-svc", func(r chi.Router) {
		r.Route("/v1", func(r chi.Router) {
			r.Post("/generate-merkle-tree", handlers.GenerateMerkleTree)
			r.Get("/generate-merkle-proof", handlers.CheckInclusion)
		})
	})

	return r
}
