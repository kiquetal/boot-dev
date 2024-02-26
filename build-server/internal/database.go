package database

import (
	"encoding/json"
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

func (db *DB) NewDB(path string) (*DB, error) {
	//check if file exists in path
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		//create file
		file, err := os.Create(path)
		if err != nil {
			return nil, err
		}
		defer file.Close()
		//create empty chirps
	}

	// We return the path to the database file.
	return &DB{
		path: path,
	}, nil
}

func (db *DB) CreateChirp(body string) (Chirp, error) {

	db.mu.Lock()
	defer db.mu.Unlock()
	//open file
	data, err := os.ReadFile(db.path)
	if err != nil {
		return Chirp{}, err
	}
	//parse json
	var databaseContent DBStructure
	json.Unmarshal(data, &databaseContent)
	if len(databaseContent.Chirps) == 0 {
		databaseContent.Chirps = make(map[int]Chirp)

		var chrip = Chirp{
			Body: body,
			Id:   1,
		}
		databaseContent.Chirps[1] = chrip
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
	newData, _ := json.Marshal(databaseContent)
	os.WriteFile(db.path, newData, 0644)
	return chrip, nil

}
