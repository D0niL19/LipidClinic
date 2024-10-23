package main

import (
	"LipidClinic/internal/config"
	"LipidClinic/internal/http-server/handlers/auth"
	"LipidClinic/internal/http-server/handlers/patient"
	"LipidClinic/internal/http-server/handlers/relations"
	"LipidClinic/internal/lib/cors"
	"LipidClinic/internal/lib/logger/sl"
	"LipidClinic/internal/middleware"
	"LipidClinic/internal/storage/postgres"
	"github.com/gin-contrib/requestid"
	"github.com/gin-gonic/gin"
	"log/slog"
	"net/http"
	"os"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {

	cfg := config.MustLoad()
	corsSettings := cors.CorsSettings()

	log := setupLogger(cfg.Env)

	log.Info("Starting application", slog.String("Env", cfg.Env))
	log.Debug("Debug messages are enabled")

	storage, err := postgres.New(cfg)
	//rdb, err := red.New(cfg)
	if err != nil {
		log.Error("failed to init storage", sl.Err(err))
		panic(err)
	}

	router := gin.New()

	router.Use(gin.Recovery())
	router.Use(gin.Logger())
	router.Use(requestid.New())

	router.Use(func(c *gin.Context) {
		corsSettings.HandlerFunc(c.Writer, c.Request)
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusOK)
			return
		}
		c.Next()
	})

	router.POST("/auth/register", auth.Register(log, storage, cfg))
	router.POST("/auth/sign-in", auth.SignIn(log, storage, cfg))
	router.POST("/auth/refresh", auth.Refresh(log, cfg))

	router.POST("/auth/confirm-account/:email/:token", auth.ConfirmAccount(log, storage))
	router.POST("/auth/reset-password", auth.ForgetPassword(log, storage, cfg))
	router.POST("/auth/change-password/:token", auth.ChangePassword(log, storage, cfg))

	patients := router.Group("/patients")
	patients.Use(middleware.AuthMiddleware(cfg.Jwt.Secret, log))
	patients.POST("/", patient.Add(log, storage))
	patients.GET("/", patient.ByEmail(log, storage))
	patients.GET("/:id", patient.ById(log, storage))
	patients.GET("/all", patient.All(log, storage))
	patients.DELETE("/:id", patient.Delete(log, storage))

	relationships := router.Group("/relationships")
	relationships.Use(middleware.AuthMiddleware(cfg.Jwt.Secret, log))
	relationships.POST("/", relations.Add(log, storage))
	relationships.GET("/:id", relations.ById(log, storage))
	relationships.GET("/all/:id", relations.AllById(log, storage))
	relationships.DELETE("/", relations.Delete(log, storage))

	log.Info("starting server", slog.String("address", cfg.Address))

	server := &http.Server{
		Addr:         cfg.Address,
		Handler:      router,
		WriteTimeout: cfg.Timeout,
		ReadTimeout:  cfg.Timeout,
		IdleTimeout:  cfg.IdleTimeout,
	}

	if err = server.ListenAndServe(); err != nil {
		log.Error("failed to start server", sl.Err(err))
	}

	log.Info("server stopped")
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelDebug,
		}))

	case envDev:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelDebug,
		}))
	case envProd:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelInfo,
		}))
	}

	return log
}
