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

	for _, user := range users {

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