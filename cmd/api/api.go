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
	"github.com/go-chi/cors"
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

	// CORS configuration
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300, // Maximum value not ignored by browsers
	}))

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

	r.Route("/api/categories", func(r chi.Router) {
		r.Get("/", app.getAllCategoryHandler)
		r.Post("/", app.createCategoryHandler)

		r.Route("/{categoryID}", func(r chi.Router) {
			r.Use(app.categoryContextMiddleware)
			r.Get("/", app.getCategoryHandler)

			r.Patch("/", app.updateCategoryHandler)
			r.Delete("/", app.deleteCategoryHandler)
		})
	})

	r.Route("/api/assets", func(r chi.Router) {
		r.Get("/", app.getAllAssetHandler)
		r.Post("/", app.createAssetHandler)

		r.Route("/{assetID}", func(r chi.Router) {
			r.Use(app.assetContextMiddleware)
			r.Get("/", app.getAssetHandler)

			r.Patch("/", app.updateAssetHandler)
			r.Delete("/", app.deleteAssetHandler)
		})
	})

	r.Route("/asset-assignments", func(r chi.Router) {
		r.Post("/", app.CreateAssetAssignmentHandler)
	})

	r.Route("/api/profile", func(r chi.Router) {
		r.Use(app.AuthTokenMiddleware)
		r.Patch("/", app.updateUserHandler)
	})

	// Public routes
	r.Route("/api/authentication", func(r chi.Router) {
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
