package server

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/telexintegrations/ekefan-go/model"
	"github.com/telexintegrations/ekefan-go/storage"
)

func TestSendErrorsToTelex(t *testing.T) {
	mockTelexServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))

		var received model.TelexResponsePayload
		err := json.NewDecoder(r.Body).Decode(&received)
		assert.NoError(t, err)

		expectedMessage := "unit testing storage"
		assert.Equal(t, expectedMessage, received.Message)

		w.WriteHeader(http.StatusOK)
	}))
	defer mockTelexServer.Close()

	mockStorage := storage.NewStorage()
	ts := &Server{
		Store:      mockStorage,
		HTTPClient: mockTelexServer.Client(), // Use mock HTTP client
	}

	msg := &model.TelexErrMsg{ErrMsg: "unit testing storage"}
	ctx := storage.WithTenant(context.Background(), "unique-channel")
	err := ts.Store.WriteErrorLog(ctx, msg)
	assert.NoError(t, err)

	telexReqPayload := model.TelexRequestPayload{
		ChannelID: "unique-channel",
		ReturnURL: mockTelexServer.URL,
		Settings:  nil,
	}

	fmt.Println(telexReqPayload.ReturnURL)
	ts.sendErrorsToTelex(ctx, telexReqPayload)

	logs, _ := ts.Store.ReadErrorLog(ctx)
	assert.Empty(t, logs, "Expected logs to be purged after sending")
}
