package service

import (
	"cuboid-challenge/app/db"
	"cuboid-challenge/app/models"
	"errors"

	"gorm.io/gorm"
)

var (
	ErrNotFound             = errors.New("not found")
	ErrInsufficientCapacity = errors.New("insufficient capacity in bag")
	ErrBagDisabled          = errors.New("bag is disabled")
)

type CuboidInput struct {
	ID     uint `json:"id"`
	Width  uint
	Height uint
	Depth  uint
	BagID  uint `json:"bagId"`
}

type CuboidService struct {
	Connection  *gorm.DB
	cuboidInput *CuboidInput
}

func NewCuboidService(conn *gorm.DB, cuboidInput *CuboidInput) *CuboidService {
	return &CuboidService{Connection: conn, cuboidInput: cuboidInput}
}

// validate Validates the cuboidInput that comes from the request.
func (cs *CuboidService) validate() (*models.Cuboid, error) {
	var (
		cuboid = models.Cuboid{
			Width:  cs.cuboidInput.Width,
			Height: cs.cuboidInput.Height,
			Depth:  cs.cuboidInput.Depth,
			BagID:  cs.cuboidInput.BagID,
		}
		bag models.Bag
	)

	if r := db.CONN.Preload("Cuboids").First(&bag, cuboid.BagID); r.Error != nil {
		return nil, r.Error
	}

	if !bag.HasCapacity(cuboid.PayloadVolume()) {
		return nil, ErrInsufficientCapacity
	}

	if bag.Disabled {
		return nil, ErrBagDisabled
	}

	return &cuboid, nil
}

func (cs *CuboidService) Create() (*models.Cuboid, error) {
	cuboidValid, err := cs.validate()
	if err != nil {
		return nil, err
	}

	if cs.Connection.Create(cuboidValid).Error != nil {
		return nil, cs.Connection.Error
	}

	return cuboidValid, nil
}

func (cs *CuboidService) Update() (*models.Cuboid, error) {
	cuboidValid, err := cs.validate()
	if err != nil {
		return nil, err
	}

	cuboidValid.ID = cs.cuboidInput.ID

	update := cs.Connection.Updates(cuboidValid)
	if update.RowsAffected == 0 {
		return nil, ErrNotFound
	}

	if update.Error != nil {
		return nil, cs.Connection.Error
	}

	return cuboidValid, nil
}
