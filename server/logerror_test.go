package server

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/telexintegrations/ekefan-go/storage"
)

func TestLogError(t *testing.T) {
	ts := NewServer(storage.NewStorage())

	testErrPayload := ErrorPayload{
		TelexChanID: "sample-channel-id",
		Errors:      []string{"error-1", "error-2"},
	}

	data, _ := json.Marshal(testErrPayload)
	reader := bytes.NewReader(data)
	tr := httptest.NewRequest(http.MethodPost, "/error-log", reader)
	tr.Header.Set("Content-Type", "application/json")
	tw := httptest.NewRecorder()
	ts.LogError(tw, tr)
	assert.Equal(t, tw.Code, http.StatusAccepted)
	// purge store
	ts.Store.Purge(storage.WithTenant(context.Background(), testErrPayload.TelexChanID))
}
