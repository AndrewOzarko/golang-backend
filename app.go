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
	Addr string
}

func New(Addr string) *App {
	return &App{
		Addr: Addr,
	}
}

func (app *App) Start() {
	// Run database migrations
	// database.Rollback()
	database.Migrate()
	database.Seed()
	// Workers
	r := chi.NewRouter()
	routes(r)
	log.Printf("Started server on %s", app.Addr)
	log.Fatalln(http.ListenAndServe(app.Addr, r))
}
