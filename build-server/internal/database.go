package database

import (
	"encoding/json"
	"errors"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"os"
	"sync"
	"time"
)

type DB struct {
	path      string
	mu        sync.RWMutex
	pathToken string
	pathUser  string
}

type Chirp struct {
	AuthorId int    `json:"author_id"`
	Body     string `json:"body"`
	Id       int    `json:"id"`
}

type User struct {
	Email       string `json:"email"`
	Id          int    `json:"id"`
	Password    string `json:"password"`
	IsChirpyRed bool   `json:"is_chirpy_red"`
}

type UserDB struct {
	Users map[int]User `json:"users"`
}

type DBStructure struct {
	Chirps map[int]Chirp `json:"chirps"`
}

type RevokedTokenDB struct {
	Tokens map[string]int `json:"tokens"`
}

func NewDB(path string) (*DB, error) {

	// We return the path to the database file.

	db := &DB{
		path:      path,
		pathToken: "tokens.json",
		pathUser:  "users.json",
	}

	err := db.ensureDB(db.path)
	if err != nil {
		return nil, err

	}
	return db, nil
}

func (db *DB) CreatePolkaWebhook(subId int) (User, error) {

	//open file
	userDatabaseContent, err := db.loadUsersDB()
	if err != nil {
		return User{}, err
	}
	var updateUser User
	updateUser = userDatabaseContent.Users[subId]
	updateUser.IsChirpyRed = true
	userDatabaseContent.Users[subId] = updateUser
	err = db.writeDB(userDatabaseContent, db.pathUser)
	if err != nil {
		return User{}, err

	}
	return updateUser, nil
}
func (db *DB) CreateChirp(body string, subId int) (Chirp, error) {

	//open file
	err := db.ensureDB(db.path)
	if err != nil {
		return Chirp{}, err
	}
	databaseContent, err := db.loadDB()
	if err != nil {
		return Chirp{}, err
	}
	if len(databaseContent.Chirps) == 0 {
		databaseContent.Chirps = make(map[int]Chirp)
		var chirp = Chirp{
			Body:     body,
			Id:       1,
			AuthorId: subId,
		}
		databaseContent.Chirps[1] = chirp
		err := db.writeDB(databaseContent, db.path)
		if err != nil {
			return Chirp{}, err
		}
		return chirp, nil
	}
	var lastId int
	for k := range databaseContent.Chirps {
		if k > lastId {
			lastId = k
		}
	}
	var chirp = Chirp{
		Body:     body,
		Id:       lastId + 1,
		AuthorId: subId,
	}
	databaseContent.Chirps[lastId+1] = chirp
	err = db.writeDB(databaseContent, db.path)
	if err != nil {
		return Chirp{}, err

	}
	return chirp, nil

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

	err := db.ensureDB(db.pathUser)
	if err != nil {
		return UserDB{}, err
	}
	data, err := os.ReadFile(db.pathUser)

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

func (db *DB) writeDB(data interface{}, dbFile string) error {
	//open file
	db.mu.Lock()
	defer db.mu.Unlock()
	newData, _ := json.Marshal(data)
	err := os.WriteFile(dbFile, newData, 0644)
	if err != nil {
		return err
	}
	return nil
}

func (db *DB) ensureDB(nameFileDb string) error {

	//open file
	_, err := os.Stat(nameFileDb)
	if err != nil {
		if os.IsNotExist(err) {
			//create file
			_, err := os.Create(nameFileDb)
			if err != nil {
				return err
			}
			fmt.Printf("Database created at: %s\n", nameFileDb)

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
		return Chirp{}, errors.New("chirp not found")
	}
	return chirp, nil
}

func (db *DB) GetChirpsByUserId(userId int) ([]Chirp, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()
	databaseContent, err := db.loadDB()
	if err != nil {
		return nil, err
	}
	if len(databaseContent.Chirps) == 0 {
		return []Chirp{}, nil
	}
	var chirps []Chirp
	for _, chirp := range databaseContent.Chirps {
		if chirp.AuthorId == userId {
			chirps = append(chirps, chirp)
		}
	}
	return chirps, nil
}

func (db *DB) CreateUser(email, password string) (User, error) {

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
		err := db.writeDB(userDatabaseContent, db.pathUser)
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
	err = db.writeDB(userDatabaseContent, db.pathUser)
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
		return User{}, errors.New("user not found")
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
	err = db.writeDB(userDatabaseContent, db.path)
	if err != nil {
		return User{}, err
	}
	return user, nil
}

func (db *DB) SaveRevokedToken(refresh string) (bool, error) {
	err := db.ensureDB(db.pathToken)
	if err != nil {
		return false, err
	}
	//open file
	allData, err := db.GetAllRevokedTokens()
	if err != nil {
		return false, err

	}
	if len(allData.Tokens) == 0 {
		allData.Tokens = make(map[string]int)
	}

	allData.Tokens[refresh] = int(time.Now().Unix())
	err = db.writeDB(allData, db.pathToken)
	if err != nil {
		return false, err
	}

	return false, nil
}

func (db *DB) IsRevokedToken(refresh string) (bool, error) {
	err := db.ensureDB(db.pathToken)
	if err != nil {
		return false, err
	}
	//open file
	allData, err := db.GetAllRevokedTokens()
	if err != nil {
		return false, err
	}
	_, ok := allData.Tokens[refresh]
	if ok {
		return true, nil
	}
	return false, nil

}

func (db *DB) GetAllRevokedTokens() (RevokedTokenDB, error) {
	//open file
	err := db.ensureDB(db.pathToken)
	if err != nil {
		return RevokedTokenDB{}, err
	}
	db.mu.RLock()
	defer db.mu.RUnlock()
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

func (db *DB) DeleteChirp(id int) error {
	db.mu.Lock()
	databaseContent, err := db.loadDB()
	db.mu.Unlock()
	if err != nil {
		return err
	}
	if len(databaseContent.Chirps) == 0 {
		return errors.New("chirp not found")
	}
	_, ok := databaseContent.Chirps[id]
	if !ok {
		return errors.New("chirp not found")
	}
	delete(databaseContent.Chirps, id)
	err = db.writeDB(databaseContent, db.path)
	if err != nil {
		return err
	}
	return nil
}
