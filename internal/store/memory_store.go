package store 
import (
	"github.com/CuriousHet/geo-nearby-service/internal/models"
)

type MemoryStore struct {
	Users []models.User
}

func NewMemoryStore() *MemoryStore {

	users := []models.User{
		{ID: 1, Latitude: 23.0225, Longitude: 72.5714},
		{ID: 2, Latitude: 23.0300, Longitude: 72.5800},
		{ID: 3, Latitude: 19.0760, Longitude: 72.8777},
	}

	return &MemoryStore{
		Users: users,
	}
}