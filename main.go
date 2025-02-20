package main

import (
	"github.com/telexintegrations/ekefan-go/server"
	"log/slog"
	"net/http"
)

func main() {

	server := server.NewServer()
	http.HandleFunc("/integration.json", server.IntegrationConfigHandler)
	http.HandleFunc("/tick", server.TickHandler)
	http.HandleFunc("/error-log", server.LogError)

	err := http.ListenAndServe(":8080", nil)
	slog.Error("Error starting ekefan-go gin APM service", "details", err)
}
