package server

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/telexintegrations/ekefan-go/model"
	"github.com/telexintegrations/ekefan-go/storage"
)

func TestTickHandler(t *testing.T) {
	//Test that tickHandler always returns a 202 http code on request
	ts := NewServer(storage.NewStorage())
	goodPayload := model.TelexRequestPayload{
		ChannelID: "sample-channel-id",
		ReturnURL: "sample-return-url",
		Settings:  []model.IntegrationSettings{},
	}
	data, _ := json.Marshal(goodPayload)
	reader := bytes.NewReader(data)
	tr := httptest.NewRequest(http.MethodPost, "/tick", reader)
	tr.Header.Set("Content-Type", "application/json")
	tw := httptest.NewRecorder()
	ts.TickHandler(tw, tr)
	assert.Equal(t, tw.Code, http.StatusAccepted)
	// purge store
	ts.Store.Purge(storage.WithTenant(context.Background(), goodPayload.ChannelID))
}
