package main

import (
	"log/slog"
	"net/http"

	"github.com/telexintegrations/ekefan-go/server"
	"github.com/telexintegrations/ekefan-go/storage"
)

func main() {

	memory := storage.NewStorage()
	server := server.NewServer(memory)
	http.HandleFunc("/integration.json", server.IntegrationConfigHandler)
	http.HandleFunc("/tick", server.TickHandler)
	http.HandleFunc("/error-log", server.LogError)

	err := http.ListenAndServe(":8080", nil)
	slog.Error("Error starting ekefan-go gin APM service", "details", err)
}
