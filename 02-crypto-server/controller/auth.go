package controller

import (
	"net/http"
	"encoding/json"
	"cryptoserver/usecase"
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

func (controller *Auth) RegisterUser(w http.ResponseWriter, r *http.Request) {
	data := &userDTO{}
	if err := json.NewDecoder(r.Body).Decode(data); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	
	if err := controller.ua.Register(data.Username, data.Password); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

}

func LoginUser(w http.ResponseWriter, r *http.Request) {
}
