package note

import (
	uuid "github.com/satori/go.uuid"
)

var (
	demoNote = Note{
		GroupID: uuid.UUID{},
		NoteID:  uuid.UUID{},
		Title:   "Note",
		Text:    "Text",
	}
)

type Note struct {
	GroupID uuid.UUID `bson:"groupid"`
	NoteID  uuid.UUID `bson:"noteid"`
	Title   string    `bson:"title"`
	Text    string    `bson:"text"`
}
