package repository

import (
	"database/sql"
	"log"
	"os"
	domain "taskserver/clean/domain/auth"
)

type SessionRepository struct {
	logger *log.Logger
	db     *sql.DB
}

func NewSessionRepository(connString string) *SessionRepository {
	sessionRepo := new(SessionRepository)
	sessionRepo.logger = log.New(os.Stdout, "postgres: ", log.Ldate|log.Ltime)

	db, err := sql.Open("pgx", connString)
	if err != nil {
		sessionRepo.logger.Fatal(err)
	}

	err = db.Ping()
	if err != nil {
		sessionRepo.logger.Fatal(err)
	}

	sessionRepo.db = db
	sessionRepo.logger.Println("Successful connection to PostgreSQL.")
	return sessionRepo
}

func (sessionRepo *SessionRepository) Delete(session *domain.Session) error {
	query := `DELETE FROM sessions WHERE session_id = $1;`
	_, err := sessionRepo.db.Exec(query, session.ID)

	if err == nil {
		sessionRepo.logger.Printf("Session (user_id -- %d, token -- %s, expires_at -- %T) deleted from db.\n",
			session.UserID,
			session.Token,
			session.ExpiresAt)
	}

	return err
}

func (sessionRepo *SessionRepository) Save(session *domain.Session) error {
	query := `INSERT INTO sessions (session_id, user_id, token, expires_at) VALUES ($1, $2, $3, $4);`
	_, err := sessionRepo.db.Exec(query, session.ID, session.UserID, session.Token, session.ExpiresAt)

	if err == nil {
		sessionRepo.logger.Printf("Session (user_id -- %d, token -- %s, expires_at -- %T) added to db.\n",
			session.UserID,
			session.Token,
			session.ExpiresAt)
	}

	return err
}

func (sessionRepo *SessionRepository) GetByToken(token string) (*domain.Session, bool) {
	query := `SELECT session_id, user_id, token, expires_at FROM sessions WHERE token = $1;`

	session := &domain.Session{}
	if err := sessionRepo.db.QueryRow(query, token).Scan(&session.ID, &session.UserID, &session.Token, &session.ExpiresAt); err != nil {
		return nil, false
	}

	sessionRepo.logger.Printf("Session (user_id -- %d, token -- %s, expires_at -- %T) fetched from db.\n",
		session.UserID,
		session.Token,
		session.ExpiresAt)

	return session, true
}

func (sessionRepo *SessionRepository) GetByUserID(userID int) (*domain.Session, bool) {
	query := `SELECT session_id, user_id, token, expires_at FROM sessions WHERE user_id = $1;`

	session := &domain.Session{}
	if err := sessionRepo.db.QueryRow(query, userID).Scan(&session.ID, &session.UserID, &session.Token, &session.ExpiresAt); err != nil {
		return nil, false
	}

	sessionRepo.logger.Printf("Session (user_id -- %d, token -- %s, expires_at -- %T) fetched from db.\n",
		session.UserID,
		session.Token,
		session.ExpiresAt)

	return session, true
}
