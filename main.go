package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"golang.org/x/sync/errgroup"
)

func main() {
	workingDirectory, err := os.Getwd()
	if err != nil {
		log.Fatalf("could not get current working directory: %s", err)
	}
	log.Printf("Current directory: %s", workingDirectory)
	host := getHost()
	log.Printf("Host: %s", host)
	port := getPort()
	log.Printf("Port: %s", port)

	mux := http.NewServeMux()
	fileHandler := http.StripPrefix("/", http.FileServer(http.Dir(workingDirectory)))
	mux.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("request: %s", r.URL)
		fileHandler.ServeHTTP(w, r)
	}))

	server := http.Server{
		Addr:    net.JoinHostPort(host, port),
		Handler: mux,
	}

	eg := errgroup.Group{}
	eg.SetLimit(2)

	eg.Go(func() error {
		signalChan := make(chan os.Signal)
		defer close(signalChan)
		signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)
		defer signal.Reset(syscall.SIGTERM)
		<-signalChan
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
		defer cancel()
		if err := server.Shutdown(ctx); err != nil {
			return fmt.Errorf("an error occurred while attempting to handle a signal to shutdown: %w", err)
		}
		return nil
	})

	eg.Go(func() error {
		if err := server.ListenAndServe(); errors.Is(err, http.ErrServerClosed) || err == nil {
			return nil
		} else {
			return err
		}
	})

	if err := eg.Wait(); err != nil {
		log.Fatalf("an unexpected error has occurred: %s", err)
	}
	log.Println("server has gracefully shutdown")
}

func getHost() string {
	host, isPresent := os.LookupEnv("HOST")
	if !isPresent {
		return "0.0.0.0"
	}
	return host
}

func getPort() string {
	port, isPresent := os.LookupEnv("PORT")
	if !isPresent {
		return "8080"
	}
	return port
}
