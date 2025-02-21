package server

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/telexintegrations/ekefan-go/model"
)

// IntegrationConfigHandler returns an integration.json response for telex to setup the integration
func (s *Server) IntegrationConfigHandler(w http.ResponseWriter, r *http.Request) {
	handleCors(w, r)
	date := model.IntegrationDates{
		CreatedAt: "2025-02-09",
		UpdatedAt: "2025-02-09",
	}
	descriptions := model.IntegrationDescriptions{
		AppName:         "ekefan-go GIN APM",
		AppDescription:  "Reports errors in gin applicaitons. Request latency and tracing comming soon",
		AppURL:          baseUrl,
		AppLogo:         "https://img.icons8.com/?size=80&id=VUif5Y3XkX2o&format=png",
		BackgroundColor: "#fff",
	}
	settings := []model.IntegrationSettings{
		// {Label: "channel-id", Type: "text", Required: true, Default: ""},
		{Label: "interval", Type: "text", Required: true, Default: "* * * * *"},
	}

	data := model.IntegrationData{
		Date:            date,
		Descriptions:    descriptions,
		IsActive:        true,
		IntegrationType: integrationType,
		KeyFeatures: []string{
			"Notify channels subscribed about errors in gin applications",
		},
		IntegrationCategory: integrationCategory,
		Author:              "Cloud Ekefan",
		Settings:            settings,
		TickURL:             fmt.Sprintf("%s/tick", baseUrl),
		TargetURL:           "",
	}
	resp := model.IntegrationConfig{
		Data: data,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, "Failed to encode JSON", http.StatusInternalServerError)
	}
}
