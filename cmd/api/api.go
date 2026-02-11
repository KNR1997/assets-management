package main

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/knr1997/assets-management-apiserver/internal/auth"
	"github.com/knr1997/assets-management-apiserver/internal/store"
	httpSwagger "github.com/swaggo/http-swagger"
	"go.uber.org/zap"
)

type application struct {
	config        config
	store         store.Storage
	logger        *zap.SugaredLogger
	authenticator auth.Authenticator
}

type config struct {
	addr        string
	db          dbConfig
	env         string
	apiURL      string
	frontendURL string
	auth        authConfig
}

type authConfig struct {
	basic basicConfig
	token tokenConfig
}

type basicConfig struct {
	user string
	pass string
}

type tokenConfig struct {
	secret string
	exp    time.Duration
	iss    string
}

type dbConfig struct {
	addr         string
	maxOpenConns int
	maxIdleConns int
	maxIdleTime  string
}

func (app *application) mount() http.Handler {
	r := chi.NewRouter()

	// A good base middleware stack
	r.Use(middleware.RequestID) // important for rate limiting
	r.Use(middleware.RealIP)    // import for rate limiting and analytics and tracing
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer) // recover from crashes
	r.Use(middleware.Timeout(60 * time.Second))

	// Swagger UI
	r.Get("/swagger/*", httpSwagger.WrapHandler)

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("all good"))
	})

	r.Route("/categories", func(r chi.Router) {
		r.Post("/", app.createCategoryHandler)

		r.Route("/{categoryID}", func(r chi.Router) {
			r.Use(app.categoryContextMiddleware)
			r.Get("/", app.getCategoryHandler)

			r.Patch("/", app.updateCategoryHandler)
			r.Delete("/", app.deleteCategoryHandler)
		})
	})

	r.Route("/assets", func(r chi.Router) {
		r.Post("/", app.createAssetHandler)

		r.Route("/{assetID}", func(r chi.Router) {
			r.Use(app.assetContextMiddleware)
			r.Get("/", app.getAssetHandler)

			r.Patch("/", app.updateAssetHandler)
			r.Delete("/", app.deleteAssetHandler)
		})
	})

	// Public routes
	r.Route("/authentication", func(r chi.Router) {
		r.Post("/user", app.registerUserHandler)
		r.Post("/token", app.createTokenHandler)
	})

	return r
}

func (app *application) run(mux http.Handler) error {
	srv := http.Server{
		Addr:         app.config.addr,
		Handler:      mux,
		WriteTimeout: time.Second * 30,
		ReadTimeout:  time.Second * 10,
		IdleTimeout:  time.Minute,
	}

	shutdown := make(chan error)

	go func() {
		quit := make(chan os.Signal, 1)

		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		s := <-quit

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		app.logger.Infow("signal caught", "signal", s.String())

		shutdown <- srv.Shutdown(ctx)
	}()

	app.logger.Infow("server has started", "addr", app.config.addr, "env", app.config.env)

	err := srv.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	err = <-shutdown
	if err != nil {
		return err
	}

	app.logger.Infow("server has stopped", "addr", app.config.addr, "env", app.config.env)

	return nil
}
