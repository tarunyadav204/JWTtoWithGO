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
)

const todoconnectionString ="mongodb://localhost:27017/"
const tododbName = "React-Go-Auth"
const todocolName = "ToDo-Lists"

var todocollection *mongo.Collection
func init(){
 
	clientOptions :=  options.Client().ApplyURI(todoconnectionString)
	cxt,err := mongo.Connect(context.TODO(),clientOptions)
	if err!=nil{
		log.Fatal(err)
		return
	}
	   fmt.Println("MongoDB connected Successfully.........")
	   todocollection = cxt.Database(tododbName).Collection(todocolName)
       fmt.Println("Collection instance is ready .........")     
	}

func GetAllTodos(w http.ResponseWriter , r *http.Request){
	w.Header().Set("Content-Type" , "application/json")
	w.Header().Set("Allow-Control-Allow-Methods","GET")
	tokenStr := r.Header.Get("Authorization")
	if tokenStr == "" {
		http.Error(w, "Authorization header is missing", http.StatusUnauthorized)
		return
	}
	if strings.HasPrefix(tokenStr , "Bearer "){
		tokenStr = strings.TrimPrefix(tokenStr , "Bearer ")
	}
	     claims, err := utils.VerifyToken(tokenStr)
		 if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	    }

		//filter := bson.M{"user-email" : claims.UserEmail}
		filter := bson.M{"email" : claims.UserEmail}

		cursor,err := todocollection.Find(context.Background(),filter)
		if err != nil {
		http.Error(w, "Failed to fetch todos", http.StatusInternalServerError)
		return
	    }
			
		defer cursor.Close(context.Background())

		var todos []model.ToDoList
		if err = cursor.All(context.Background(), &todos); err != nil {
		http.Error(w, "Failed to parse todos", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(todos)
}

func CreateTodos(w http.ResponseWriter , r *http.Request){
	w.Header().Set("Content-Type" , "application/json")
	w.Header().Set("Allow-Control-Allow-Methods","POST")
	tokenStr := r.Header.Get("Authorization")
	if tokenStr == "" {
		http.Error(w, "Authorization header is missing", http.StatusUnauthorized)
		return
	}
	if strings.HasPrefix(tokenStr, "Bearer ") {
		tokenStr = strings.TrimPrefix(tokenStr, "Bearer ")
	}
		claims, err := utils.VerifyToken(tokenStr)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	var todo model.ToDoList
	if err := json.NewDecoder(r.Body).Decode(&todo); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}
	todo.Email = claims.UserEmail

	 insertResult, err := todocollection.InsertOne(context.Background(), todo)
	if err != nil {
		http.Error(w, "Failed to create todo", http.StatusInternalServerError)
		return
	}
  
	todo.ID = insertResult.InsertedID.(primitive.ObjectID)
    json.NewEncoder(w).Encode(todo)
	
}
/*
func DeleteTodos(w http.ResponseWriter , r *http.Request){
	w.Header().Set("Content-Type" , "application/json")
	w.Header().Set("Allow-Control-Allow-Methods","DELETE")
	tokenStr := r.Header.Get("Authorization")
	if tokenStr == "" {
		http.Error(w, "Authorization header is missing", http.StatusUnauthorized)
		return
	}
	if strings.HasPrefix(tokenStr, "Bearer ") {
		tokenStr = strings.TrimPrefix(tokenStr, "Bearer ")
	}

		claims, err := utils.VerifyToken(tokenStr)
       if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	var todo model.ToDoList
	if err := json.NewDecoder(r.Body).Decode(&todo); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// Ensure the todo item belongs to the authenticated user
	filter := bson.M{"_id": todo.ID, "email": claims.UserEmail}

	result, err := todocollection.DeleteOne(context.Background(), filter)
	if err != nil {
		http.Error(w, "Failed to delete todo", http.StatusInternalServerError)
		return
	}
	if result.DeletedCount == 0 {
		http.Error(w, "Todo not found or unauthorized", http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusNoContent)
	response := map[string]string{"message" : "Todo deleted successfully"}
   json.NewEncoder(w).Encode(response)

}
*/
func DeleteTodos(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Methods", "DELETE")

	tokenStr := r.Header.Get("Authorization")
	if tokenStr == "" {
		http.Error(w, `{"message": "Authorization header is missing"}`, http.StatusUnauthorized)
		return
	}

	if strings.HasPrefix(tokenStr, "Bearer ") {
		tokenStr = strings.TrimPrefix(tokenStr, "Bearer ")
	}

	claims, err := utils.VerifyToken(tokenStr)
	if err != nil {
		http.Error(w, `{"message": "Unauthorized"}`, http.StatusUnauthorized)
		return
	}

	todoID := r.URL.Query().Get("id")
	if todoID == "" {
		http.Error(w, `{"message": "Todo ID is required"}`, http.StatusBadRequest)
		return
	}

	objID, err := primitive.ObjectIDFromHex(todoID)
	if err != nil {
		http.Error(w, `{"message": "Invalid Todo ID"}`, http.StatusBadRequest)
		return
	}

	filter := bson.M{"_id": objID, "email": claims.UserEmail}

	result, err := todocollection.DeleteOne(context.Background(), filter)
	if err != nil {
		http.Error(w, `{"message": "Failed to delete todo"}`, http.StatusInternalServerError)
		return
	}

	if result.DeletedCount == 0 {
		http.Error(w, `{"message": "Todo not found or unauthorized"}`, http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
func UpdateToDo(w http.ResponseWriter , r *http.Request){
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Allow-Control-Allow-Methods", "PUT")
	tokenStr := r.Header.Get("Authorization")
	if tokenStr == "" {
		http.Error(w, "Authorization header is missing", http.StatusUnauthorized)
		return
	}
	if strings.HasPrefix(tokenStr, "Bearer ") {
		tokenStr = strings.TrimPrefix(tokenStr, "Bearer ")
	}
	claims, err := utils.VerifyToken(tokenStr)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
    var updatedTodo model.ToDoList
	if err := json.NewDecoder(r.Body).Decode(&updatedTodo); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}
	 if updatedTodo.ID.IsZero() {
        http.Error(w, `{"message": "Todo ID is required"}`, http.StatusBadRequest)
        return
    }
	// Ensure the todo item belongs to the authenticated user
	filter := bson.M{"_id": updatedTodo.ID, "email": claims.UserEmail}
	updateBSON := bson.M{}
	if updatedTodo.Title != ""{
		updateBSON["title"] = updatedTodo.Title
	}
	if updatedTodo.Content != ""{
		updateBSON["content"] = updatedTodo.Content
	}
	if updatedTodo.Completed != nil{
		updateBSON["completed"] = &updatedTodo.Completed
	}

	updateDocs := bson.M{"$set" : updateBSON}
	result, err := todocollection.UpdateOne(context.Background(), filter, updateDocs)

	if err != nil {
		http.Error(w, "Failed to update todo", http.StatusInternalServerError)
		return
	}

	if result.MatchedCount == 0 {
		http.Error(w, "Todo not found or unauthorized", http.StatusNotFound)
		return
	}

		json.NewEncoder(w).Encode(updatedTodo)

	/*
	update := bson.M{
		"$set": bson.M{
			"title":      updatedTodo.Title,
			"content" :  updatedTodo.Content,
			"completed": updatedTodo.Completed,
		},
	}
		result, err := todocollection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		http.Error(w, "Failed to update todo", http.StatusInternalServerError)
		return
	}

	if result.MatchedCount == 0 {
		http.Error(w, "Todo not found or unauthorized", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(updatedTodo) */
	}
