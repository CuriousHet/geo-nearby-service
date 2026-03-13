package geo

import (
	"github.com/golang/geo/s1"
	"github.com/golang/geo/s2"
)

// S2Level defines the precision of our S2 cells.
// Level 13 provides a cell size of roughly 1.27km x 1.27km.
const S2Level = 13

// LatLonToS2CellID converts a latitude and longitude to an S2 CellID at a specific level.
func LatLonToS2CellID(lat, lon float64, level int) s2.CellID {
	level = validateLevel(level)
	ll := s2.LatLngFromDegrees(lat, lon)
	return s2.CellIDFromLatLng(ll).Parent(level)
}

// LatLonToCellToken returns the S2 cell token string for a given coordinate and level.
func LatLonToCellToken(lat, lon float64, level int) string {
	return LatLonToS2CellID(lat, lon, level).ToToken()
}

// GetS2LevelForRadius chooses an appropriate S2 level based on radius to avoid cell explosion.
func GetS2LevelForRadius(radiusKm float64) int {
	switch {
	case radiusKm <= 1.3:
		return 13
	case radiusKm <= 2.6:
		return 12
	case radiusKm <= 5.3:
		return 11
	case radiusKm <= 10.6:
		return 10
	case radiusKm <= 21:
		return 9
	case radiusKm <= 42:
		return 8
	case radiusKm <= 85:
		return 7
	case radiusKm <= 170:
		return 6
	default:
		return 5
	}
}

// GetS2CoveringCells calculates the list of S2 CellIDs required to cover a circular region at a specific level.
func GetS2CoveringCells(lat, lon, radiusKm float64, level int) s2.CellUnion {
	level = validateLevel(level)
	center := s2.LatLngFromDegrees(lat, lon)

	// S2 Cap represents a spherical cap (a circle on the sphere)
	// We convert our radius in KM to an angle in radians on the Earth's surface
	const earthRadiusKm = 6371.01

	// Create the cap
	cap := s2.CapFromCenterAngle(
		s2.PointFromLatLng(center),
		s1.Angle(radiusKm/earthRadiusKm),
	)

	// RegionCoverer is the tool that finds which cells best tile our circle
	rc := &s2.RegionCoverer{
		MaxCells: 128, // Increased accuracy for covering
		MinLevel: level,
		MaxLevel: level,
	}

	// Return the covering cell union
	return rc.Covering(cap)
}

// GetS2CellPolygon returns the 5 vertices (closed loop) of an S2 Cell for visualization on the map.
func GetS2CellPolygon(cellID s2.CellID) [][]float64 {
	cell := s2.CellFromCellID(cellID)
	vertices := make([][]float64, 5)
	for i := 0; i < 4; i++ {
		v := s2.LatLngFromPoint(cell.Vertex(i))
		vertices[i] = []float64{v.Lat.Degrees(), v.Lng.Degrees()}
	}
	// Close the polygon by repeating the first vertex
	vertices[4] = vertices[0]
	return vertices
}

// validateLevel ensures the S2 level is within the valid range [0, 30].
func validateLevel(level int) int {
	if level < 0 || level > 30 {
		return S2Level
	}
	return level
}
