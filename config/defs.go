package config

var (
	// DefaultMongoDBConnectionURI default mongodb uri, can be changed trough cli flag
	DefaultMongoDBConnectionURI = "mongodb://localhost:27017"

	// DefaultDatabaseTableName is the default database table name
	DefaultDatabaseTableName = "notepad"

	// DefaultDatabaseNoteCollectionName is the default database collection name for Notes
	DefaultDatabaseNoteCollectionName = "notes"

	// DefaultMonitoringPort is the default port on which we should expose the collected metrics
	DefaultMonitoringPort = 2112
)
