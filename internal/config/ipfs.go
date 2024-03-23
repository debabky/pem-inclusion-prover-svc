package config

import (
	"github.com/pkg/errors"
	"gitlab.com/distributed_lab/figure"
	"gitlab.com/distributed_lab/kit/comfig"
	"gitlab.com/distributed_lab/kit/kv"
)

type IpfsConfiger interface {
	IpfsConfig() *IpfsConfig
}

type IpfsConfig struct {
	Url string `fig:"url,required"`
}

type ipfs struct {
	once   comfig.Once
	getter kv.Getter
}

func NewIpfsConfiger(getter kv.Getter) IpfsConfiger {
	return &ipfs{
		getter: getter,
	}
}

func (c *ipfs) IpfsConfig() *IpfsConfig {
	return c.once.Do(func() interface{} {
		var result IpfsConfig
		err := figure.
			Out(&result).
			With(figure.BaseHooks).
			From(kv.MustGetStringMap(c.getter, "ipfs")).
			Please()

		if err != nil {
			panic(errors.WithMessage(err, "failed to figure out"))
		}

		return &result
	}).(*IpfsConfig)
}
