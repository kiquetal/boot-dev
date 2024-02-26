package database

import "testing"

func TestDB_NewDB(t *testing.T) {

	// Test the NewDB function
	// Create a new database
	db := &DB{}
	databaseCon, err := db.NewDB("test.json")
	if err != nil {
		t.Errorf("Error creating new database: %v", err)
	}
	if databaseCon.path != "test.json" {
		t.Errorf("Path is not correct")
	}
	// create a chrip
	_, err = databaseCon.CreateChirp("Hello World")
	if err != nil {
		t.Errorf("Error creating new chirp: %v", err)
	}
	// read the file

}
