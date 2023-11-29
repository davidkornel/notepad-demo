package main

import (
	"github.com/davidkornel/notepad-demo/note"
	"github.com/gin-gonic/gin"
	"github.com/go-logr/logr"
)

func main() {
	logger := logr.Logger{}
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

	err := router.Run()
	if err != nil {
		return
	} // listen and serve on 0.0.0.0:8080
	setupLog.Info("Server set up succesfully, serving on http://localhost:8080")
}
