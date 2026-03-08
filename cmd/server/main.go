package main

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/CuriousHet/geo-nearby-service/internal/store"
	"github.com/CuriousHet/geo-nearby-service/internal/service"
)

func main() {

	store := store.NewMemoryStore()

	http.HandleFunc("/nearby", func(w http.ResponseWriter, r *http.Request) {

		lat, _ := strconv.ParseFloat(r.URL.Query().Get("lat"), 64)
		lon, _ := strconv.ParseFloat(r.URL.Query().Get("lon"), 64)
		radius, _ := strconv.ParseFloat(r.URL.Query().Get("radius"), 64)

		users := service.FindNearbyUsers(
			store.Users,
			lat,
			lon,
			radius,
		)

		json.NewEncoder(w).Encode(users)
	})

	http.ListenAndServe(":8080", nil)
}