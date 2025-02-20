package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/telexintegrations/ekefan-go/model"
)

func TestIntegrationConfigHandler(t *testing.T) {
	var ts *Server

	date := model.IntegrationDates{
		CreatedAt: "2025-02-09",
		UpdatedAt: "2025-02-09",
	}
	descriptions := model.IntegrationDescriptions{
		AppName:         "ekefan-go GIN APM",
		AppDescription:  "Reports errors in gin applicaitons, request latency and tracing comming soon",
		AppURL:          baseUrl,
		AppLogo:         "https://i.imgur.com/IZqvffp.png",
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
			"Log Errors from gin applications",
		},
		IntegrationCategory: integrationCategory,
		Author:              "<The name of Your Organisation>",
		Settings:            settings,
		TickURL:             fmt.Sprintf("%s/tick", baseUrl),
		TargetURL:           "",
	}
	expectedJson := model.IntegrationConfig{
		Data: data,
	}
	tr := httptest.NewRequest(http.MethodGet, "/integration.json", nil)
	tr.Header.Set("Content-Type", "application/json")
	tw := httptest.NewRecorder()
	ts.IntegrationConfigHandler(tw, tr)

	var actualJson model.IntegrationConfig

	err := json.NewDecoder(tw.Body).Decode(&actualJson)
	assert.NoError(t, err)

	assert.EqualValues(t, expectedJson, actualJson)

}
