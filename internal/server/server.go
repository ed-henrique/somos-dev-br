package server

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/ed-henrique/somos-dev-br/assets"
	"github.com/ed-henrique/somos-dev-br/internal/db"
	"github.com/ed-henrique/somos-dev-br/internal/models"
	"github.com/ed-henrique/somos-dev-br/internal/server/handlers"
)

type Config struct {
	IsDev bool
	Addr  string

	// Ignored if IsDev is true
	DBURI string
}

type Server struct {
	isDev bool
	addr  string
	r     *http.ServeMux
	q     *models.Queries
}

func New(sc Config) *Server {
	var (
		dbURI string
		seed  bool
	)

	if sc.IsDev {
		dbURI = ":memory:"
		slog.SetLogLoggerLevel(slog.LevelDebug)
	} else {
		dbURI = sc.DBURI
	}

	seed = sc.IsDev
	dbConn, err := db.New(dbURI, seed)
	if err != nil {
		slog.Error("could not open DB", slog.String("err", err.Error()))
		os.Exit(1)
	}

	return &Server{
		isDev: sc.IsDev,
		addr:  sc.Addr,
		r:     http.NewServeMux(),
		q:     models.New(dbConn),
	}
}

func (s *Server) AddRoutes() {
	assetsHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if s.isDev {
			w.Header().Set("Cache-Control", "no-store")

		}

		var fs http.Handler
		if s.isDev {
			fs = http.FileServer(http.Dir("./assets"))
		} else {
			fs = http.FileServer(http.FS(assets.Assets))
		}

		fs.ServeHTTP(w, r)
	})

	s.r.Handle("GET /assets/", http.StripPrefix("/assets/", assetsHandler))

	s.r.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/posts/", http.StatusMovedPermanently)
	})

	s.r.Handle("GET /posts/", handlers.GetAllPosts(s.q))
	s.r.HandleFunc("GET /users/", func(w http.ResponseWriter, r *http.Request) {})
}

func (s *Server) Run() {
	slog.Debug("Starting to listen...", slog.String("addr", s.addr))
	if err := http.ListenAndServe(s.addr, s.r); err != nil {
		slog.Error("could not start server", slog.String("err", err.Error()))
		os.Exit(1)
	}
}
