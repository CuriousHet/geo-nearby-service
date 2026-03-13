package service

import (
	"fmt"
	"time"

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

type TimingStats struct {
	GeohashTime     string `json:"geohash_time"`
	S2Time          string `json:"s2_time"`
	BoundingBoxTime string `json:"bounding_box_time"`
	HaversineTime   string `json:"haversine_time"`
}

type NearbyResult struct {
	Users       []models.User  `json:"users"`
	Grid        []GeoBounds    `json:"grid"`     // For Geohash (rectangles)
	Polygons    [][]models.Point `json:"polygons"` // For S2 (true cell geometry)
	TotalDBSize int            `json:"total_db_size"`
	Timing      TimingStats    `json:"timing"`
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

func getS2CandidatesPhase(storeInstance *store.MemoryStore, myLat float64, myLon float64, radius float64) ([]models.User, [][]models.Point, string) {
	startTime := time.Now()

	// Choose appropriate S2 level based on radius
	level := geo.GetS2LevelForRadius(radius)

	// Get the cells that cover our search circle at this level
	covering := geo.GetS2CoveringCells(myLat, myLon, radius, level)
	
	var candidates []models.User
	var polygons [][]models.Point

	if indexAtLevel, ok := storeInstance.S2Index[level]; ok {
		for _, cellID := range covering {
			if users, ok := indexAtLevel[cellID]; ok {
				candidates = append(candidates, users...)
			}

			// Prepare polygons for visualization
			vertices := geo.GetS2CellPolygon(cellID)
			poly := make([]models.Point, len(vertices))
			for i, v := range vertices {
				poly[i] = models.Point{Lat: v[0], Lon: v[1]}
			}
			polygons = append(polygons, poly)
		}
	}

	duration := fmt.Sprintf("%.4f ms", time.Since(startTime).Seconds()*1000.0)
	return candidates, polygons, duration
}

func getGeohashCandidatesPhase(storeInstance *store.MemoryStore, myLat float64, myLon float64, radius float64) ([]models.User, []GeoBounds, string) {
	startTime := time.Now()

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

	duration := fmt.Sprintf("%.4f ms", time.Since(startTime).Seconds()*1000.0)
	return candidates, grid, duration
}

func filterByBoundingBoxPhase(candidates []models.User, myLat float64, myLon float64, radius float64) ([]models.User, string) {
	startTime := time.Now()
	var bboxCandidates []models.User
	minLat, maxLat, minLon, maxLon := geo.BoundingBox(myLat, myLon, radius)

	for _, user := range candidates {
		if user.Latitude < minLat || user.Latitude > maxLat || user.Longitude < minLon || user.Longitude > maxLon {
			continue
		}
		bboxCandidates = append(bboxCandidates, user)
	}

	duration := fmt.Sprintf("%.4f ms", time.Since(startTime).Seconds()*1000.0)
	return bboxCandidates, duration
}

func filterByHaversinePhase(candidates []models.User, myLat float64, myLon float64, radius float64) ([]models.User, string) {
	startTime := time.Now()
	var nearby []models.User
	for _, user := range candidates {
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

	duration := fmt.Sprintf("%.4f ms", time.Since(startTime).Seconds()*1000.0)
	return nearby, duration
}

func FindNearbyUsers(
	storeInstance *store.MemoryStore,
	myLat float64,
	myLon float64,
	radius float64,
	algorithm string,
) NearbyResult {

	var candidates []models.User
	var grid []GeoBounds
	polygons := [][]models.Point{}
	var geohashDuration, s2Duration, bboxDuration, havDuration string

	geohashDuration = "-"
	s2Duration = "-"
	bboxDuration = "-"
	havDuration = "-"

	if algorithm == "naive" {
		// Phase 1: Pure Haversine O(N) over entire DB
		candidates = storeInstance.Users
		nearby, dur := filterByHaversinePhase(candidates, myLat, myLon, radius)
		havDuration = dur
		
		return NearbyResult{
			Users:       nearby,
			Grid:        []GeoBounds{},
			TotalDBSize: len(storeInstance.Users),
			Timing: TimingStats{
				GeohashTime:     geohashDuration,
				S2Time:          s2Duration,
				BoundingBoxTime: bboxDuration,
				HaversineTime:   havDuration,
			},
		}
	} else if algorithm == "bounding_box" {
		// Phase 2: Bounding Box O(N) over entire DB, then Haversine
		candidates = storeInstance.Users
		bboxCandidates, dur1 := filterByBoundingBoxPhase(candidates, myLat, myLon, radius)
		bboxDuration = dur1
		
		nearby, dur2 := filterByHaversinePhase(bboxCandidates, myLat, myLon, radius)
		havDuration = dur2
		
		return NearbyResult{
			Users:       nearby,
			Grid:        []GeoBounds{},
			TotalDBSize: len(storeInstance.Users),
			Timing: TimingStats{
				GeohashTime:     geohashDuration,
				S2Time:          s2Duration,
				BoundingBoxTime: bboxDuration,
				HaversineTime:   havDuration,
			},
		}
	} else if algorithm == "s2" {
		// Phase 4: Google S2 O(1), then Bounding Box O(K), then Haversine O(M)
		candidates, polygons, s2Duration = getS2CandidatesPhase(storeInstance, myLat, myLon, radius)
		bboxCandidates, dur1 := filterByBoundingBoxPhase(candidates, myLat, myLon, radius)
		bboxDuration = dur1
		
		nearby, dur2 := filterByHaversinePhase(bboxCandidates, myLat, myLon, radius)
		havDuration = dur2
		
		return NearbyResult{
			Users:       nearby,
			Grid:        []GeoBounds{},
			Polygons:    polygons,
			TotalDBSize: len(storeInstance.Users),
			Timing: TimingStats{
				GeohashTime:     geohashDuration,
				S2Time:          s2Duration,
				BoundingBoxTime: bboxDuration,
				HaversineTime:   havDuration,
			},
		}
	} else {
		// Phase 3: Geohash O(1), then Bounding Box O(K), then Haversine O(M)
		candidates, grid, geohashDuration = getGeohashCandidatesPhase(storeInstance, myLat, myLon, radius)
		bboxCandidates, dur1 := filterByBoundingBoxPhase(candidates, myLat, myLon, radius)
		bboxDuration = dur1
		
		nearby, dur2 := filterByHaversinePhase(bboxCandidates, myLat, myLon, radius)
		havDuration = dur2
		
		return NearbyResult{
			Users:       nearby,
			Grid:        grid,
			Polygons:    polygons,
			TotalDBSize: len(storeInstance.Users),
			Timing: TimingStats{
				GeohashTime:     geohashDuration,
				S2Time:          s2Duration,
				BoundingBoxTime: bboxDuration,
				HaversineTime:   havDuration,
			},
		}
	}
}