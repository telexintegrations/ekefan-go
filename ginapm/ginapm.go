package ginamp

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Config struct {
	TelexChanID string `json:"telex_channel_id"`
}

const serverURL = "https://ekefan-go.onrender.com/error-log"

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
		fmt.Println("Failed to marshal error payload:", err)
		return
	}

	// Send errors to APM server
	resp, err := http.Post(serverURL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("Failed to send error log:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusAccepted {
		fmt.Println("APM server responded with status:", resp.Status)
	}
}
