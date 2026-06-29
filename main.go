package main

import (
	"cmp"
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	stats_api "github.com/fukata/golang-stats-api-handler"
	"golang.org/x/sync/errgroup"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "text/plain; charset=utf-8")
		_, _ = fmt.Fprintf(w, "Hello, %q", r.URL.Path)
	})
	mux.HandleFunc("/api/stats", stats_api.Handler)

	server := &http.Server{
		Addr:    ":" + cmp.Or(os.Getenv("SERVER_PORT"), "8080"),
		Handler: accesslog(mux),
	}

	if err := graceful(context.Background(), server); err != nil {
		slog.Error("server error", "error", err)
		os.Exit(1)
	}
}

func graceful(ctx context.Context, server *http.Server) error {
	listener, err := net.Listen("tcp", server.Addr)
	if err != nil {
		return fmt.Errorf("failed to listen: %w", err)
	}

	ctx, stop := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	eg, ctx := errgroup.WithContext(ctx)
	eg.Go(func() error {
		slog.Info("starting the server", "address", listener.Addr().String())
		if err := server.Serve(listener); err != nil && !errors.Is(err, http.ErrServerClosed) {
			return fmt.Errorf("failed to serve: %w", err)
		}
		return nil
	})

	eg.Go(func() error {
		<-ctx.Done()
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		slog.Info("shutting down the server")
		if err := server.Shutdown(ctx); err != nil {
			return fmt.Errorf("failed to shutdown the server: %w", err)
		}
		return nil
	})

	if err := eg.Wait(); err != nil {
		return err
	}
	slog.Info("server shutdown completed")
	return nil
}

func accesslog(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rw := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}

		defer func() {
			slog.Info("access",
				slog.String("protocol", r.Proto),
				slog.String("host", r.Host),
				slog.String("method", r.Method),
				slog.String("path", r.URL.Path),
				slog.String("user_agent", r.UserAgent()),
				slog.Int("status_code", rw.statusCode),
				slog.Int64("request_size", r.ContentLength),
				slog.Int64("response_size", rw.size),
				slog.Float64("response_time", time.Since(start).Seconds()),
			)
		}()

		next.ServeHTTP(rw, r)
	})
}

type responseWriter struct {
	http.ResponseWriter
	statusCode int
	size       int64
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func (rw *responseWriter) Write(b []byte) (int, error) {
	n, err := rw.ResponseWriter.Write(b)
	rw.size += int64(n)
	return n, err
}

func (rw *responseWriter) Unwrap() http.ResponseWriter {
	return rw.ResponseWriter
}
