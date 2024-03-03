package database

import (
	"encoding/json"
	"errors"
	"fmt"
	"golang.org/x/crypto/bcrypt"
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

type User struct {
	Email    string `json:"email"`
	Id       int    `json:"id"`
	Password string `json:"password"`
}

type UserDB struct {
	Users map[int]User `json:"users"`
}

type DBStructure struct {
	Chirps map[int]Chirp `json:"chirps"`
}

type RevokedTokenDB struct {
	Tokens map[string]bool `json:"tokens"`
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

func (db *DB) loadUsersDB() (UserDB, error) {
	//open file

	db.ensureDB()
	data, err := os.ReadFile(db.path)

	if err != nil {
		return UserDB{}, err

	}
	if len(data) == 0 {
		return UserDB{}, nil
	}
	var databaseContent UserDB
	err = json.Unmarshal(data, &databaseContent)
	if err != nil {
		fmt.Printf("Error unmarshalling data: %v", err)
		return UserDB{}, err
	}

	return databaseContent, nil
}

func (db *DB) writeDB(data interface{}) error {
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

func (db *DB) GetChirp(id int) (Chirp, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()
	databaseContent, err := db.loadDB()
	if err != nil {
		return Chirp{}, err
	}
	if len(databaseContent.Chirps) == 0 {
		return Chirp{}, nil
	}
	chirp, ok := databaseContent.Chirps[id]
	if !ok {
		return Chirp{}, errors.New("Chirp not found")
	}
	return chirp, nil
}

func (db *DB) CreateUser(email, password string) (User, error) {
	db.mu.Lock()
	defer db.mu.Unlock()
	//open file
	userDatabaseContent, err := db.loadUsersDB()
	if err != nil {
		return User{}, err
	}
	//hash password
	hashedPassword, err := db.generateHashFromPassword(password)
	if err != nil {
		return User{}, err
	}
	if len(userDatabaseContent.Users) == 0 {
		userDatabaseContent.Users = make(map[int]User)
		var user = User{
			Email:    email,
			Id:       1,
			Password: string(hashedPassword),
		}
		userDatabaseContent.Users[1] = user
		err := db.writeDB(userDatabaseContent)
		if err != nil {
			return User{}, err
		}
		return user, nil
	}
	var lastId int
	for k := range userDatabaseContent.Users {
		if k > lastId {
			lastId = k
		}

	}
	var user = User{
		Email:    email,
		Id:       lastId + 1,
		Password: string(hashedPassword),
	}
	userDatabaseContent.Users[lastId+1] = user
	err = db.writeDB(userDatabaseContent)
	if err != nil {
		return User{}, err
	}
	return user, nil

}

func (db *DB) generateHashFromPassword(password string) ([]byte, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	fmt.Printf("Hashed password: %s\n", hashed)
	return hashed, nil
}

func (db *DB) Login(email string, password string) (User, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()
	userDatabaseContent, err := db.loadUsersDB()
	if err != nil {
		return User{}, err
	}
	if len(userDatabaseContent.Users) == 0 {
		return User{}, errors.New("User not found")
	}
	var foundUser User
	for _, user := range userDatabaseContent.Users {
		if user.Email == email {
			foundUser = user
			break
		}
	}
	if foundUser == (User{}) {
		return User{}, errors.New("user not found")

	}
	if len(foundUser.Password) > 0 {
		_, err := db.checkPassword(password, []byte(foundUser.Password))
		if err != nil {
			return User{}, errors.New("credentials not valid")
		}
	}
	return foundUser, nil
}

func (db *DB) checkPassword(password string, hashedPassword []byte) (bool, error) {
	err := bcrypt.CompareHashAndPassword(hashedPassword, []byte(password))
	if err != nil {
		return false, err
	}
	return true, nil
}

func (db *DB) GetUser(id int) (User, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()
	userDatabaseContent, err := db.loadUsersDB()
	if err != nil {
		return User{}, err
	}
	if len(userDatabaseContent.Users) == 0 {
		return User{}, errors.New("user not found")
	}
	user, ok := userDatabaseContent.Users[id]
	if !ok {
		return User{}, errors.New("user not found")
	}
	return user, nil
}

func (db *DB) UpdateUser(user User) (User, error) {
	db.mu.Lock()
	defer db.mu.Unlock()
	userDatabaseContent, err := db.loadUsersDB()
	if err != nil {
		return User{}, err
	}
	if len(userDatabaseContent.Users) == 0 {
		return User{}, errors.New("user not found")
	}
	_, ok := userDatabaseContent.Users[user.Id]
	if !ok {
		return User{}, errors.New("user not found")
	}
	var newPassword, e = db.generateHashFromPassword(user.Password)
	if e != nil {
		return User{}, e
	}

	userDatabaseContent.Users[user.Id] = User{
		Email:    user.Email,
		Password: string(newPassword),
		Id:       user.Id,
	}
	err = db.writeDB(userDatabaseContent)
	if err != nil {
		return User{}, err
	}
	return user, nil
}

func (db *DB) SaveRevokedToken(refresh string) (bool, error) {
	db.mu.Lock()
	db.ensureDB()
	defer db.mu.Unlock()
	//open file
	allData, error := db.GetAllRevokedTokens()
	if error != nil {
		return false, error

	}
	if len(allData.Tokens) == 0 {
		allData.Tokens = make(map[string]bool)
	}

	allData.Tokens[refresh] = true
	err := db.writeDB(allData)
	if err != nil {
		return false, err
	}

	return false, nil
}

func (db *DB) IsRevokedToken(refresh string) (bool, error) {
	defer db.mu.RUnlock()
	db.ensureDB()
	//open file
	allData, error := db.GetAllRevokedTokens()
	if error != nil {
		return false, error
	}
	db.mu.RLock()
	_, ok := allData.Tokens[refresh]
	if ok {
		return true, nil
	}
	return false, nil

}

func (db *DB) GetAllRevokedTokens() (RevokedTokenDB, error) {
	//open file
	db.ensureDB()

	data, err := os.ReadFile(db.path)
	if err != nil {
		return RevokedTokenDB{}, err
	}
	if len(data) == 0 {
		return RevokedTokenDB{}, nil
	}
	var revokedTokens RevokedTokenDB
	err = json.Unmarshal(data, &revokedTokens)
	if err != nil {
		fmt.Printf("Error unmarshalling data: %v", err)
		return RevokedTokenDB{}, err
	}

	return revokedTokens, nil
}
