package repository

import (
	"database/sql"
	"errors"
	"log"
	"os"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type ImageRepository struct {
	logger *log.Logger
	db     *sql.DB
}

var (
	ErrLinkExpired = errors.New("Image link has expired.")
)

func NewImageRepository(connString string) *ImageRepository {
	imageRepo := new(ImageRepository)
	imageRepo.logger = log.New(os.Stdout, "image postgres: ", log.Ldate|log.Ltime)

	db, err := sql.Open("pgx", connString)
	if err != nil {
		imageRepo.logger.Fatal(err)
	}

	err = db.Ping()
	if err != nil {
		imageRepo.logger.Fatal(err)
	}

	imageRepo.db = db
	imageRepo.logger.Println("Successful connection to PostgreSQL.")
	return imageRepo
}

func (imageRepo *ImageRepository) delete(uuid string) error {
	query := `DELETE FROM images WHERE uuid = $1;`
	_, err := imageRepo.db.Exec(query, uuid)

	if err == nil {
		imageRepo.logger.Printf("Image with task uuid %s deleted from db.\n", uuid)
	}

	return err
}

func (imageRepo *ImageRepository) Get(uuid string) (string, error) {
	query := `SELECT url, expires_at FROM images WHERE uuid = $1;`

	var (
		imageUrl  string
		expiresAt time.Time
	)

	if err := imageRepo.db.QueryRow(query, uuid).Scan(&imageUrl, &expiresAt); err != nil {
		return "", err
	}

	if expiresAt.After(time.Now()) {
		imageRepo.delete(uuid)
		return "", ErrLinkExpired
	}

	imageRepo.logger.Printf("Image (image_url -- %s, expires_at -- %v) fetched from db.\n",
		imageUrl,
		expiresAt,
	)

	return imageUrl, nil
}

func (imageRepo *ImageRepository) Set(uuid, imageUrl string) error {
	return nil
}
