package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"jwtAuth/model"
	"jwtAuth/utils"
	"log"
	"net/http"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
)


const connectionString ="mongodb://localhost:27017/"
const dbName = "React-Go-Auth"
const colName = "Auth-GO"

var collection *mongo.Collection
func init(){
 
	clientOptions :=  options.Client().ApplyURI(connectionString)
	cxt,err := mongo.Connect(context.TODO(),clientOptions)
	if err!=nil{
		log.Fatal(err)
		return
	}
	   fmt.Println("MongoDB connected Successfully.........")
	   collection = cxt.Database(dbName).Collection(colName)
       fmt.Println("Collection instance is ready .........")     
	}

	

func ServeHome(w http.ResponseWriter , r *http.Request) {
    w.Header().Set("Content-Type" , "application.json")
	w.Write([]byte("Hello World !"))
}

func registerUser(user model.Users) (model.Users, error) {
	// Hash the user's password
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return user, fmt.Errorf("internal server error: %w", err)
	}
	user.Password = string(hashPassword)

	token, err := utils.GenerateToken(user.Email)
	if err != nil {
		return user, fmt.Errorf("internal server error: %w", err)	
	}
	user.Token = token

	// Insert the user into the MongoDB collection
	insertResult, err := collection.InsertOne(context.Background(), user)
	if err != nil {
		return user, fmt.Errorf("failed to insert user: %w", err)
	}

	user.ID = insertResult.InsertedID.(primitive.ObjectID)
	return user, nil
}


func Register(w http.ResponseWriter, r *http.Request) {
	var user model.Users
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	registeredUser, err := registerUser(user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	token, err := utils.GenerateToken(registeredUser.Email)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

    
	response := map[string]string{
		"token": token,
	}
	/* response := model.Token{
		UserEmail : registeredUser.Email,
		Token:      token,
	}*/
	
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func Login(w http.ResponseWriter , r *http.Request){
	var user model.Users
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	var foundUser model.Users
	err := collection.FindOne(context.Background(), bson.M{"email": user.Email}).Decode(&foundUser)
	if err != nil {
		http.Error(w, "User not found", http.StatusUnauthorized)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(foundUser.Password), []byte(user.Password))
	if err != nil {
		http.Error(w, "Invalid password", http.StatusUnauthorized)
		return
	}

	token, err := utils.GenerateToken(foundUser.Email)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := map[string]string{
		"token": token,
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
	   
}
func Welcome(w http.ResponseWriter, r *http.Request) {
    tokenStr := r.Header.Get("Authorization")
    if tokenStr == "" {
        http.Error(w, "Authorization header is missing", http.StatusUnauthorized)
        return
    }

    // Remove "Bearer " prefix if it exists
    if strings.HasPrefix(tokenStr, "Bearer ") {
        tokenStr = strings.TrimPrefix(tokenStr, "Bearer ")
    }

    claims, err := utils.VerifyToken(tokenStr)
    if err != nil {
        http.Error(w, "Unauthorized", http.StatusUnauthorized)
        return
    }

    response := map[string]string{
       // "message": fmt.Sprintf("Welcome %s!", claims.UserEmail),
	   "message" : claims.UserEmail,
    }

    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(response)
}
