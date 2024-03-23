package service

import (
	"net"
	"net/http"

	"github.com/debabky/pem-inclusion-prover-svc/internal/config"
	"gitlab.com/distributed_lab/kit/copus/types"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"

	shell "github.com/ipfs/go-ipfs-api"
)

type service struct {
	log      *logan.Entry
	copus    types.Copus
	listener net.Listener
	ipfs     *shell.Shell

	cfg config.Config
}

func (s *service) run() error {
	s.log.Info("Service started")
	r := s.router()

	if err := s.copus.RegisterChi(r); err != nil {
		return errors.Wrap(err, "cop failed")
	}

	return http.Serve(s.listener, r)
}

func newService(cfg config.Config) *service {

	return &service{
		log:      cfg.Log(),
		copus:    cfg.Copus(),
		listener: cfg.Listener(),
		ipfs:     shell.NewShell(cfg.IpfsConfig().Url),

		cfg: cfg,
	}
}

func Run(cfg config.Config) {
	if err := newService(cfg).run(); err != nil {
		panic(err)
	}
}
