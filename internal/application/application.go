package application

import (
	"context"
	"errors"
	"net"
	"net/http"
	"sync"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"

	"github.com/Vikot10/viarticles/internal/dto"
	"github.com/Vikot10/viarticles/internal/storage"
)

type Application struct {
	logger *zap.Logger
	store  *storage.Storage
}

type VkProvider interface {
	GetFaves() ([]*dto.Fave, error)
	SynchronizeFaves() error
}

func New(store *storage.Storage) *Application {
	app := &Application{
		store: store,
	}

	return app
}

func (app *Application) Run(ctx context.Context, cancel context.CancelFunc, wg *sync.WaitGroup, ln net.Listener) {
	defer wg.Done()

	r := chi.NewRouter()
	r.Use(middleware.Recoverer)
	r.Get("/", func(http.ResponseWriter, *http.Request) {}) // liveness probe

	//app.registerRoutes(r)

	server := http.Server{
		Handler: r,
	}

	wg.Add(1)
	go func(wg *sync.WaitGroup) {
		defer wg.Done()

		<-ctx.Done()

		//logger.Info("server shutdown")
		errShutdown := server.Shutdown(context.Background())
		if errShutdown != nil {
			//logger.Error("server shutdown error", zap.Error(errShutdown))
		}
	}(wg)

	//logger.Info("server start", zap.String("address", ln.Addr().String()))
	errServe := server.Serve(ln)
	if errServe != nil {
		if !errors.Is(errServe, http.ErrServerClosed) {
			//logger.Error("server serve error", zap.Error(errServe))
		}
		cancel()
	}
}
