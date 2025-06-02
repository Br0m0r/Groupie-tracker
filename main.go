package main

import (
	"fmt"
	"log"
	"net/http"

	"groupie/config"
	"groupie/server"
)

func main() {


	// Configure and start HTTP server using config values
	mux := server.SetupServer()
	httpServer := &http.Server{
		Addr:         config.PORT,
		Handler:      mux,
		ReadTimeout:  config.SERVER_READ_TIMEOUT,
		WriteTimeout: config.SERVER_WRITE_TIMEOUT,
		IdleTimeout:  config.SERVER_IDLE_TIMEOUT,
	}

	fmt.Printf("Server starting on http://localhost%s\n", config.PORT)
	fmt.Println("Press Ctrl+C to stop the server")

	if err := httpServer.ListenAndServe(); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}