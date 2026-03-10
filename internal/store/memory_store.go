package store

import (
	"github.com/CuriousHet/geo-nearby-service/internal/models"
	"github.com/mmcloughlin/geohash"
)

type MemoryStore struct {
	Users        []models.User
	GeohashIndex map[uint]map[string][]models.User
}

func NewMemoryStore() *MemoryStore {

	users := []models.User{
		{ID: 1, Latitude: 23.0225, Longitude: 72.5714},
		{ID: 2, Latitude: 23.0300, Longitude: 72.5800},
		{ID: 3, Latitude: 19.0760, Longitude: 72.8777},
		{ID: 4, Latitude: 28.7041, Longitude: 77.1025},
		{ID: 5, Latitude: 12.9716, Longitude: 77.5946},
		{ID: 6, Latitude: 13.0827, Longitude: 80.2707},
		{ID: 7, Latitude: 22.5726, Longitude: 88.3639},
		{ID: 8, Latitude: 17.3850, Longitude: 78.4867},
		{ID: 9, Latitude: 26.9124, Longitude: 75.7873},
		{ID: 10, Latitude: 18.5204, Longitude: 73.8567},
		{ID: 11, Latitude: 23.2599, Longitude: 77.4126},
		{ID: 12, Latitude: 21.1702, Longitude: 72.8311},
		{ID: 13, Latitude: 22.3072, Longitude: 73.1812},
		{ID: 14, Latitude: 26.8467, Longitude: 80.9462},
		{ID: 15, Latitude: 25.3176, Longitude: 82.9739},
		{ID: 16, Latitude: 24.5854, Longitude: 73.7125},
		{ID: 17, Latitude: 30.7333, Longitude: 76.7794},
		{ID: 18, Latitude: 15.2993, Longitude: 74.1240},
		{ID: 19, Latitude: 11.0168, Longitude: 76.9558},
		{ID: 20, Latitude: 9.9312, Longitude: 76.2673},
		{ID: 21, Latitude: 16.5062, Longitude: 80.6480},
		{ID: 22, Latitude: 23.3441, Longitude: 85.3096},
		{ID: 23, Latitude: 26.1445, Longitude: 91.7362},
		{ID: 24, Latitude: 27.1767, Longitude: 78.0081},
		{ID: 25, Latitude: 29.9457, Longitude: 78.1642},
	}

	index := make(map[uint]map[string][]models.User)
	for p := uint(1); p <= 6; p++ {
		index[p] = make(map[string][]models.User)
		for _, u := range users {
			hash := geohash.EncodeWithPrecision(u.Latitude, u.Longitude, p)
			index[p][hash] = append(index[p][hash], u)
		}
	}

	return &MemoryStore{
		Users:        users,
		GeohashIndex: index,
	}
}
