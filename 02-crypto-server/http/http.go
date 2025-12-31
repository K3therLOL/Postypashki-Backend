package http


import (
	"fmt"
	"net/http"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

func authRoute(r chi.Router) {
	r.Route("/auth", func(r chi.Router) {
		r.Post("/register", registerUser) // POST /auth/register
		r.Post("/login",    loginUser)	  // POST /auth/login
	})
}

func cryptoRoute(r chi.Router) {
	r.Route("/crypto", func(r chi.Router) {
		r.Get("/",  listCryptos) // GET  /crypto
		r.Post("/", addCrypto)   // POST /crypto

		r.Route("/{symbol}", func(r chi.Router) {
			r.Get("/",        getCrypto) 	// GET    /crypto/{symbol}
			r.Put("/refresh", updateCrypto) // PUT /crypto/{symbol}/refresh
			r.Get("/history", getHistory)   // GET /crypto/{symbol}/history
			r.Get("/stats",   getStats)     // GET /crypto/{symbol}/stats
			r.Delete("/",     deleteCrypto) // DELETE /crypto/{symbol}
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

	return http.ListenAndServe(":8080", r)
}
