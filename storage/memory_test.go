package storage

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/telexintegrations/ekefan-go/model"
)

func TestMemory_WriteErrorLog(t *testing.T) {
	m := NewMemory()
	msg := &model.TelexErrMsg{
		ErrMsg: "unit testing storage",
	}
	ctx := WithTenant(context.Background(), "unique-channel")
	var expectedErr error
	actalErr := m.WriteErrorLog(ctx, msg)

	assert.Equal(t, expectedErr, actalErr)
}
