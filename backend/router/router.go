package router

import (
	"jwtAuth/controller"
	//"net/http"

	"github.com/gorilla/mux"
)

func Router() *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("/", controller.ServeHome).Methods("GET")
	router.HandleFunc("/register", controller.Register).Methods("POST", "OPTIONS")
	router.HandleFunc("/login", controller.Login).Methods("POST", "OPTIONS")
	router.HandleFunc("/welcome", controller.Welcome).Methods("GET")

	// Global CORS handling
	/*router.Methods(http.MethodOptions).HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})*/

	// Todo routes
	router.HandleFunc("/todos", controller.GetAllTodos).Methods("GET")
	router.HandleFunc("/createtodo",controller.CreateTodos).Methods("POST")
	router.HandleFunc("/deletetodo",controller.DeleteTodos).Methods("DELETE")
	router.HandleFunc("/updatetodo",controller.UpdateToDo).Methods("PUT")


	return router
}
