package server

import (
	"encoding/json"
	"net/http"

	"github.com/telexintegrations/ekefan-go/model"
	"github.com/telexintegrations/ekefan-go/storage"
)

// TickHandler handles telex tick request,
// sends errors to telex payload channel id
// and purges the memory to free up more space
func (s *Server) TickHandler(w http.ResponseWriter, r *http.Request) {
	handleCors(w, r)

	var telexPayload model.TelexRequestPayload
	err := json.NewDecoder(r.Body).Decode(&telexPayload)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Create context with ChannelID
	ctx := storage.WithTenant(r.Context(), telexPayload.ChannelID)

	// Run monitorTasks asynchronously
	go s.sendErrorsToTelex(ctx, telexPayload)

	w.WriteHeader(http.StatusAccepted)
}
