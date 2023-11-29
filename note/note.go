package note

import (
	"github.com/davidkornel/notepad-demo/view"
	"github.com/gin-gonic/gin"
	"github.com/go-logr/logr"
	uuid "github.com/satori/go.uuid"
	"net/http"
)

type Routes struct {
	logger logr.Logger
}

func NewRoutes(l logr.Logger) *Routes {
	return &Routes{
		logger: l,
	}
}

func (r *Routes) RegisterRoutes(router *gin.Engine) {
	noteRoutes := router.Group("/note")
	noteRoutes.Use() //Authentication middleware should go here
	{
		noteRoutes.GET("/all", r.showAllNotes)
		noteRoutes.POST("/create", r.createNote)
		//noteRoutes.GET("/view/:id", viewNote)
		//noteRoutes.POST("/edit/:id", editNote)
		//noteRoutes.DELETE("/delete/:id", deleteNote)
	}
	r.logger.WithName("RegisterRoutes").V(1).Info("routes registered for /note")

}

func (r *Routes) createNote(c *gin.Context) {
	log := r.logger.WithName("createNote")

	var noteData POSTData

	if err := c.ShouldBindJSON(&noteData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	noteID := uuid.NewV4()
	groupID := uuid.NewV4()
	payload := map[string]string{
		"title":   noteData.Title,
		"group":   noteData.Group,
		"text":    noteData.Text,
		"noteid":  noteID.String(),
		"groupid": groupID.String(),
	}

	//saveNote

	view.Render(c, gin.H{
		"payload": payload,
	}, "")
	log.Info("new note received", "note", payload)
}

func (r *Routes) showAllNotes(c *gin.Context) {
	log := r.logger.WithName("showAllNotes")

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
	log.V(1).Info("notes have been fetched", "notes", notes)
}
