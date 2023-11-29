package main

import (
	"flag"
	"github.com/davidkornel/notepad-demo/note"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap/zapcore"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
)

func main() {
	opts := zap.Options{
		Development: true,
	}
	opts.BindFlags(flag.CommandLine)
	flag.Parse()

	logger := zap.New(zap.UseFlagOptions(&opts), func(o *zap.Options) {
		o.TimeEncoder = zapcore.RFC3339NanoTimeEncoder
	})
	setupLog := logger.WithName("setupLog")
	router := gin.Default()

	router.LoadHTMLGlob("templates/*.html")
	router.Static("/assets", "./assets")

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	noteRoute := note.NewRoutes(logger)

	noteRoute.RegisterRoutes(router)

	setupLog.Info("Server set up succesfully, serving on http://localhost:8080")
	err := router.Run()
	if err != nil {
		return
	} // listen and serve on 0.0.0.0:8080

}
