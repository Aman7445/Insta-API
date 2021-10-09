package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)
var client *mongo.Client


func main(){
	fmt.Println("Starting the Application... ")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	clientOptions:= options.Client().ApplyURI("mongodb+srv://aman123:aman123@cluster0.pkzni.mongodb.net/myFirstDatabase?retryWrites=true&w=majority")
	client, _=mongo.Connect(ctx,clientOptions)
	router :=mux.NewRouter();
	router.HandleFunc("/users", CreateUserEndpoint).Methods("POST")
	router.HandleFunc("/users/{id}", GetUserEndpoint).Methods("GET")
	router.HandleFunc("/posts", CreatePostEndpoint).Methods("POST")
	router.HandleFunc("/posts/{id}", GetPostEndpoint).Methods("GET")
	router.HandleFunc("/posts/users/{id}", GetUserPostsEndpoint).Methods("GET")
	http.ListenAndServe(":8080",router)
}
