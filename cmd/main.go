package main

import (
	"database/sql"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
	_ "github.com/lib/pq"
	"context"

	"notestamp/auth"
	"notestamp/project"
	"notestamp/user"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type App struct {
  Router *mux.Router
  UserStore user.UserStore
  MetadataStore project.MetadataStore
  NotesStore project.NotesStore
  MediaStore project.MediaStore
  RevokedStore auth.RevokedStore
}

func (app *App) Initialize(db *sql.DB, s3Client *s3.Client, notesBucket string, mediaBucket string) {
  // Initialize stores
	app.UserStore = user.NewUserDB(db)
	app.MetadataStore = project.NewProjectDB(db)
	app.RevokedStore = auth.NewRevokedDB(db)
	app.MediaStore = project.NewMediaBucket(mediaBucket, s3Client)
	app.NotesStore = project.NewNotesBucket(notesBucket, s3Client)

  // Register routes
	app.Router = mux.NewRouter()

	auth := auth.NewAuthHandler(app.UserStore, app.RevokedStore)
	project := project.NewProjectHandler(
    app.MetadataStore, 
    app.UserStore, 
    app.MediaStore, 
    app.NotesStore, 
    app.RevokedStore,
  )

	app.Router.HandleFunc("/auth", auth.ServeHTTP)
	app.Router.HandleFunc("/auth/register", auth.Register).Methods("POST")
	app.Router.HandleFunc("/auth/login", auth.Login).Methods("POST")
	app.Router.HandleFunc("/auth/logout", auth.Logout).Methods("POST")
	app.Router.HandleFunc("/auth/unregister", auth.Unregister).Methods("POST")

	app.Router.HandleFunc("/project", project.ServeHTTP)
	app.Router.HandleFunc("/project/save", project.Save).Methods("POST")
	app.Router.HandleFunc("/project/get/{title}", project.Get).Methods("GET")
	app.Router.HandleFunc("/project/list", project.List).Methods("GET")
	app.Router.HandleFunc("/project/delete/{title}", project.Delete).Methods("DELETE")

	app.Router.HandleFunc("/media/download/{title}", project.DownloadMedia).Methods("GET")
	app.Router.HandleFunc("/media/stream/{title}", project.StreamMedia).Methods("GET")
}

func (app *App) Run(port string) {
	http.ListenAndServe(port, app.Router)
}

func main() {
  // Start psql service
	if err := godotenv.Load("../.env"); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
	dbUser := os.Getenv("DB_USER")
	dbName := os.Getenv("DB_NAME")
	dbPassword := os.Getenv("DB_PW")
	connStr := "user=" + dbUser + " dbname=" + dbName + " password=" + dbPassword
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}

  // Start s3 service
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatal(err)
	}
	s3Client := s3.NewFromConfig(cfg)
  notesBucket, mediaBucket := "timestampdocsbucket", "timestampdocsbucket"

  // Start app
  app := App{}
  app.Initialize(db, s3Client, mediaBucket, notesBucket)
  app.Run(":8000")
}

