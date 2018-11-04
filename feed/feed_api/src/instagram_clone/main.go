package main

import (
  "os"
  "net/http"
  "github.com/gorilla/mux"
  "instagram_clone/handlers"
)

func main() {



  handlers.Set_up_client()
  defer handlers.Redis_client.Close()

  router := mux.NewRouter()

  router.HandleFunc("/timeline/{id:[0-9]+}", handlers.TimelineHandler).Methods("GET")

  http.ListenAndServe(port(), router)

}


func port() string {
	port := os.Getenv("PORT")
	if len(port) == 0 {
		port = "8080"
	}
	return ":" + port
}
