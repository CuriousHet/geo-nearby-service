package models

type User struct {
	ID        int
	Latitude  float64
	Longitude float64
}

// Point represents a lat/lon coordinate for visualization.
type Point struct {
	Lat float64 `json:"lat"`
	Lon float64 `json:"lon"`
}