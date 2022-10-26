package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sha256service/internal/handler"
	"sha256service/internal/httpclient"
	"sha256service/pkg/sha256"
	"syscall"
)

func ContinuouslyHandleRequests(port *string, mux http.Handler) {
	log.Printf("start listening on port: %s", *port)
	log.Printf("http://localhost:%s", *port)
	if err := http.ListenAndServe(":"+*port, mux); err != nil {
		log.Fatalf("server listening error: %v", err)
	}
}

func main() {
	port := flag.String("port", "8080", "bing the HTTP server to this port")
	flag.Parse()
	ctx := context.Background()

	requestHandler := handler.NewRequestHandler(&handler.Config{
		Ctx:        ctx,
		HashClient: sha256.New(),
		HttpClient: httpclient.NewHTTPClient(ctx),
	})
	routesHandler := handler.GetRoutesHandler(requestHandler)

	go ContinuouslyHandleRequests(port, routesHandler)

	exitSignal := make(chan os.Signal)
	signal.Notify(exitSignal, syscall.SIGINT, syscall.SIGTERM)
	<-exitSignal
}
