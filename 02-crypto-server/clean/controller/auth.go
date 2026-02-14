package controller

import (
	"cryptoserver/clean/usecase"
	"cryptoserver/errorfmt"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

var (
	ErrInvalidJson = errors.New("Invalid json.")
	ErrInvalidDTO  = errors.New("Username and password required.")
)

type userDTO struct { // DATA TRANSFER OBJECT
	Username string `json:"username"`
	Password string `json:"password"`
}

type Auth struct {
	ua *usecase.Auth
}

func NewAuth(ua *usecase.Auth) *Auth {
	return &Auth{ua: ua}
}

type tokenJson struct {
	Token string `json:"token"`
}

func newTokenJson(token string) tokenJson {
	return tokenJson{Token: token}
}

func formateToken(token string) string {
	tokenStruct := newTokenJson(token)
	tokenJson, _ := json.Marshal(tokenStruct)
	return string(tokenJson)
}

func (controller *Auth) RegisterUser(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	data := &userDTO{}
	if err := json.NewDecoder(r.Body).Decode(data); err != nil {
		http.Error(w, errorfmt.Jsonize(ErrInvalidJson), http.StatusBadRequest)
		return
	}

	if data.Username == "" || data.Password == "" {
		http.Error(w, errorfmt.Jsonize(ErrInvalidDTO), http.StatusBadRequest)
		return
	}

	if err := controller.ua.Register(data.Username, data.Password); err != nil {
		http.Error(w, errorfmt.Jsonize(err), http.StatusConflict)
		return
	}

	tokenString, err := createToken()
	if err != nil {
		http.Error(w, errorfmt.Jsonize(err), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
	fmt.Fprintln(w, formateToken(string(tokenString)))
}

func (controller *Auth) LoginUser(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	data := &userDTO{}
	if err := json.NewDecoder(r.Body).Decode(data); err != nil {
		http.Error(w, errorfmt.Jsonize(ErrInvalidJson), http.StatusBadRequest)
		return
	}

	if data.Username == "" || data.Password == "" {
		http.Error(w, errorfmt.Jsonize(ErrInvalidDTO), http.StatusBadRequest)
		return
	}

	err := controller.ua.Login(data.Username, data.Password)
	if err != nil {
		http.Error(w, errorfmt.Jsonize(err), http.StatusUnauthorized)
		return
	}

	tokenString, err := createToken()
	if err != nil {
		http.Error(w, errorfmt.Jsonize(err), http.StatusUnauthorized)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, formateToken(string(tokenString)))
}

func createToken() (string, error) {
	claims := &jwt.RegisteredClaims{
		Subject:   uuid.NewString(),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(30 * time.Minute)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	secret := []byte(os.Getenv("JWT_SECRET"))
	return token.SignedString(secret)
}
