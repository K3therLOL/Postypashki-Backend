package controller

import (
	"fmt"
	"errors"
	"net/http"
	"encoding/json"
	"cryptoserver/clean/usecase"
	"github.com/golang-jwt/jwt/v5"
)

var (
	ErrInvalidJson = errors.New("Invalid json")
	ErrInvalidDTO  = errors.New("Username and password required")
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

type errorJson struct {
	Err string `json:"error"`
}

func newErrorJson(err string) errorJson {
	return errorJson{Err: err}
}

func formateError(err error) string {
	errStruct := newErrorJson(err.Error())
	errJson, _ := json.Marshal(errStruct)
	return string(errJson)
}

func (controller *Auth) RegisterUser(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	data := &userDTO{}
	if err := json.NewDecoder(r.Body).Decode(data); err != nil {
		http.Error(w, formateError(ErrInvalidJson), http.StatusBadRequest)
		return
	}

	if data.Username == "" || data.Password == "" {
		http.Error(w, formateError(ErrInvalidDTO), http.StatusBadRequest)
		return
	}
	
	if err := controller.ua.Register(data.Username, data.Password); err != nil {
		http.Error(w, formateError(err), http.StatusConflict)
		return
	}

	token := jwt.New(jwt.SigningMethodHS256)
	secret := []byte("super_secret_key_that_should_be_long_and_random") // fix this
	tokenString, err := token.SignedString(secret)
	if err != nil {
		http.Error(w, formateError(err), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
	fmt.Fprintln(w, formateToken(string(tokenString)))
}

func (controller *Auth) LoginUser(w http.ResponseWriter, r *http.Request) {
}
