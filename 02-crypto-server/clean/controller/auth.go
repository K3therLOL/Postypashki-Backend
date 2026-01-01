package controller

import (
	"fmt"
	"net/http"
	"encoding/json"
	"cryptoserver/usecase"
	"github.com/golang-jwt/jwt/v5"
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

func NewTokenJson(token string) tokenJson {
	return tokenJson{ Token: token }
}

type errorJson struct {
	Err string
}

func NewErrorJson(err string) errorJson {
	return errorJson{ Err: err }
}

func (controller *Auth) RegisterUser(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	token := jwt.New(jwt.SigningMethodHS256)
	secret := []byte("super_secret_key_that_should_be_long_and_random")
	tokenString, err := token.SignedString(secret)
	if err != nil {
		errJson, _ := json.Marshal(NewErrorJson(err.Error()))
		http.Error(w, string(errJson), http.StatusBadRequest)
		return
	}

	tokenJson, err := json.Marshal(NewTokenJson(tokenString))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	data := &userDTO{}
	if err := json.NewDecoder(r.Body).Decode(data); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	
	if err := controller.ua.Register(data.Username, data.Password); err != nil {
		http.Error(w, err.Error(), http.StatusConflict)
		return
	}

	w.WriteHeader(http.StatusCreated)
	fmt.Fprintln(w, string(tokenJson))
}

func (controller *Auth) LoginUser(w http.ResponseWriter, r *http.Request) {
}
