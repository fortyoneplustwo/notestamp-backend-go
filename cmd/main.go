package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"notestamp/auth"
	"notestamp/handler"
	"notestamp/media"
	"notestamp/metadata"
	staging "notestamp/metadata/staging_area"
	"notestamp/middleware"
	"notestamp/notes"
	"notestamp/user"
	"os"
	"time"

	"cloud.google.com/go/storage"
	"github.com/go-redis/redis"
	"github.com/joho/godotenv"

	firebase "firebase.google.com/go"
	"google.golang.org/api/option"
)

func main() {
	if err := godotenv.Load("../.env"); err != nil {
		log.Print("No .env file found")
	}

	redisClient := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_ADDR"),
		Password: "",
		DB:       0,
	})

	ctx := context.TODO()
	sa := option.WithCredentialsFile(os.Getenv("FIREBASE_CONF"))
	firebase, err := firebase.NewApp(ctx, nil, sa)
	if err != nil {
		log.Fatalln(err)
	}

	firestoreClient, err := firebase.Firestore(ctx)
	if err != nil {
		log.Fatalln(err)
	}
	defer func() {
		err := firestoreClient.Close()
		if err != nil {
			log.Println(err)
		}
	}()

	storageClient, err := storage.NewClient(
		ctx,
		option.WithCredentialsFile(os.Getenv("FIREBASE_CONF")),
	)
	if err != nil {
		log.Fatalln(err)
	}
	defer func() {
		err := storageClient.Close()
		if err != nil {
			log.Println(err)
		}
	}()

	userStore := user.NewUserStore(firestoreClient)
	revokedStore := auth.NewRevokedTokenStore(firestoreClient)
	stagingArea := staging.NewStagingArea(redisClient, time.Hour)
	metadataStore := metadata.NewMetadataStore(firestoreClient)
	mediaStore := media.NewMediaStore(storageClient, os.Getenv("NOTES_BUCKET"))
	notesStore := notes.NewNotesStore(storageClient, os.Getenv("MEDIA_BUCKET"))

	mux := http.NewServeMux()
	mux.HandleFunc("POST /register", handler.RegisterUser(userStore))
	mux.HandleFunc("POST /login", handler.LoginUser(userStore, revokedStore))
	mux.HandleFunc("DELETE /logout", handler.LogoutUser(revokedStore))
	mux.HandleFunc(
		"DELETE /deregister",
		handler.DeregisterUser(
			userStore,
			metadataStore,
			mediaStore,
			notesStore,
			revokedStore,
		),
	)

	protectedMux := http.NewServeMux()
	protectedMux.HandleFunc("GET /list", handler.ListMetadata(metadataStore))
	protectedMux.HandleFunc(
		"DELETE /remove",
		handler.DeleteProject(metadataStore, mediaStore, notesStore),
	)
	protectedMux.HandleFunc(
		"POST /save-without-media",
		handler.StageWithoutMedia(stagingArea, notesStore),
	)
	protectedMux.HandleFunc(
		"POST /save-with-media",
		handler.StageWithMedia(stagingArea, mediaStore, notesStore),
	)
	protectedMux.HandleFunc(
		"PUT /update-notes",
		handler.UpdateProject(notesStore),
	)
	protectedMux.HandleFunc(
		"POST /commit",
		handler.PostSaveProject(
			stagingArea,
			metadataStore,
			mediaStore,
			notesStore,
		),
	)

	mux.Handle("/api/", middleware.Authenticate(
		http.StripPrefix("/api", protectedMux),
		revokedStore,
	))

	server := http.Server{
		Addr:              os.Getenv("PORT"),
		Handler:           mux,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
		ReadHeaderTimeout: 200 * time.Millisecond,
	}

	err = server.ListenAndServe()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error starting server: %v\n", err)
	}
}
