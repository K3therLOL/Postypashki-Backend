package repository

import (
	"cryptoserver/clean/domain"
	"database/sql"
	"log"

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

func (s *Database) Save(user *domain.User) error {
	return nil
}

func (s *Database) Exist(username string) *domain.User {
	passwordHash := ""
	//hash, ok := r.s[username]
	//if !ok {
	//	return nil
	//}
	return &domain.User{Username: username, PasswordHash: passwordHash}
}
