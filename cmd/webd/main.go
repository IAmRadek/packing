package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/IAmRadek/go-kit/envconfig"
	"github.com/IAmRadek/packing/internal/algorithms/dp"
	"github.com/IAmRadek/packing/internal/app/allocation"
	"github.com/IAmRadek/packing/internal/app/inventory"
	"github.com/IAmRadek/packing/internal/domain/pack"
	"github.com/IAmRadek/packing/internal/handlers"
	"github.com/IAmRadek/packing/internal/infra"
	"github.com/IAmRadek/packing/internal/templates"
	gorillaHandlers "github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
)

type Config struct {
	Addr                     string        `env:"ADDR" default:":8080"`
	ReadTimeout              time.Duration `env:"READ_TIMEOUT" default:"10s"`
	ReadHeaderTimeout        time.Duration `env:"READ_HEADER_TIMEOUT" default:"10s"`
	WriteTimeout             time.Duration `env:"WRITE_TIMEOUT" default:"10s"`
	IdleTimeout              time.Duration `env:"IDLE_TIMEOUT" default:"10s"`
	MaxHeaderBytes           int           `env:"MAX_HEADER_BYTES" default:"1024"`
	GracefulShutdownDuration time.Duration `env:"GRACEFUL_SHUTDOWN_DURATION" default:"5s"`
}

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	defer cancel()

	log := slog.Default()

	var cfg Config
	err := envconfig.Read(&cfg, os.LookupEnv)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "reading config: %v", err)
		return
	}

	log.Info("Config Loaded")

	router := mux.NewRouter()

	render, err := templates.NewTemplates()
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "creating templates: %v", err)
		return
	}

	dec := schema.NewDecoder()

	log.Info("Templates Loaded", "templates", render.Templates())

	dpAlgo := dp.Allocator{}

	memRepo := infra.NewMemoryRepo()

	inv := pack.NewInventory("tires", pack.Sizes{
		pack.Size{
			ID:       "S",
			Capacity: 23,
			Label:    "S",
		},
		pack.Size{
			ID:       "L",
			Capacity: 31,
			Label:    "L",
		},
		pack.Size{
			ID:       "XL",
			Capacity: 53,
			Label:    "XL",
		},
	})
	memRepo.Save(ctx, inv)

	allocSrv := allocation.NewService(memRepo, dpAlgo)
	allocHandler := handlers.NewAllocationHandler(allocSrv)

	invSrv := inventory.NewService(memRepo)
	invHandlers := handlers.NewInventoryHandler(invSrv, allocSrv, render, dec)

	idxHandler := handlers.NewIndexHandler(render)

	registerRoutes(router, idxHandler, allocHandler, invHandlers)
	log.Info("Routes Registered")

	loggedRouter := gorillaHandlers.CustomLoggingHandler(
		os.Stdout,
		router,
		func(writer io.Writer, params gorillaHandlers.LogFormatterParams) {
			slog.LogAttrs(context.Background(), slog.LevelInfo, "http_request",
				slog.String("method", params.Request.Method),
				slog.String("uri", params.URL.String()),
				slog.Int("status", params.StatusCode),
				slog.Int("bytes", params.Size),
			)
		},
	)

	httpSrv := &http.Server{
		Addr:              cfg.Addr,
		Handler:           loggedRouter,
		ReadTimeout:       cfg.ReadTimeout,
		ReadHeaderTimeout: cfg.ReadHeaderTimeout,
		WriteTimeout:      cfg.WriteTimeout,
		IdleTimeout:       cfg.IdleTimeout,
		MaxHeaderBytes:    cfg.MaxHeaderBytes,
		BaseContext: func(listener net.Listener) context.Context {
			return ctx
		},
	}

	log.Info("Server Started", "addr", cfg.Addr)

	go func() {
		if err := httpSrv.ListenAndServe(); err != nil {
			if !errors.Is(err, http.ErrServerClosed) {
				log.Error("Server Failed", "err", err)
			}
		}
	}()

	<-ctx.Done()

	teardownCtx, cancel := context.WithTimeout(context.Background(), cfg.GracefulShutdownDuration)
	defer cancel()

	if err := httpSrv.Shutdown(teardownCtx); err != nil {
		log.Error("Server Failed to Shutdown", "err", err)
	}

	log.Info("Server Stopped")
}

func registerRoutes(
	router *mux.Router,
	handler *handlers.IndexHandler,
	allocHandler *handlers.AllocationHandler,
	invHandlers *handlers.InventoryHandler,
) {
	routes := []struct {
		path    string
		methods []string
		h       http.HandlerFunc
		// TODO: allow adding middlewares.
	}{
		{
			path:    "/api/allocate",
			methods: []string{"GET"},
			h:       allocHandler.HandleAllocate,
		},

		{
			path:    "/inventory/create",
			methods: []string{"GET", "POST"},
			h:       invHandlers.HandleCreate,
		},
		{
			path:    "/inventory/{sku}/update",
			methods: []string{"POST"},
			h:       invHandlers.HandleUpdate,
		},
		{
			path:    "/inventory/{sku}/delete",
			methods: []string{"POST"},
			h:       invHandlers.HandleDelete,
		},
		{
			path:    "/inventory/{sku}",
			methods: []string{"GET", "POST"},
			h:       invHandlers.HandleGet,
		},
		{
			path:    "/inventory",
			methods: []string{"GET"},
			h:       invHandlers.HandleList,
		},
		{
			path:    "/",
			methods: []string{"GET"},
			h:       handler.HandleIndex,
		},
	}

	for _, r := range routes {
		router.PathPrefix(r.path).Methods(r.methods...).Handler(r.h)
	}
}
