package cmd

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"
	"github.com/wjb-iv/rest-api-skeleton/rest"
)

// serveCmd represents the serve command which starts the HTTP server
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Run the service",
	Run: func(cmd *cobra.Command, args []string) {
		serve()
	},
}

func serve() {
	restServer := rest.New()
	go func() {
		err := restServer.Serve(httpPort, logRequests)
		if err != nil {
			log.Fatal(err)
		}
	}()
	// After setting everything up, wait for a SIGINT (perhaps triggered
	// by user with CTRL-C) Run cleanup when signal is received
	signalChan := make(chan os.Signal, 1)
	cleanupDone := make(chan bool)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)
	go func() {
		for _ = range signalChan {
			log.Println("Received an interrupt, stopping services...")
			cleanupDone <- true
		}
	}()
	<-cleanupDone
}

func init() {
	RootCmd.AddCommand(serveCmd)
	serveCmd.Flags().IntVar(&httpPort, "http", defaultHttpPort, "Port for the HTTP/REST service.")
	serveCmd.Flags().BoolVar(&logRequests, "log-requests", defaultReqLog, "Log all HTTP/REST requests.")
}
