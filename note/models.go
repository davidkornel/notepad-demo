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
	GroupID uuid.UUID `bson:"groupid"`
	NoteID  uuid.UUID `bson:"noteid"`
	Title   string    `bson:"title"`
	Text    string    `bson:"text"`
}

type RenderNote struct {
	GroupID string
	NoteID  string
	Title   string
	Text    string
}
