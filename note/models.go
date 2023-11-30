package note

import (
	uuid "github.com/satori/go.uuid"
)

type POSTData struct {
	Title string `json:"title"`
	Group string `json:"group"`
	Text  string `json:"text"`
}

type Note struct {
	NoteID uuid.UUID `bson:"noteid" json:"noteid"`
	Group  string    `bson:"group" json:"group"`
	Title  string    `bson:"title" json:"title"`
	Text   string    `bson:"text" json:"text"`
}

type RenderNote struct {
	Group  string
	NoteID string
	Title  string
	Text   string
}

func POSTData2Note(uuid uuid.UUID, data POSTData) Note {
	return Note{
		NoteID: uuid,
		Group:  data.Group,
		Title:  data.Title,
		Text:   data.Text,
	}
}

func note2Map(note Note) map[string]string {
	return map[string]string{
		"title":  note.Title,
		"group":  note.Group,
		"text":   note.Text,
		"noteid": note.NoteID.String(),
	}
}
