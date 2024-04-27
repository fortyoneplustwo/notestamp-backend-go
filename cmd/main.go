package main

import (
	"database/sql"
	"os"
	"log"
	"net/http"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"

	"notestamp/auth"
	"notestamp/project"
	"notestamp/user"

  "context"
  // "github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func main() {
  if err := godotenv.Load("../.env"); err != nil {
    log.Fatalf("Error loading .env file: %v", err)
  }

  // Initialize db service
  dbUser := os.Getenv("DB_USER")
  dbName := os.Getenv("DB_NAME")
  dbPassword := os.Getenv("DB_PW")
  connStr := "user=" + dbUser + " dbname=" + dbName +  " password=" + dbPassword
  db, err := sql.Open("postgres", connStr)
  if err != nil {
    panic(err)
  }

  // Initialize s3 service
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatal(err)
	}
	s3Client := s3.NewFromConfig(cfg)

  // Initialize stores
  userStore := user.NewUserDB(db)
  projectStore := project.NewProjectDB(db)
  revokedStore := auth.NewRevokedDB(db)
  mediaStore := project.NewMediaBucket("timestampdocsbucket", s3Client)
  notesStore := project.NewNotesBucket("timestampdocsbucket", s3Client)

  // Create router
  router := mux.NewRouter()

  // Create handlers
  home := HomeHandler{}
  auth := auth.NewAuthHandler(userStore, revokedStore)
  project := project.NewProjectHandler(projectStore, userStore, mediaStore, notesStore, revokedStore)

  // Register routes
  router.HandleFunc("/", home.ServeHTTP)

  router.HandleFunc("/auth", auth.ServeHTTP)
  router.HandleFunc("/auth/register", auth.Register).Methods("POST")
  router.HandleFunc("/auth/login", auth.Login).Methods("POST")
  router.HandleFunc("/auth/logout", auth.Logout).Methods("POST")
  router.HandleFunc("/auth/unregister", auth.Unregister).Methods("POST")

  router.HandleFunc("/project", project.ServeHTTP)
  router.HandleFunc("/project/save", project.Save).Methods("POST")
  router.HandleFunc("/project/get/{title}", project.Get).Methods("GET")
  router.HandleFunc("/project/list", project.List).Methods("GET")
  router.HandleFunc("/project/delete/{title}", project.Delete).Methods("DELETE")
  router.HandleFunc("/media/download/{title}", project.DownloadMedia).Methods("GET")
  router.HandleFunc("/media/stream/{title}", project.ServeHTTP).Methods("GET")

  // Serve
  http.ListenAndServe(":8000", router)
}


// Home Handler
type HomeHandler struct{}

func (h *HomeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte("This is my home page\n"))
}

