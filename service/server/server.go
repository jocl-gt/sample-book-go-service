package server

import (
	"fmt"
	"log"
	"net/http"
	"readcommend/configuration"
	db "readcommend/database"
	"readcommend/handlers"
	"readcommend/services"

	"github.com/jmoiron/sqlx"
	"github.com/rs/cors"
)

type ServerInstance struct {
	Config *configuration.Configuration
	DB     *sqlx.DB
	Router *http.ServeMux
	Cors   *cors.Cors
	Server *http.Server
}

func NewServerInstance(configurationFilePath string) (*ServerInstance, error) {

	config := configuration.LoadConfiguration(configurationFilePath)
	database := db.InitDB(config)
	router := http.NewServeMux()

	corsHandler := cors.New(cors.Options{
		AllowedOrigins: config.CORSOrigins,
	})

	bookService := services.NewBookService(database)
	authorService := services.NewAuthorService(database)
	genreService := services.NewGenreService(database)
	sizeService := services.NewSizeService(database)
	erasService := services.NewErasService(database)

	bookHandler := handlers.NewBookHandler(bookService)
	authorHandler := handlers.NewAuthorHandler(authorService)
	genreHandler := handlers.NewGenreHandler(genreService)
	sizeHandler := handlers.NewSizeHandler(sizeService)
	erasHandler := handlers.NewErasHandler(erasService)

	registerRoutes(router, bookHandler, authorHandler, genreHandler, sizeHandler, erasHandler)

	app := &ServerInstance{
		Config: &config,
		DB:     database,
		Router: router,
		Cors:   corsHandler,
		Server: nil,
	}

	return app, nil
}

func registerRoutes(
	mux *http.ServeMux,
	bookHandler *handlers.BookHandler,
	authorHandler *handlers.AuthorHandler,
	genreHandler *handlers.GenreHandler,
	sizeHandler *handlers.SizeHandler,
	erasHandler *handlers.ErasHandler,
) {
	mux.HandleFunc("/api/v1/genres", func(w http.ResponseWriter, r *http.Request) {
		genreHandler.GetGenres(w, r)
	})
	mux.HandleFunc("/api/v1/authors", func(w http.ResponseWriter, r *http.Request) {
		authorHandler.GetAuthors(w, r)
	})
	mux.HandleFunc("/api/v1/books", func(w http.ResponseWriter, r *http.Request) {
		bookHandler.GetBooks(w, r)
	})
	mux.HandleFunc("/api/v1/sizes", func(w http.ResponseWriter, r *http.Request) {
		sizeHandler.GetSizes(w, r)
	})
	mux.HandleFunc("/api/v1/eras", func(w http.ResponseWriter, r *http.Request) {
		erasHandler.GetEras(w, r)
	})
}

func (app *ServerInstance) Run() error {
	addr := fmt.Sprintf(":%d", app.Config.ServerPort)
	log.Printf("Server is running on port %d", app.Config.ServerPort)
	return http.ListenAndServe(addr, app.Cors.Handler(app.Router))
}

func (app *ServerInstance) Start() {
	addr := fmt.Sprintf(":%d", app.Config.ServerPort)
	if app.Server == nil {
		app.Server = &http.Server{
			Addr:    addr,
			Handler: app.Cors.Handler(app.Router),
		}
		log.Printf("Server is running on port %d", app.Config.ServerPort)

		go func() {
			if err := app.Server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				log.Fatalf("Error starting server: %v", err)
			}
		}()
	} else {
		log.Println("Server is already running")
	}
}

func (app *ServerInstance) Close() {
	if err := app.Server.Close(); err != nil {
		log.Fatalf("Error closing server: %v", err)
	}
}
