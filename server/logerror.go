package server

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/telexintegrations/ekefan-go/model"
	"github.com/telexintegrations/ekefan-go/storage"
)

type ErrorPayload struct {
	TelexChanID string                 `json:"telex_channel_id"`
	Errors      []string               `json:"errors"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
}

// LogError handles requests from client applications, it stores the errors until telex reads them
func (s *Server) LogError(w http.ResponseWriter, r *http.Request) {
	var payload ErrorPayload

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if payload.TelexChanID == "" || len(payload.Errors) == 0 {
		http.Error(w, "Missing channel ID or errors", http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusAccepted)

	ctx := storage.WithTenant(context.Background(), payload.TelexChanID)

	for _, errMsg := range payload.Errors {
		if errMsg == "" {
			continue
		}
		errLog := &model.TelexErrMsg{
			ErrMsg: errMsg,
		}
		if err := s.Store.WriteErrorLog(ctx, errLog); err != nil {
			http.Error(w, fmt.Sprintf("Failed to write error: %v", err), http.StatusInternalServerError)
			return
		}
		fmt.Printf("Logged error for %s: %s\n", payload.TelexChanID, errMsg)
	}

	fmt.Fprintln(w, "Error log accepted")
}
