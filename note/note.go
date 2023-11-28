package note

import "github.com/gin-gonic/gin"

func RegisterRoutes(router *gin.Engine) {

	noteRoutes := router.Group("/note")
	noteRoutes.Use() //Authentication middleware should go here
	{
		noteRoutes.GET("/all", showAllNotes)
		//noteRoutes.POST("/create", createNote)
		//noteRoutes.GET("/view/:id", viewNote)
		//noteRoutes.POST("/edit/:id", editNote)
		//noteRoutes.DELETE("/delete/:id", deleteNote)
	}
}

func showAllNotes(c *gin.Context) {
	var notes []Note

	notes = append(notes, demoNote)

}
