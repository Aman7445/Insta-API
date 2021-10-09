package main

import (
	"context"
	"encoding/json"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var lock sync.Mutex
//Posts Data Model
type Posts struct {
	ID              primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	UserID          string             `json:"userid,omitempty" bson:"userid,omitempty"`
	ImageURL        string             `json:"imageurl,omitempty" bson:"imageurl,omitempty"`
	Caption         string             `json:"caption,omitempty" bson:"caption,omitempty"`
	PostedTimeStamp time.Time          `json:"postedtimestamp,omitempty" bson:"postedtimestamp,omitempty"`
}

//Create Posts FUNCTION
func CreatePostEndpoint(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("content-type", "application/json")
	var post Posts
	_ = json.NewDecoder(request.Body).Decode(&post)
	post.PostedTimeStamp = time.Now()
	collection := client.Database("appointy").Collection("post")
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	result, _ := collection.InsertOne(ctx, post)
	json.NewEncoder(response).Encode(result)
}

//Get Posts by id Function 
func GetPostEndpoint(response http.ResponseWriter, request *http.Request) {
	lock.Lock()
    defer lock.Unlock()
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
	time.Sleep(1 * time.Second)
}

//Get Each Users Posts Function

func GetUserPostsEndpoint(response http.ResponseWriter, request *http.Request) {
	lock.Lock()
    defer lock.Unlock()
	response.Header().Add("content-type", "application/json")
	params := mux.Vars(request)
	id, _ := params["id"]
	var posts []Posts
	collection := client.Database("appointy").Collection("post")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	cursor, err := collection.Find(ctx, Posts{UserID: id})
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{"message": "` + err.Error() + `"}`))
		return
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var post Posts
		cursor.Decode(&post)
		posts = append(posts, post)
	}
	if err := cursor.Err(); err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{"message": "` + err.Error() + `"}`))
		return
	}
	json.NewEncoder(response).Encode(posts)
	time.Sleep(1 * time.Second)
}