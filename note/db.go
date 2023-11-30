package note

import (
	"context"
	"github.com/davidkornel/notepad-demo/config"
)

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
