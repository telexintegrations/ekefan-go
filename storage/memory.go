package storage

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"sync"

	"github.com/telexintegrations/ekefan-go/model"
)

type Memory struct {
	sync.RWMutex
	tenants map[string]Tenant
}

type Tenant struct {
	ChannelID string
	ErrLogs   []*model.TelexErrMsg
}

// tenantKeyType is a custom type for the key "tenant", following context.Context convention
type tenantKeyType string

var (
	ErrNotExist             = errors.New("tenant does not exist")
	ErrTenantIDNotInContext = errors.New("no channel id provided in context")
)

const (
	// tenantKey holds tenancy for spans
	tenantKey = tenantKeyType("channel_id")
)

func NewMemory() *Memory {
	return &Memory{
		tenants: make(map[string]Tenant),
	}
}

// WithTenant creates a Context with a tenant association
func WithTenant(ctx context.Context, tenant string) context.Context {
	return context.WithValue(ctx, tenantKey, tenant)
}

// GetErrorLog reads the tenant value from ctx, and returns the tenant errLog
func (m *Memory) ReadErrorLog(ctx context.Context) ([]*model.TelexErrMsg, error) {
	m.RWMutex.RLock()
	defer m.RWMutex.RUnlock()
	ctxValue := ctx.Value(tenantKey)
	telexChanID, ok := ctxValue.(string)
	if !ok {
		return nil, fmt.Errorf("%s, %s", ErrTenantIDNotInContext, telexChanID)
	}

	tenant, ok := m.tenants[telexChanID]
	if !ok {
		return nil, ErrNotExist
	}
	return tenant.ErrLogs, nil
}

// WriteError writes a telexErr a specific tenant, returns nil on success
// telexErr must not be nil or empty
func (m *Memory) WriteErrorLog(ctx context.Context, telexErr *model.TelexErrMsg) error {
	ctxValue := ctx.Value(tenantKey)
	telexChanID, ok := ctxValue.(string)
	if !ok {
		return ErrTenantIDNotInContext
	}
	if telexErr == nil {
		slog.Info("Nil telexErr received but not written")
		return nil
	}
	m.RWMutex.Lock()
	defer m.RWMutex.Unlock()

	tenant, ok := m.tenants[telexChanID]
	if !ok {
		m.tenants[telexChanID] = Tenant{
			ChannelID: telexChanID,
			ErrLogs: []*model.TelexErrMsg{
				telexErr,
			},
		}
		return nil
	}

	tenant.ErrLogs = append(tenant.ErrLogs, telexErr)
	m.tenants[telexChanID] = tenant

	return nil
}

// Purge cleans in-memory stuff **must be called only when telex makes a tick request
func (m *Memory) Purge(ctx context.Context) error {
	ctxValue := ctx.Value(tenantKey)
	telexChanID, ok := ctxValue.(string)
	if !ok {
		return ErrTenantIDNotInContext
	}

	m.RWMutex.Lock()
	defer m.RWMutex.Unlock()

	tenant, ok := m.tenants[telexChanID]
	if !ok {
		return ErrNotExist
	}

	tenant.ErrLogs = tenant.ErrLogs[:0]
	m.tenants[telexChanID] = tenant
	return nil
}
