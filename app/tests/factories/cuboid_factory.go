package factories

import (
	"cuboid-challenge/app/models"

	"github.com/brianvoe/gofakeit/v6"
)

func Cuboid() *models.Cuboid {
	return &models.Cuboid{
		Width:  uint(gofakeit.Uint32()),
		Height: uint(gofakeit.Uint32()),
		Depth:  uint(gofakeit.Uint32()),
	}
}
