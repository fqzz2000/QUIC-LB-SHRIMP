package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/fqzz2000/QUIC-LB-SHRIMP/pkg/lb"
)

// Configuration flags
var (
	configFile string
	listenAddr string
	debugMode  bool
)

func init() {
	// Parse command line flags
	flag.StringVar(&configFile, "config", "config.yaml", "Path to configuration file")
	flag.StringVar(&listenAddr, "listen", ":8080", "Address to listen on")
	flag.BoolVar(&debugMode, "debug", false, "Enable debug mode")
}

func main() {
	// Parse flags
	flag.Parse()

	// Initialize logger
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	// Load configuration
	// TODO: Implement configuration loading
	backends := []string{"backend1", "backend2", "backend3"} // TODO: Implement configuration loading

	// Initialize load balancer
	lb, err := lb.InitLoadBalancer(listenAddr, backends)
	if err != nil {
		log.Fatalf("Failed to initialize load balancer: %v", err)
	}

	// Setup signal handling for graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Start the load balancer
	go func() {
		if err := lb.Start(); err != nil {
			log.Fatalf("Load balancer error: %v", err)
		}
	}()

	log.Printf("QUIC Load Balancer started on %s", listenAddr)

	// Wait for shutdown signal
	<-sigChan
	log.Println("Shutting down...")

	// Perform cleanup
	if err := lb.Shutdown(); err != nil {
		log.Printf("Error during shutdown: %v", err)
	}
}
