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

// type Users struct{
// 	ID primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
// 	Name string `json:"name,omitempty" bson:"name,omitempty"`
// 	Email string `json:"email,omitempty" bson:"email,omitempty"`
// 	Password  string  `json:"password,omitempty" bson:"password,omitempty"`
// }
type Posts struct{
	ID primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	UserID string `json:"userid,omitempty" bson:"userid,omitempty"`
	ImageURL string `json:"imageurl,omitempty" bson:"imageurl,omitempty"`
	Caption  string `json:"caption,omitempty" bson:"caption,omitempty"`
	PostTime time.Duration `json:"time,omitempty" bson:"time,omitempty"`
}
// func CreateUserEndpoint(response http.ResponseWriter, request *http.Request){
// 	response.Header().Add("content-type", "application/json")
// 	var user Users
// 	json.NewDecoder(request.Body).Decode(&user)
// 	hash, _ := HashPassword(user.Password)
// 	user.Password = hash
// 	collection := client.Database("appointy").Collection("users")
// 	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
// 	result, _ := collection.InsertOne(ctx, user)
// 	json.NewEncoder(response).Encode(result)
// }
// func GetUserEndpoint(response http.ResponseWriter, request *http.Request) {
// 	response.Header().Set("content-type", "application/json")
// 	params := mux.Vars(request)
// 	id, _ := primitive.ObjectIDFromHex(params["id"])
// 	var user Users
// 	collection := client.Database("appointy").Collection("users")
// 	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
// 	err := collection.FindOne(ctx, Users{ID: id}).Decode(&user)
// 	if err != nil {
// 		response.WriteHeader(http.StatusInternalServerError)
// 		response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
// 		return
// 	}
// 	json.NewEncoder(response).Encode(user)
// }
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

func GetUserPostsEndpoint(response http.ResponseWriter, request *http.Request){
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
    for cursor.Next(ctx){
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

}

// func HashPassword(password string) (string, error) {
// 	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
//     return string(bytes), err
// }

// func CheckPasswordHash(password, hash string) bool {
//     err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
//     return err == nil
// }

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
