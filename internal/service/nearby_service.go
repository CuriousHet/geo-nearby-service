package service

import(
	"github.com/CuriousHet/geo-nearby-service/internal/geo"
	"github.com/CuriousHet/geo-nearby-service/internal/models"
)

func FindNearbyUsers(

	users []models.User,
	myLat float64,
	myLon float64,
	radius float64,
	
) []models.User {
	
	var nearby []models.User

	minLat, maxLat, minLon, maxLon := geo.BoundingBox(myLat, myLon, radius)

	for _, user := range users {
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

	return nearby
}	