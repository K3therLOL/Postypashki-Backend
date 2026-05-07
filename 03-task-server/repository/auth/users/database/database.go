package repository

import (
	"crypto/rand"
	"database/sql"
	"encoding/binary"
	"log"
	"math"
	"os"
	domain "taskserver/clean/domain/auth"
)

type UserRepository struct {
	logger *log.Logger
	db     *sql.DB
}

func NewUserRepository(connString string) *UserRepository {
	userRepo := new(UserRepository)
	userRepo.logger = log.New(os.Stdout, "postgres: ", log.Ldate|log.Ltime)

	db, err := sql.Open("pgx", connString)
	if err != nil {
		userRepo.logger.Fatal(err)
	}

	err = db.Ping()
	if err != nil {
		userRepo.logger.Fatal(err)
	}

	userRepo.db = db
	userRepo.logger.Println("Successful connection to PostgreSQL.")
	return userRepo
}

func generateRandomId() int32 {
	var id uint64
	binary.Read(rand.Reader, binary.BigEndian, &id)

	return int32(id % uint64(math.MaxInt32))
}

func (userRepo *UserRepository) Save(user *domain.User) error {
	query := `INSERT INTO users (user_id, username, password_hash) VALUES ($1, $2, $3);`
	_, err := userRepo.db.Exec(query, user.ID, user.Username, user.PasswordHash)

	userRepo.logger.Printf("User (user_id -- %d, username -- %s) added to db.\n",
		user.ID,
		user.Username)

	return err
}

func (userRepo *UserRepository) Exist(username string) (*domain.User, bool) {
	query := `SELECT user_id, user_id, token, expires_at FROM users WHERE username = $1;`

	user := &domain.User{}
	if err := userRepo.db.QueryRow(query, username).Scan(&user.ID, &user.ID, &user.Username, &user.PasswordHash); err != nil {
		return nil, false
	}

	userRepo.logger.Printf("User (user_id -- %d, username -- %s) fetched from db.\n",
		user.ID,
		user.Username)

	return user, true
}
