package geo

import "math"

func BoundingBox(lat, lon, radius float64) (float64, float64, float64, float64) {

	latDelta := radius / 111.0

	lonDelta := radius / (111.0 * math.Cos(lat*math.Pi/180))

	minLat := lat - latDelta
	maxLat := lat + latDelta
	minLon := lon - lonDelta
	maxLon := lon + lonDelta

	return minLat, maxLat, minLon, maxLon
}