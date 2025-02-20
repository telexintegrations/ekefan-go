package storage

import (
	"context"

	"github.com/telexintegrations/ekefan-go/model"
)

type Store interface {
	ReadErrorLog(ctx context.Context) ([]*model.TelexErrMsg, error)
	WriteErrorLog(ctx context.Context, telexErr *model.TelexErrMsg) error
	Purge(ctx context.Context) error
}

// v1 of gin-apm uses only in-memory storage backend
func NewStorage() Store {
	return &Memory{
		tenants: map[string]Tenant{},
	}
}
