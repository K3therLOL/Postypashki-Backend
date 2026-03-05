package repository

import (
	"database/sql"
	"log"
	"task/clean/domain"

	"github.com/google/uuid"
)

type Database struct {
	db *sql.DB
}

func NewRepository(connString string) *Database {
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

func (repo *Database) Save(uuid uuid.UUID) error {
}

func (repo *Database) Get(uuid uuid.UUID) (*domain.TaskObject, error) {
}
