package store

import (
	"math/rand"

	"github.com/CuriousHet/geo-nearby-service/internal/geo"
	"github.com/CuriousHet/geo-nearby-service/internal/models"
	"github.com/golang/geo/s2"
	"github.com/mmcloughlin/geohash"
)

type MemoryStore struct {
	Users        []models.User
	GeohashIndex map[uint]map[string][]models.User
	S2Index      map[int]map[s2.CellID][]models.User
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

	ghIndex := make(map[uint]map[string][]models.User)
	for p := uint(1); p <= 6; p++ {
		ghIndex[p] = make(map[string][]models.User)
	}

	s2Index := make(map[int]map[s2.CellID][]models.User)
	// We'll index from Level 3 (~700km) to Level 13 (~800m)
	for lv := 3; lv <= 13; lv++ {
		s2Index[lv] = make(map[s2.CellID][]models.User)
	}

	// Index all 100,000 users
	for _, u := range users {
		// Populate Geohash Index
		for p := uint(1); p <= 6; p++ {
			hash := geohash.EncodeWithPrecision(u.Latitude, u.Longitude, p)
			ghIndex[p][hash] = append(ghIndex[p][hash], u)
		}

		// Populate S2 Index for each level
		for lv := 3; lv <= 13; lv++ {
			cellID := geo.LatLonToS2CellID(u.Latitude, u.Longitude, lv)
			s2Index[lv][cellID] = append(s2Index[lv][cellID], u)
		}
	}

	return &MemoryStore{
		Users:        users,
		GeohashIndex: ghIndex,
		S2Index:      s2Index,
	}
}
