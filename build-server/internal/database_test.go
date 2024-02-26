package database

import (
	"fmt"
	"testing"
)

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

}

func TestDB_getAllChirps(t *testing.T) {
	// Test the getAllChirps function
	// Create a new database
	databaseCon, err := NewDB("database.json")
	if err != nil {
		t.Errorf("Error creating new database: %v", err)
	}

	data, err := databaseCon.GetChirps()
	fmt.Printf("Data: %v\n", data)
	if err != nil {
		t.Errorf("Error getting chirps: %v", err)
	}
	if len(data) < 1 {
		t.Errorf("Chirps not found")
	}
}
