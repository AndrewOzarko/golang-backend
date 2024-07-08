package golangbackend

import (
	"golang-backend/internal/database"
	"golang-backend/internal/workers"
	"log"
	"net/http"

	"github.com/go-chi/chi"
)

func init() {
	go workers.JwtTokenWorker()

}

type App struct {
	Addr     string
	Migrate  bool
	Seed     bool
	Rollback bool
}

func New(Addr string, Migrate bool, Seed bool, Rollback bool) *App {
	return &App{
		Addr:     Addr,
		Migrate:  Migrate,
		Seed:     Seed,
		Rollback: Rollback,
	}
}

func (app *App) Start() {
	// Run database migrations
	if app.Migrate {
		database.Rollback()
		database.Migrate()
	}
	if app.Seed {
		database.Seed()
	}
	// Workers
	r := chi.NewRouter()
	routes(r)
	log.Printf("Started server on %s", app.Addr)
	log.Fatalln(http.ListenAndServe(app.Addr, r))
}
