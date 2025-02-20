package ginamp

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Config holds required fields to accurately your error logs to telex
type Config struct {
	// TelexChanID - the id of the telex channel where logs pop up
	TelexChanID string `json:"telex_channel_id"`
	// ApmServerUrl - the url of the apm server sends logs to telex at intervals
	// currently: "https://ekefan-go.onrender.com/error-log"
	ApmServerUrl string `json:"apm_server_url"`
}

// GinAPM creates middleware to capture errors and send them to the APM server
func GinAPM(config Config) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		defer sendErrors(ctx, config)
		ctx.Next()
	}
}

// sendErrors sends captured errors to the APM server
func sendErrors(ctx *gin.Context, config Config) {
	if len(ctx.Errors) == 0 {
		return
	}

	// Prepare error payload
	payload := map[string]interface{}{
		"telex_channel_id": config.TelexChanID,
		"errors":           ctx.Errors.Errors(),
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		slog.Error("Failed to marshal error payload", "detail", err.Error())
		return
	}

	// Send errors to APM server
	resp, err := http.Post(config.ApmServerUrl, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		slog.Error("Failed to send error log", "details", err.Error())
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusAccepted {
		slog.Error("APM server responded", "status", resp.Status)
	}
}
