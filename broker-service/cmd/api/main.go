package main

import (
	"fmt"
	"log"
	"net/http"
	"tools"
)

const port = 80

type Config struct {
	tools.Tools
}

func main() {
	app := Config{
		Tools: tools.New(),
	}

	log.Printf("Starting broker service on port %d\n", port)

	// Define http server
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: app.routes(),
	}

	// Start the server
	err := srv.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}
}
