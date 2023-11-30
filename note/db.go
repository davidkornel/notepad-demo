package note

import (
	"context"
	"errors"
	"github.com/davidkornel/notepad-demo/config"
	uuid "github.com/satori/go.uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func (r *Routes) fetchAllNoteFromDB() []Note {
	log := r.logger.WithName("fetchAllNoteFromDB")

	var results []Note
	log.Info("Trying to fetch all notes from database")
	coll := r.dbClient.Database(config.DefaultDatabaseTableName).Collection(config.DefaultDatabaseNoteCollectionName)

	cursor, err := coll.Find(context.TODO(), bson.D{})
	if errors.Is(err, mongo.ErrNoDocuments) {
		log.Info("MongoDB table is empty, no document found")
		return nil
	}
	if err != nil {
		log.Error(err, "error happened while finding notes from DB")
	}
	if err = cursor.All(context.TODO(), &results); err != nil {
		panic(err)
	}
	return results
}

func (r *Routes) saveNoteIntoDB(note Note) error {
	log := r.logger.WithName("saveNoteIntoDB")

	log.Info("Trying to save note into database", "note", note)
	coll := r.dbClient.Database(config.DefaultDatabaseTableName).Collection(config.DefaultDatabaseNoteCollectionName)

	_, err := coll.InsertOne(context.TODO(), note)
	if err != nil {
		log.Error(err, "could not save")
		return err
	}
	log.Info("Successfully saved note into db", "note", note)
	return nil
}

func (r *Routes) editNoteInDB(note Note) error {
	log := r.logger.WithName("saveNoteIntoDB")

	log.Info("Trying to save note into database", "note", note)
	coll := r.dbClient.Database(config.DefaultDatabaseTableName).Collection(config.DefaultDatabaseNoteCollectionName)

	update := bson.D{{"$set", bson.D{
		{"text", note.Text},
		{"group", note.Group},
	}}}
	result, err := coll.UpdateOne(context.TODO(), bson.D{{"noteid", note.NoteID}}, update)
	if err != nil {
		return err
	}

	if result.MatchedCount == 1 {
		log.Info("matched and updated successfully an existing note in DB")
	}
	if result.UpsertedCount != 0 {
		log.Info("upserted a new note", "id", note.NoteID, "title", note.Title)
	}

	return nil
}

func (r *Routes) deleteNoteFromDB(id uuid.UUID) error {
	log := r.logger.WithName("deleteNoteFromDB")

	log.Info("Trying to delete note from database", "noteid", id.String())
	coll := r.dbClient.Database(config.DefaultDatabaseTableName).Collection(config.DefaultDatabaseNoteCollectionName)

	result, err := coll.DeleteOne(context.TODO(), bson.D{{"noteid", id}})
	if err != nil {
		return err
	}
	if result.DeletedCount == 1 {
		log.V(1).Info("Successfully deleted note")
	}
	return nil
}
