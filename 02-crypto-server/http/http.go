package http

import (
	"cryptoserver/crypto"
	"cryptoserver/clean/composure"
	"fmt"
	"net/http"
	//"context"
	"strings"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/golang-jwt/jwt/v5"
)

func authRoute(r chi.Router) {
	auth := composure.NewAuth()
	r.Route("/auth", func(r chi.Router) {
		r.Post("/register", auth.RegisterUser) // POST /auth/register
		r.Post("/login",    auth.LoginUser)	   // POST /auth/login
	})
}

func authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			http.Error(w, "You are not authorized.", http.StatusUnauthorized)
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		secret := []byte("super_secret_key_that_should_be_long_and_random") // fix this, need security
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
			if token.Method.Alg() != "HS256" {
				return nil, fmt.Errorf("algo %v, expected HS256\n", token.Header["alg"])
			}
			return secret, nil
		})

		if err != nil || !token.Valid {
			http.Error(w, "Invalid token.", http.StatusUnauthorized)
			return
		}

		//const tokenKey = "token"
		//ctx := context.WithValue(context.Background(), tokenKey, token)
		next.ServeHTTP(w, r/*.WithContext(ctx)*/)
	})
}

func cryptoRoute(r chi.Router) {
	r.Route("/crypto", func(r chi.Router) {
		r.Use(authMiddleware)
		api := crypto.NewAPI()
		r.Get("/",  api.ListCryptos) // GET  /crypto
//		r.Post("/", addCrypto)   // POST /crypto
//
		r.Route("/{symbol}", func(r chi.Router) {
			r.Get("/",        api.GetCrypto) 	// GET    /crypto/{symbol}
//			r.Put("/refresh", updateCrypto) // PUT /crypto/{symbol}/refresh
//			r.Get("/history", getHistory)   // GET /crypto/{symbol}/history
//			r.Get("/stats",   getStats)     // GET /crypto/{symbol}/stats
//			r.Delete("/",     deleteCrypto) // DELETE /crypto/{symbol}
		})
	})
}

func CreateAndRun() error {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.URLFormat)
	r.Use(middleware.RequestID)
	r.Use(render.SetContentType(render.ContentTypeJSON))

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Root of cryptoserver.")
	})

	r.Get("/panic", func(w http.ResponseWriter, r *http.Request) {
		panic("test")
	})

	authRoute(r)
	cryptoRoute(r)
	return http.ListenAndServe(":8080", r)
}
