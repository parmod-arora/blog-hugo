package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	v1 "stripe-integration/cmd/serverd/v1"
	"stripe-integration/internal/db"
	"stripe-integration/pkg/envar"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

func main() {
	log.Printf("Start:")
	envar.Validate()
	r := chi.NewRouter()
	// A good base middleware stack
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	db.InitDatabase(os.Getenv("DATABASE_URL"))

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("welcome"))
	})

	// API routes
	r.Group(v1.Router())

	// File server
	workDir, err := os.Getwd()
	if err != nil {
		log.Fatalf("error is %v", err.Error())
	}
	filesDir := http.Dir(filepath.Join(workDir, "public"))
	FileServer(r, "/", filesDir)

	log.Printf("Starting server at :3000")
	http.ListenAndServe(":3000", r)
}

// FileServer conveniently sets up a http.FileServer handler to serve
// static files from a http.FileSystem.
func FileServer(r chi.Router, path string, root http.FileSystem) {
	if strings.ContainsAny(path, "{}*") {
		panic("FileServer does not permit any URL parameters.")
	}

	if path != "/" && path[len(path)-1] != '/' {
		r.Get(path, http.RedirectHandler(path+"/", 301).ServeHTTP)
		path += "/"
	}
	path += "*"

	r.Get(path, func(w http.ResponseWriter, r *http.Request) {
		rctx := chi.RouteContext(r.Context())
		pathPrefix := strings.TrimSuffix(rctx.RoutePattern(), "/*")
		fs := http.StripPrefix(pathPrefix, http.FileServer(root))
		fs.ServeHTTP(w, r)
	})
}
