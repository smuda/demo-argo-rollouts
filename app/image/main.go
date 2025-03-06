package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

var (
	opsGoodProcessed = promauto.NewCounter(prometheus.CounterOpts{
		Name: "cider_good_processed_apples_total",
		Help: "The good number of processed apples",
	})

	opsBadProcessed = promauto.NewCounter(prometheus.CounterOpts{
		Name: "cider_bad_processed_apples_total",
		Help: "The bad number of processed apples",
	})
)

func main() {
	server := &http.Server{
		Addr: ":8080",
	}

	// Expose metrics
	http.Handle("/metrics", promhttp.Handler())

	version := os.Getenv("VERSION")
	good := os.Getenv("GOOD")

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if good == "yes" {
			fmt.Printf("Environmentvar is %s, increasing good apples", good)
			opsGoodProcessed.Inc()
		} else {
			fmt.Printf("Environmentvar is %s, increasing bad apples", good)
			opsBadProcessed.Inc()
		}
		t := time.Now()
		fmt.Fprintf(w, "%s Hello, %q\n", t.Format("15:04:05"), version)
	})
	http.HandleFunc("/wait/", wait)

	log.Println("Starting http server on port 8080...")
	go func() {
		if err := server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("HTTP server error: %v", err)
		}
		log.Println("Stopped serving new connections on port 8080.")
	}()

	log.Println("Started web server")
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	shutdownCtx, shutdownRelease := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownRelease()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("HTTP shutdown error: %v", err)
	}
	log.Println("Graceful shutdown complete.")
}

func wait(w http.ResponseWriter, r *http.Request) {
	milliSecondsString := r.URL.Query().Get("ms")
	fmt.Fprintf(w, "Hi %s\n", milliSecondsString)

	log.Println("Ingress IP: ", r.RemoteAddr)

	if milliSecondsString == "" {
		fmt.Fprintf(w, "No wait time requested. Returning.")
		return
	}

	i, err := strconv.Atoi(milliSecondsString)
	if err != nil {
		// Ops. Handle error
		fmt.Fprintf(w, "Error while converting ms to int: %v", err)
		return
	}

	time.Sleep(time.Duration(i) * time.Millisecond)
	fmt.Fprintf(w, "Returning after %s ms\n", milliSecondsString)
}
