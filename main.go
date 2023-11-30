package main

import (
	"context"
	"flag"
	"github.com/davidkornel/notepad-demo/config"
	"github.com/davidkornel/notepad-demo/database"
	"github.com/davidkornel/notepad-demo/ginserver"
	"github.com/davidkornel/notepad-demo/monitoring"
	"go.uber.org/zap/zapcore"
	"os"
	"os/signal"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
	"syscall"
	"time"
)

func main() {
	var mongoDBUri string
	var monitoringPort int
	// parse cli flags
	flag.StringVar(&mongoDBUri, "mongo-uri", config.DefaultMongoDBConnectionURI, "The mongodb uri to be used when connecting")
	flag.IntVar(&monitoringPort, "monitoring-port", config.DefaultMonitoringPort, "The port used to expose the metrics")
	flag.Parse()

	// setup logging
	opts := zap.Options{
		Development: true,
	}
	opts.BindFlags(flag.CommandLine)
	flag.Parse()
	logger := zap.New(zap.UseFlagOptions(&opts), func(o *zap.Options) {
		o.TimeEncoder = zapcore.RFC3339NanoTimeEncoder
	})

	setupLog := logger.WithName("setupLog")

	// gracful shutdown
	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, syscall.SIGINT, syscall.SIGTERM)

	// setup to mongodb connection
	setupLog.Info("Trying to setup mongodb connection", "uri", mongoDBUri)
	mongoDB := database.NewMongoDB(logger, mongoDBUri)
	err := mongoDB.Connect2MongoDB()
	if err != nil {
		setupLog.Error(err, "failed to connect to mongodb")
	}

	// setup monitoring
	setupLog.Info("Trying to setup Monitoring", "port", monitoringPort)
	metricServer := monitoring.NewMetricServer(logger, monitoringPort)
	setupLog.Info("Starting metrics server")
	metricServer.Start()

	// initialize gin engine
	setupLog.Info("Trying to setup Gin")
	ginServer := ginserver.NewGinServer(logger, metricServer, mongoDB)
	setupLog.Info("Starting Gin server")
	ginServer.Start()

	// Wait for a signal to shut down the server
	sig := <-signalCh
	logger.WithName("shutdownLog").Info("Received signal", "signal", sig)
	// Create a context with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	ginServer.Shutdown(ctx)
	metricServer.Shutdown(ctx)
	mongoDB.CloseMongoDBConnection(ctx)
}
