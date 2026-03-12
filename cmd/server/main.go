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
		algorithm := r.URL.Query().Get("algorithm")

		result := service.FindNearbyUsers(
			store,
			lat,
			lon,
			radius,
			algorithm,
		)

		json.NewEncoder(w).Encode(result)
	})

	http.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Content-Type", "application/json")
		
		// Only return a tiny subset to prevent crashing the browser on full reload
		subset := store.Users
		if len(subset) > 100 {
			subset = subset[:100]
		}
		json.NewEncoder(w).Encode(subset)
	})

	http.Handle("/", http.FileServer(http.Dir("./web")))

	http.ListenAndServe(":8080", nil)
}