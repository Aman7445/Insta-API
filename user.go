package main

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

//Users Data Model
type Users struct{
	ID primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name string `json:"name,omitempty" bson:"name,omitempty"`
	Email string `json:"email,omitempty" bson:"email,omitempty"`
	Password  string  `json:"password,omitempty" bson:"password,omitempty"`
}

//Create User Function
func CreateUserEndpoint(response http.ResponseWriter, request *http.Request){
	response.Header().Add("content-type", "application/json")
	var user Users
	json.NewDecoder(request.Body).Decode(&user)
	hash, _ := HashPassword(user.Password)
	user.Password = hash
	collection := client.Database("appointy").Collection("users")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	result, _ := collection.InsertOne(ctx, user)
	json.NewEncoder(response).Encode(result)
}
//Get User by id function
func GetUserEndpoint(response http.ResponseWriter, request *http.Request) {
	lock.Lock()
    defer lock.Unlock()
	response.Header().Set("content-type", "application/json")
	params := mux.Vars(request)
	id, _ := primitive.ObjectIDFromHex(params["id"])
	var user Users
	collection := client.Database("appointy").Collection("users")
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	err := collection.FindOne(ctx, Users{ID: id}).Decode(&user)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}
	json.NewEncoder(response).Encode(user)
	time.Sleep(1 * time.Second)
}

func HashPassword(password string) (string, error) {
	lock.Lock()
    defer lock.Unlock()
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	time.Sleep(1 * time.Second)
    return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	lock.Lock()
    defer lock.Unlock()
    err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	time.Sleep(1 * time.Second)
    return err == nil
	
}