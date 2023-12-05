package note

import (
	"github.com/davidkornel/notepad-demo/config"
	"github.com/davidkornel/notepad-demo/monitoring"
	"github.com/davidkornel/notepad-demo/view"
	"github.com/gin-gonic/gin"
	"github.com/go-logr/logr"
	uuid "github.com/satori/go.uuid"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
)

type Routes struct {
	coll    *mongo.Collection
	metrics *monitoring.Metrics
	logger  logr.Logger
}

func NewRoutes(l logr.Logger, metrics *monitoring.Metrics, client *mongo.Client) *Routes {
	return &Routes{
		coll:    client.Database(config.DefaultDatabaseTableName).Collection(config.DefaultDatabaseNoteCollectionName),
		metrics: metrics,
		logger:  l,
	}
}

func (r *Routes) RegisterRoutes(router *gin.Engine) {
	noteRoutes := router.Group("/note")
	noteRoutes.Use() //Authentication middleware should go here
	{
		noteRoutes.GET("/all", r.showAllNotes)
		noteRoutes.POST("/create", r.createNote)
		//noteRoutes.GET("/view/:id", viewNote)
		noteRoutes.POST("/edit/:id", r.editNote)
		noteRoutes.DELETE("/delete/:id", r.deleteNote)
	}
	r.logger.WithName("RegisterRoutes").V(1).Info("routes registered for /note")

}

func (r *Routes) deleteNote(c *gin.Context) {
	log := r.logger.WithName("deleteNote")

	r.metrics.IncrementRequestTotal("/note/delete", c.Request.Method)

	id, err := uuid.FromString(c.Param("id"))
	if err != nil {
		log.Error(err, "error happened while fetching id")
	}
	log.V(1).Info("NoteID fetched from url", "noteid", id.String())
	log.V(1).Info("Trying to delete note", "noteid", id.String())

	if err = r.deleteNoteFromDB(id); err != nil {
		log.Error(err, "error happened while updating note in DB")
		return
	}
	log.V(1).Info("Successfully deleted note", "id", id.String())
	log.Info("Route served successfully", "route", "/note/delete/"+id.String())

}

func (r *Routes) editNote(c *gin.Context) {
	log := r.logger.WithName("editNote")

	r.metrics.IncrementRequestTotal("/note/edit", c.Request.Method)

	id, err := uuid.FromString(c.Param("id"))
	if err != nil {
		log.Error(err, "error happened while fetching id")
	}
	log.Info("NoteID fetched from url", "noteid", id.String())
	log.V(1).Info("Trying to update note", "noteid", id.String())
	var noteData POSTData

	if err := c.ShouldBindJSON(&noteData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	note := postData2Note(id, noteData)

	err = r.editNoteInDB(note)
	if err != nil {
		log.Error(err, "error happened while updating note in DB")
		return
	}

	payload := note2Map(note)

	view.Render(c, gin.H{
		"payload": payload,
	}, "")
	log.V(1).Info("Successfully updated note", "id", note.NoteID.String())
	log.Info("Route served successfully", "route", "/note/edit/"+note.NoteID.String())
}

func (r *Routes) createNote(c *gin.Context) {
	log := r.logger.WithName("createNote")

	r.metrics.IncrementRequestTotal(c.Request.URL.Path, c.Request.Method)

	var noteData POSTData

	if err := c.ShouldBindJSON(&noteData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	noteID := uuid.NewV4()

	newNote := postData2Note(noteID, noteData)
	//saveNote
	err := r.saveNoteIntoDB(newNote)
	if err != nil {
		return
	}

	payload := note2Map(newNote)

	view.Render(c, gin.H{
		"payload": payload,
	}, "")
	log.V(1).Info("New note received", "note", payload)
	log.Info("Route served successfully", "route", "/note/create")
}

func (r *Routes) showAllNotes(c *gin.Context) {
	log := r.logger.WithName("showAllNotes")

	r.metrics.IncrementRequestTotal("/note/all", c.Request.Method)

	var renderNotes []RenderNote

	for _, n := range r.fetchAllNoteFromDB() {
		renderNotes = append(renderNotes, RenderNote{
			Group:  n.Group,
			NoteID: n.NoteID.String(),
			Title:  n.Title,
			Text:   n.Text,
		})
	}

	view.Render(c, gin.H{
		"title": "demo",
		"notes": renderNotes,
	}, "note-all.html")
	log.V(1).Info("Notes have been fetched", "notes", renderNotes)
	log.Info("Route served successfully", "route", "/note/all")
}
