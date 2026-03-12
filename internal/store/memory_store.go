package store

import (
	"math/rand"

	"github.com/CuriousHet/geo-nearby-service/internal/models"
	"github.com/mmcloughlin/geohash"
)

type MemoryStore struct {
	Users        []models.User
	GeohashIndex map[uint]map[string][]models.User
}

func NewMemoryStore() *MemoryStore {

	var users []models.User
	
	// Create 100,000 mock users globally distributed for stress-testing O(1) retrieval
	for i := 1; i <= 100000; i++ {
		// Latitudes range from -90 to 90
		lat := -90.0 + (rand.Float64() * 180.0)
		// Longitudes range from -180 to 180
		lon := -180.0 + (rand.Float64() * 360.0)
		
		users = append(users, models.User{
			ID:        i,
			Latitude:  lat,
			Longitude: lon,
		})
	}

	index := make(map[uint]map[string][]models.User)
	for p := uint(1); p <= 6; p++ {
		index[p] = make(map[string][]models.User)
	}

	// Index all 100,000 users
	for _, u := range users {
		for p := uint(1); p <= 6; p++ {
			hash := geohash.EncodeWithPrecision(u.Latitude, u.Longitude, p)
			index[p][hash] = append(index[p][hash], u)
		}
	}

	return &MemoryStore{
		Users:        users,
		GeohashIndex: index,
	}
}
