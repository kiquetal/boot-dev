package database

import "testing"

func TestDB_NewDB(t *testing.T) {

	databaseCon, err := NewDB("test.json")
	if err != nil {
		t.Errorf("Error creating new database: %v", err)
	}
	if databaseCon.path != "test.json" {
		t.Errorf("Path is not correct")
	}

}

func TestDB_CreateChirp(t *testing.T) {
	// Test the CreateChirp function
	// Create a new database
	databaseCon, err := NewDB("test.json")
	if err != nil {
		t.Errorf("Error creating new database: %v", err)
	}
	// create a chrip
	_, err = databaseCon.CreateChirp("Hello World")
	if err != nil {
		t.Errorf("Error creating new chirp: %v", err)
	}
	databaseCon.CreateChirp("A second chirp")
	// read the file
}
