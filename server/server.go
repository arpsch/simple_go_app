package server

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"golang.org/x/sys/unix"

	api "github.com/arpsch/ha/api/http"
	dns "github.com/arpsch/ha/dns"
)

// InitAndRun initializes the server and runs it
func InitAndRun(ctx context.Context, port, sid string) error {

	// Business logic object
	appl := dns.NewApp()

	os.Setenv(api.SECTOR_ID, sid)

	router, err := api.NewRouter(appl)
	if err != nil {
		return err
	}

	srv := &http.Server{
		Addr:    port,
		Handler: router,
	}

	go func() {
		log.Printf("starting the server, listening on %s\n", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, unix.SIGINT, unix.SIGTERM)
	<-quit

	log.Print("Server shutting down")

	ctxWithTimeout, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctxWithTimeout); err != nil {
		return err
	}

	log.Print("Server exited")
	return nil
}
