package handlers

import (
	"context"
	"net/http"

	shell "github.com/ipfs/go-ipfs-api"
	"gitlab.com/distributed_lab/logan/v3"
)

type ctxKey int

const (
	logCtxKey ctxKey = iota
	ipfsCtxKey
)

func CtxLog(entry *logan.Entry) func(context.Context) context.Context {
	return func(ctx context.Context) context.Context {
		return context.WithValue(ctx, logCtxKey, entry)
	}
}

func Log(r *http.Request) *logan.Entry {
	return r.Context().Value(logCtxKey).(*logan.Entry)
}

func CtxIpfs(entry *shell.Shell) func(context.Context) context.Context {
	return func(ctx context.Context) context.Context {
		return context.WithValue(ctx, ipfsCtxKey, entry)
	}
}

func Ipfs(r *http.Request) *shell.Shell {
	return r.Context().Value(ipfsCtxKey).(*shell.Shell)
}
