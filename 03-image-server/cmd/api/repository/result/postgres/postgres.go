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

func (imageRepo *ImageRepository) delete(task_id string) error {
	query := `DELETE FROM images WHERE task_id = $1;`
	_, err := imageRepo.db.Exec(query, task_id)

	if err == nil {
		imageRepo.logger.Printf("Image with task task_id %s deleted from db.\n", task_id)
	}

	return err
}

func (imageRepo *ImageRepository) Get(task_id string) (string, error) {
	query := `SELECT url, expires_at FROM images WHERE task_id = $1;`

	var (
		imageUrl  string
		expiresAt time.Time
	)

	if err := imageRepo.db.QueryRow(query, task_id).Scan(&imageUrl, &expiresAt); err != nil {
		imageRepo.logger.Println("Something happened with postgres.")
		return "", err
	}

	//if expiresAt.After(time.Now()) {
	//	imageRepo.logger.Printf("Link %s expired\n", imageUrl)
	//	imageRepo.delete(task_id)
	//	return "", ErrLinkExpired
	//}

	imageRepo.logger.Printf("Image (image_url -- %s, expires_at -- %v) fetched from db.\n",
		imageUrl,
		expiresAt,
	)

	return imageUrl, nil
}

func (imageRepo *ImageRepository) Set(task_id, imageUrl string) error {
	return nil
}
