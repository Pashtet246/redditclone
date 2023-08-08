package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	commentsHandler "redditclone/internal/api/server/comments/delivery"
	commentsRepo "redditclone/internal/api/server/comments/repo"
	middleware "redditclone/internal/api/server/middlewares"
	"redditclone/internal/api/server/posts"
	postsRepo "redditclone/internal/api/server/posts/repo"
	sessionRepo "redditclone/internal/api/server/session/repo"
	userHandler "redditclone/internal/api/server/users/delivery"
	"redditclone/internal/api/server/users/repo"
	"time"

	"github.com/go-redis/redis"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var ctx = context.Background()

func main() {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379",
		Password: "",
		DB:       0,
	})
	session := sessionRepo.NewSessionRepo(rdb)

	mysqlDb, err := sqlx.Connect("mysql", "root@tcp(localhost:3306)/coursera")
	if err != nil {
		panic(err)
	}

	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017/")
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}

	postsCollection := client.Database("redditclone").Collection("posts")
	commentsCollection := client.Database("redditclone").Collection("comments")

	postsRepo := postsRepo.NewPostsRepo(ctx, postsCollection)
	userHandler := &userHandler.UserHandler{
		UserRepo: repo.NewUserRepo(mysqlDb),
		Session:  session,
	}
	postHandler := &posts.PostsHandler{
		PostsRepo: postsRepo,
		UsersRepo: userHandler.UserRepo,
		Session:   session,
	}
	commentsHandler := &commentsHandler.CommentsHandler{
		CommentsRepo: commentsRepo.NewCommentsRepo(ctx, commentsCollection),
		Session:      session,
		PostsRepo:    postsRepo,
	}

	r := mux.NewRouter()
	// NOT REQUIRED SESSION
	r.HandleFunc("/api/login", userHandler.Login).Methods("POST")
	r.HandleFunc("/api/register", userHandler.Registration).Methods("POST")
	r.HandleFunc("/api/posts/", postHandler.ListAll).Methods("GET")
	r.HandleFunc("/api/posts/{CATEGORY_NAME}", postHandler.ListRefCategory).Methods("GET")
	r.HandleFunc("/api/post/{POST_ID}", postHandler.GetPost).Methods("GET")
	r.HandleFunc("/api/user/{USER_LOGIN}", postHandler.ListRefUser).Methods("GET")
	//REQUIRED SESSION
	r.HandleFunc("/api/posts", postHandler.Create).Methods("POST")
	r.HandleFunc("/api/post/{POST_ID}", postHandler.Remove).Methods("DELETE")
	r.HandleFunc("/api/post/{POST_ID}", commentsHandler.AddComment).Methods("POST")
	r.HandleFunc("/api/post/{POST_ID}/upvote", postHandler.VoteUp).Methods("GET")
	r.HandleFunc("/api/post/{POST_ID}/downvote", postHandler.VoteDown).Methods("GET")
	r.HandleFunc("/api/post/{POST_ID}/unvote", postHandler.RemoveVote).Methods("GET")
	r.HandleFunc("/api/post/{POST_ID}/{COMMENT_ID}", commentsHandler.DeleteComment).Methods("DELETE")

	r.PathPrefix("/").Handler(http.FileServer(http.Dir("../../template/")))

	mux := middleware.Auth(session, r)

	srv := &http.Server{
		Handler:      mux,
		Addr:         "127.0.0.1:8000",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	fmt.Println("starting server at :8000")
	log.Fatal(srv.ListenAndServe())
}
