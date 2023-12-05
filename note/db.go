package note

import (
	"context"
	"errors"
	uuid "github.com/satori/go.uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (r *Routes) fetchAllNoteFromDB() []Note {
	defer r.updateMetric(r.coll)
	log := r.logger.WithName("fetchAllNoteFromDB")

	var results []Note
	log.V(1).Info("Trying to fetch all notes from database")

	cursor, err := r.coll.Find(context.TODO(), bson.D{})
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
	defer r.updateMetric(r.coll)
	log := r.logger.WithName("saveNoteIntoDB")

	log.V(1).Info("Trying to save note into database", "note", note)

	_, err := r.coll.InsertOne(context.TODO(), note)
	if err != nil {
		log.Error(err, "could not save")
		return err
	}
	log.Info("Successfully saved note into db", "note", note)
	return nil
}

func (r *Routes) editNoteInDB(note Note) error {
	defer r.updateMetric(r.coll)
	log := r.logger.WithName("saveNoteIntoDB")

	log.V(1).Info("Trying to save note into database", "note", note)

	update := bson.D{{"$set", bson.D{
		{"text", note.Text},
		{"group", note.Group},
	}}}
	result, err := r.coll.UpdateOne(context.TODO(), bson.D{{"noteid", note.NoteID}}, update)
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
	defer r.updateMetric(r.coll)
	log := r.logger.WithName("deleteNoteFromDB")
	log.V(1).Info("Trying to delete note from database", "noteid", id.String())

	result, err := r.coll.DeleteOne(context.TODO(), bson.D{{"noteid", id}})
	if err != nil {
		return err
	}
	if result.DeletedCount == 1 {
		log.V(1).Info("Successfully deleted note")
	}
	//r.updateMetric(coll)
	return nil
}

func (r *Routes) updateMetric(coll *mongo.Collection) {
	log := r.logger.WithName("updateMetric")

	log.V(1).Info("Trying to get the number of notes")

	opts := options.Count().SetHint("_id_")
	count, err := coll.CountDocuments(context.TODO(), bson.D{}, opts)
	if err != nil {
		panic(err)
	}
	log.V(1).Info("Number of notes", "num", count)
	r.metrics.SetNumberOfActiveNotes(count)
}
