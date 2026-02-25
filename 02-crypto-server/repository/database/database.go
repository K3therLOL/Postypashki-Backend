package repository

import (
	"crypto/rand"
	"cryptoserver/clean/domain"
	"database/sql"
	"encoding/binary"
	"log"
	"math"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type Database struct {
	db *sql.DB
}

func NewStorage(connString string) *Database {
	db, err := sql.Open("pgx", connString)
	if err != nil {
		log.Fatal(err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Successful connection to postgres.")
	return &Database{db: db}
}

func generateRandomId() int32 {
	var id uint64
	binary.Read(rand.Reader, binary.BigEndian, &id)

	return int32(id % uint64(math.MaxInt32))
}

func (s *Database) Save(user *domain.User) error {
	query := `INSERT INTO users (user_id, username, password_hash) VALUES ($1, $2, $3);`
	id := generateRandomId()
	_, err := s.db.Exec(query, id, user.Username, user.PasswordHash)
	log.Printf("user registered: %s --- %s\n", user.Username, user.PasswordHash)
	return err
}

func (s *Database) Exist(username string) *domain.User {
	query := `SELECT username, password_hash FROM users WHERE username = $1;`

	var passwordHash string
	if err := s.db.QueryRow(query, username).Scan(&username, &passwordHash); err != nil {
		return nil
	}

	log.Printf("user %s exists\n", username)
	return &domain.User{Username: username, PasswordHash: passwordHash}
}
