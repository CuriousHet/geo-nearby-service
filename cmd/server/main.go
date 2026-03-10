package main

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/CuriousHet/geo-nearby-service/internal/service"
	"github.com/CuriousHet/geo-nearby-service/internal/store"
)

func main() {

	store := store.NewMemoryStore()

	http.HandleFunc("/nearby", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Content-Type", "application/json")

		lat, _ := strconv.ParseFloat(r.URL.Query().Get("lat"), 64)
		lon, _ := strconv.ParseFloat(r.URL.Query().Get("lon"), 64)
		radius, _ := strconv.ParseFloat(r.URL.Query().Get("radius"), 64)

		result := service.FindNearbyUsers(
			store,
			lat,
			lon,
			radius,
		)

		json.NewEncoder(w).Encode(result)
	})

	http.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(store.Users)
	})

	http.Handle("/", http.FileServer(http.Dir("./web")))

	http.ListenAndServe(":8080", nil)
}