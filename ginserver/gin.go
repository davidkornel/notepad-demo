package ginserver

import (
	"context"
	"errors"
	"github.com/davidkornel/notepad-demo/database"
	"github.com/davidkornel/notepad-demo/monitoring"
	"github.com/davidkornel/notepad-demo/note"
	"github.com/gin-gonic/gin"
	"github.com/go-logr/logr"
	"net/http"
)

type Server struct {
	logger       logr.Logger
	metricServer *monitoring.MetricServer
	mongoDB      *database.MongoDB
	srv          *http.Server
}

func NewGinServer(logger logr.Logger, metricServer *monitoring.MetricServer, db *database.MongoDB) *Server {
	return &Server{
		logger:       logger,
		metricServer: metricServer,
		mongoDB:      db,
	}
}

func (s *Server) Start() {
	log := s.logger.WithName("MetricServer")
	log.Info("Starting MetricServer HTTP server")
	router := gin.Default()
	router.LoadHTMLGlob("templates/*.html")
	router.Static("/assets", "./assets")

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	noteRoute := note.NewRoutes(s.logger, s.metricServer.GetMetrics(), s.mongoDB.GetClient())

	noteRoute.RegisterRoutes(router)

	s.srv = &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	go func() {
		log.Info("Server set up successfully, serving on http://localhost:8080")
		err := s.srv.ListenAndServe()
		if errors.Is(err, http.ErrServerClosed) {
			log.Info("gin server: normal shutdown")
		} else {
			log.Error(err, "error running Gin HTTP server on", "addr", s.srv.Addr)
		}
	}()
}

func (s *Server) Shutdown(ctx context.Context) {
	if err := s.srv.Shutdown(ctx); err != nil {
		s.logger.WithName("GinServerShutdown").Error(err, "error happened while shutting down gin server")
	}
}
