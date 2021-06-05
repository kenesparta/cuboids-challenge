package factories

import (
	"cuboid-challenge/app/models"

	"github.com/brianvoe/gofakeit/v6"
)

func Bag() *models.Bag {
	return &models.Bag{
		Title:  gofakeit.Name(),
		Volume: uint(gofakeit.Uint32()),
	}
}
