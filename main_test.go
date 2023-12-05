package main

import (
	"github.com/davidkornel/notepad-demo/database"
	"github.com/davidkornel/notepad-demo/ginserver"
	"github.com/davidkornel/notepad-demo/monitoring"
	"net/http"
	"net/http/httptest"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNoteAllRoute(t *testing.T) {
	logger := zap.New()

	mongoDB := database.NewMongoDB(logger, "mongodb://localhost:27017")
	err := mongoDB.Connect2MongoDB()
	assert.Equal(t, err, nil)

	// setup monitoring
	metricServer := monitoring.NewMetricServer(logger, 2112)
	metricServer.Start()

	router := ginserver.NewGinServer(logger, metricServer, mongoDB)
	router.Start()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/note/all", nil)
	router.Router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
}
