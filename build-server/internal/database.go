package database

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
)

type DB struct {
	path string
	mu   sync.RWMutex
}

type Chirp struct {
	Body string `json:"body"`
	Id   int    `json:"id"`
}

type DBStructure struct {
	Chirps map[int]Chirp `json:"chirps"`
}

func NewDB(path string) (*DB, error) {

	// We return the path to the database file.

	db := &DB{
		path: path,
	}

	error := db.ensureDB()
	if error != nil {
		return nil, error

	}
	return db, nil
}

func (db *DB) CreateChirp(body string) (Chirp, error) {

	db.mu.Lock()
	defer db.mu.Unlock()
	//open file
	db.ensureDB()
	databaseContent, err := db.loadDB()
	if err != nil {
		return Chirp{}, err
	}
	if len(databaseContent.Chirps) == 0 {
		databaseContent.Chirps = make(map[int]Chirp)
		var chrip = Chirp{
			Body: body,
			Id:   1,
		}
		databaseContent.Chirps[1] = chrip
		err := db.writeDB(databaseContent)
		if err != nil {
			return Chirp{}, err
		}
		return chrip, nil
	}
	var lastId int
	for k := range databaseContent.Chirps {
		if k > lastId {
			lastId = k
		}
	}
	var chrip = Chirp{
		Body: body,
		Id:   lastId + 1,
	}
	databaseContent.Chirps[lastId+1] = chrip
	err = db.writeDB(databaseContent)
	if err != nil {
		return Chirp{}, err

	}
	return chrip, nil

}

func (db *DB) loadDB() (DBStructure, error) {
	//open file
	data, err := os.ReadFile(db.path)
	if err != nil {
		return DBStructure{}, err
	}
	if len(data) == 0 {
		return DBStructure{}, nil
	}

	var databaseContent DBStructure
	err = json.Unmarshal(data, &databaseContent)
	if err != nil {
		fmt.Printf("Error unmarshalling data: %v", err)
		return DBStructure{}, err
	}
	return databaseContent, nil
}

func (db *DB) writeDB(data DBStructure) error {
	//open file
	newData, _ := json.Marshal(data)
	err := os.WriteFile(db.path, newData, 0644)
	if err != nil {
		return err
	}
	return nil
}
func (db *DB) ensureDB() error {

	//open file
	_, err := os.Stat(db.path)
	if err != nil {
		if os.IsNotExist(err) {
			//create file
			_, err := os.Create(db.path)
			if err != nil {
				return err
			}
			fmt.Printf("Database created at: %s\n", db.path)

		}

	}
	fmt.Printf("Database found at: %s\n", db.path)
	return nil
}

func (db *DB) GetChirps() ([]Chirp, error) {
	db.mu.RLock()
	var chirps []Chirp
	defer db.mu.RUnlock()
	databaseContent, err := db.loadDB()
	if err != nil {
		return nil, err
	}
	if len(databaseContent.Chirps) == 0 {
		return nil, nil
	}
	for _, chirp := range databaseContent.Chirps {
		chirps = append(chirps, chirp)
	}

	return chirps, nil
}
