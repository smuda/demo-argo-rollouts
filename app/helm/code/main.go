package main

import (
  "context"
  "errors"
  "fmt"
  "log"
  "net/http"
  "os"
  "os/signal"
  "syscall"
  "strconv"
  "time"
)

func main() {
  server := &http.Server{
    Addr: ":8080",
  }

  version := os.Getenv("VERSION")

  http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
    t := time.Now()
    fmt.Fprintf(w, "%s Hello, %q\n",t.Format("15:04:05"), version)
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

func wait(w http.ResponseWriter, r *http.Request){
  milliSecondsString := r.URL.Query().Get("ms")
  fmt.Fprintf(w, "Hi %s\n", milliSecondsString)

  log.Println("Ingress IP: ", r.RemoteAddr)

  if milliSecondsString == "" {
    fmt.Fprintf(w, "No wait time requested. Returning.")
    return;
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
