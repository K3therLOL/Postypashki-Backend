package provider

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/binary"
	"math"
	domain "taskserver/clean/domain/auth"
	"time"
)

type SessionProvider struct {
	ttlHours int
}

func NewSessionProvider(ttlHours int) *SessionProvider {
	return &SessionProvider{ttlHours: ttlHours}
}

func generateRandomId() int32 {
	var id uint64
	binary.Read(rand.Reader, binary.BigEndian, &id)

	return int32(id % uint64(math.MaxInt32))
}

func generateToken() (string, error) {
	b := make([]byte, 32)

	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}

	return base64.URLEncoding.EncodeToString(b), nil
}

func (sessionProvider *SessionProvider) Create() (*domain.Session, error) {
	token, err := generateToken()
	if err != nil {
		return nil, err
	}

	session := &domain.Session{
		ID:        int(generateRandomId()),
		UserID:    int(generateRandomId()),
		Token:     token,
		ExpiresAt: time.Now().UTC().Add(time.Duration(sessionProvider.ttlHours) * time.Hour),
	}

	return session, nil
}
