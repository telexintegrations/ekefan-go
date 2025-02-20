package server

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/telexintegrations/ekefan-go/model"
	"github.com/telexintegrations/ekefan-go/storage"
)

const (
	telexGinUsername    = "ekefan-go Gin APM test"
	telexGinEventName   = "Error Log"
	telexGinErrorStatus = "error"
	jsonAppType         = "application/json"
	baseUrl             = "https://ekefan-go.onrender.com"
	integrationType     = "interval"
	integrationCategory = "Performance Monitoring"
)

type Server struct {
	Store      storage.Store
	HTTPClient *http.Client
}

func NewServer(store storage.Store) *Server {
	return &Server{
		Store:      store,
		HTTPClient: &http.Client{Timeout: 10 * time.Second},
	}
}

func handleCors(w http.ResponseWriter, r *http.Request) {
	origin := r.Header.Get("Origin")
	if origin == "https://telex.im/" {
		w.Header().Set("Access-Control-Allow-Origin", "https://telex.im/")
	} else {
		w.Header().Set("Access-Control-Allow-Origin", "*")
	}

	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Access-Control-Allow-Methods", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.WriteHeader(http.StatusOK)
}

func (s *Server) sendErrorsToTelex(ctx context.Context, payload model.TelexRequestPayload) {
	errLogs, err := s.Store.ReadErrorLog(ctx)
	if err != nil {
		fmt.Printf("failed to read error logs: %v\n", err)
		return
	}

	if len(errLogs) == 0 {
		fmt.Println("No error logs to report.")
		return
	}

	// http.Client{Timeout: 10 * time.Second}
	client := s.HTTPClient

	for _, errLog := range errLogs {
		respPayload := model.TelexResponsePayload{
			Message:   errLog.ErrMsg,
			Username:  telexGinUsername,
			EventName: telexGinEventName,
			Status:    telexGinErrorStatus,
		}

		data, err := json.Marshal(respPayload)
		if err != nil {
			fmt.Printf("error marshalling struct: %v\n", err)
			continue
		}
		reader := bytes.NewReader(data)

		res, err := client.Post(payload.ReturnURL, "application/json", reader)
		if err != nil {
			fmt.Printf("error sending POST request to Telex: %v\n", err)
			continue
		}
		defer res.Body.Close()

		fmt.Printf("Sent error log: %s, Response status: %s\n", errLog.ErrMsg, res.Status)
	}

	// Purge logs after processing
	if err := s.Store.Purge(ctx); err != nil {
		fmt.Printf("failed to purge error logs: %v\n", err)
	}
}
