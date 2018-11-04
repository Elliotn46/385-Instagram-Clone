package handlers

import (
  "github.com/gorilla/mux"
  "net/http"
  "encoding/json"

)

func TimelineHandler(w http.ResponseWriter, r *http.Request) {
  id := mux.Vars(r)
  res := get_timeline(id["id"])
  json.NewEncoder(w).Encode(res)
}

func ServerUnavailableHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusServiceUnavailable)
	w.Write([]byte("Service Unavailable"))
}
