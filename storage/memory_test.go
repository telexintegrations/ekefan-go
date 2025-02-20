package storage

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/telexintegrations/ekefan-go/model"
)

func TestWriteErrorLog_SucessfulWrite(t *testing.T) {
	m := NewMemory()
	msg := &model.TelexErrMsg{
		ErrMsg: "unit testing storage",
	}
	ctx := WithTenant(context.Background(), "unique-channel")
	var expectedErr error
	actalErr := m.WriteErrorLog(ctx, msg)

	assert.Equal(t, expectedErr, actalErr)
	expectedMsg := m.tenants["unique-channel"].ErrLogs[0]
	assert.Equal(t, expectedMsg, msg)
	delete(m.tenants, "unique-channel")
}

func TestWriteErrorLog_NoContextValueWithChannelID(t *testing.T) {
	m := NewMemory()
	msg := &model.TelexErrMsg{
		ErrMsg: "no context value error",
	}
	expectedErr := ErrTenantIDNotInContext
	actualErr := m.WriteErrorLog(context.Background(), msg)
	assert.Equal(t, expectedErr, actualErr)
	assert.Empty(t, m.tenants)
}

func TestReadErrorLog_SuccessfulRead(t *testing.T) {
	m := NewMemory()
	msg := &model.TelexErrMsg{
		ErrMsg: "unit testing storage",
	}
	ctx := WithTenant(context.Background(), "unique-channel")
	m.WriteErrorLog(ctx, msg)
	var expectedErr error
	errs, actualErr := m.ReadErrorLog(ctx)

	assert.Equal(t, expectedErr, actualErr)
	assert.Len(t, errs, 1)
	assert.Equal(t, errs[0], msg)
	delete(m.tenants, "unique-channel")
}

func TestReadErrorLog_NoContextValueWithChannelID(t *testing.T) {
	m := NewMemory()
	msg := &model.TelexErrMsg{
		ErrMsg: "unit testing storage",
	}
	ctx := WithTenant(context.Background(), "unique-channel")
	m.WriteErrorLog(ctx, msg)
	expectedErr := ErrTenantIDNotInContext
	errs, actualErr := m.ReadErrorLog(context.Background())

	assert.Equal(t, expectedErr, actualErr)
	assert.Empty(t, errs)
	delete(m.tenants, "unique-channel")
}

func TestReadErrorLog_TenantNotExist(t *testing.T) {
	m := NewMemory()

	expectedErr := ErrNotExist
	errs, actualErr := m.ReadErrorLog(WithTenant(context.Background(), "not-exit-channel"))

	assert.Equal(t, expectedErr, actualErr)
	assert.Empty(t, errs)
}

func TestPurge_Successful(t *testing.T) {
	m := NewMemory()
	msg := &model.TelexErrMsg{
		ErrMsg: "unit testing storage",
	}
	ctx := WithTenant(context.Background(), "unique-channel")
	m.WriteErrorLog(ctx, msg)
	var expectedErr error
	actualErr := m.Purge(ctx)

	assert.Equal(t, expectedErr, actualErr)
	assert.Empty(t, m.tenants["unique-channel"].ErrLogs)
	delete(m.tenants, "unique-channel")
}
func TestPurge_NoExist(t *testing.T) {
	m := NewMemory()
	expectedErr := ErrNotExist
	actualErr := m.Purge(WithTenant(context.Background(), "not-exist-channel"))

	assert.Equal(t, expectedErr, actualErr)
}

func TestPurge_NoContextValueWithChannelID(t *testing.T) {
	m := NewMemory()
	msg := &model.TelexErrMsg{
		ErrMsg: "unit testing storage",
	}
	ctx := WithTenant(context.Background(), "unique-channel")
	m.WriteErrorLog(ctx, msg)
	expectedErr := ErrTenantIDNotInContext
	actualErr := m.Purge(context.Background())

	assert.Equal(t, expectedErr, actualErr)
	delete(m.tenants, "unique-channel")
}
