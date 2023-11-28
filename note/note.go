package note

import (
	"github.com/davidkornel/notepad-demo/view"
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
)

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
	var notes []RenderNote

	notes = append(notes, RenderNote{
		GroupID: uuid.NewV4().String(),
		NoteID:  uuid.NewV4().String(),
		Title:   "Note",
		Text:    "Blandit tempus porttitor aasfs. Integer posuere erat a ante venenatis.",
	})

	view.Render(c, gin.H{
		"title": "demo",
		"notes": notes,
	}, "note-all.html")

}
