package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)
var client *mongo.Client

type Person struct{
	ID primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name string `json:"name,omitempty" bson:"name,omitempty"`
}
type Posts struct{
	ID primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	UserID string `json:"_id,omitempty" bson:"_id,omitempty"`
	Post string `json:"post,omitempty" bson:"post,omitempty"`
}

func CreateUserEndpoint(response http.ResponseWriter,request *http.Request){
	response.Header().Set("content-type","application/json")
	var user Person
	_=json.NewDecoder(request.Body).Decode(&user)
	collection := client.Database("appointy").Collection("person")
	ctx, _:= context.WithTimeout(context.Background(), 5*time.Second)
	result, _:=collection.InsertOne(ctx,user);
	json.NewEncoder(response).Encode(result)
}
func GetUserEndpoint(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("content-type", "application/json")
	params := mux.Vars(request)
	id, _ := primitive.ObjectIDFromHex(params["id"])
	var user Person
	collection := client.Database("appointy").Collection("person")
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	err := collection.FindOne(ctx, Person{ID: id}).Decode(&user)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}
	json.NewEncoder(response).Encode(user)
}
func CreatePostEndpoint(response http.ResponseWriter,request *http.Request){
	response.Header().Set("content-type","application/json")
	var post Posts
	_=json.NewDecoder(request.Body).Decode(&post)
	collection := client.Database("appointy").Collection("post")
	ctx, _:= context.WithTimeout(context.Background(), 5*time.Second)
	result, _:=collection.InsertOne(ctx,post);
	json.NewEncoder(response).Encode(result)
}
func GetPostEndpoint(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("content-type", "application/json")
	params := mux.Vars(request)
	id, _ := primitive.ObjectIDFromHex(params["id"])
	var post Posts
	collection := client.Database("appointy").Collection("post")
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	err := collection.FindOne(ctx, Posts{ID: id}).Decode(&post)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}
	json.NewEncoder(response).Encode(post)
}
func GetUserPostEndpoint(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("content-type", "application/json")
	params := mux.Vars(request)
	userid  := string(params["userid"])
	var upost Posts
	collection := client.Database("appointy").Collection("post")
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	err := collection.FindOne(ctx,Posts{UserID:userid}).Decode(&upost)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}
	json.NewEncoder(response).Encode(upost)
}


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
	router.HandleFunc("/posts/users/{userid}", GetUserPostEndpoint).Methods("GET")
	http.ListenAndServe(":8080",router)
}
