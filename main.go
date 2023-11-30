package main

import (
	"flag"
	"github.com/davidkornel/notepad-demo/config"
	"github.com/davidkornel/notepad-demo/database"
	"github.com/davidkornel/notepad-demo/note"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap/zapcore"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
)

func main() {
	var mongoDBUri string
	// parse cli flags
	flag.StringVar(&mongoDBUri, "mongo-uri", config.DefaultMongoDBConnectionURI, "The mongodb uri to be used when connecting")
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
	setupLog.Info("mongodb", "uri", mongoDBUri)

	// connect to mongodb
	mongoDB := database.NewMongoDB(logger, mongoDBUri)
	err := mongoDB.Connect2MongoDB()
	if err != nil {
		setupLog.Error(err, "failed to connect to mongodb")
	}

	// initialize gin engine
	router := gin.Default()

	router.LoadHTMLGlob("templates/*.html")
	router.Static("/assets", "./assets")

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	noteRoute := note.NewRoutes(logger, mongoDB.GetClient())

	noteRoute.RegisterRoutes(router)

	setupLog.Info("Server set up successfully, serving on http://localhost:8080")
	err = router.Run()
	if err != nil {
		setupLog.Error(err, "could not start/run Gin engine router")
		return
	} // listen and serve on 0.0.0.0:8080

}
