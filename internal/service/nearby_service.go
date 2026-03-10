package service

import (
	"github.com/CuriousHet/geo-nearby-service/internal/geo"
	"github.com/CuriousHet/geo-nearby-service/internal/models"
	"github.com/CuriousHet/geo-nearby-service/internal/store"
	"github.com/mmcloughlin/geohash"
)

type GeoBounds struct {
	MinLat float64 `json:"min_lat"`
	MaxLat float64 `json:"max_lat"`
	MinLon float64 `json:"min_lon"`
	MaxLon float64 `json:"max_lon"`
}

type NearbyResult struct {
	Users []models.User `json:"users"`
	Grid  []GeoBounds   `json:"grid"`
}

func getPrecisionForRadius(radiusKm float64) uint {
	// Geohash precisions (approximate cell dimensions at equator):
	// 1: 5,009.4km x 4,992.6km
	// 2: 1,252.3km x 624.1km
	// 3: 156.5km x 156km
	// 4: 39.1km x 19.5km
	// 5: 4.9km x 4.9km
	// 6: 1.2km x 0.61km
	// 7: 152.8m x 152.8m

	// To ensure a 3x3 grid covers a circle of `radiusKm`, 
	// the cell width/height should theoretically be > radiusKm 
	// (Actually, to be safe, cell size > radiusKm * 2, but evaluating 9 cells gives 3x3 length so safe limit is 1x size).
	
	if radiusKm <= 0.6 {
		return 6
	}
	if radiusKm <= 4.9 {
		return 5
	}
	if radiusKm <= 19.5 {
		return 4
	}
	if radiusKm <= 156.0 {
		return 3
	}
	if radiusKm <= 624.1 {
		return 2
	}
	return 1
}

func FindNearbyUsers(
	storeInstance *store.MemoryStore,
	myLat float64,
	myLon float64,
	radius float64,
) NearbyResult {

	precision := getPrecisionForRadius(radius)
	centerHash := geohash.EncodeWithPrecision(myLat, myLon, precision)
	neighbors := geohash.Neighbors(centerHash)

	allHashes := append(neighbors, centerHash)

	uniqueHashes := make(map[string]bool)
	var finalHashes []string
	for _, h := range allHashes {
		if !uniqueHashes[h] {
			uniqueHashes[h] = true
			finalHashes = append(finalHashes, h)
		}
	}

	grid := make([]GeoBounds, 0, len(finalHashes))
	for _, h := range finalHashes {
		box := geohash.BoundingBox(h)
		grid = append(grid, GeoBounds{
			MinLat: box.MinLat,
			MaxLat: box.MaxLat,
			MinLon: box.MinLng,
			MaxLon: box.MaxLng,
		})
	}

	var candidates []models.User
	if indexAtPrecision, ok := storeInstance.GeohashIndex[precision]; ok {
		for _, h := range allHashes {
			candidates = append(candidates, indexAtPrecision[h]...)
		}
	}

	var nearby []models.User
	minLat, maxLat, minLon, maxLon := geo.BoundingBox(myLat, myLon, radius)

	for _, user := range candidates {
		if user.Latitude < minLat || user.Latitude > maxLat || user.Longitude < minLon || user.Longitude > maxLon {
			continue
		}

		distance := geo.Haversine(
			myLat,
			myLon,
			user.Latitude,
			user.Longitude,
		)

		if distance <= radius {
			nearby = append(nearby, user)
		}
	}

	if nearby == nil {
		nearby = []models.User{}
	}

	return NearbyResult{
		Users: nearby,
		Grid:  grid,
	}
}