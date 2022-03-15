package models

import (
	"encoding/json"
	"fmt"
)

type Bag struct {
	Model

	Title  string `validate:"required,max=255"`
	Volume uint   `validate:"gt=0"`

	Cuboids []Cuboid
}

func (b *Bag) PayloadVolume() uint {
	var totalVolume uint
	for _, cuboid := range b.Cuboids {
		totalVolume += cuboid.PayloadVolume()
	}
	return totalVolume
}

func (b *Bag) AvailableVolume() uint {
	return b.Volume - b.PayloadVolume()
}

func (b *Bag) SetDisabled(value bool) {
}

func (b *Bag) MarshalJSON() ([]byte, error) {
	j, err := json.Marshal(struct {
		ID              uint     `json:"id"`
		Title           string   `json:"title"`
		Volume          uint     `json:"volume"`
		PayloadVolume   uint     `json:"payloadVolume"`
		AvailableVolume uint     `json:"availableVolume"`
		Cuboids         []Cuboid `json:"cuboids"`
	}{
		b.ID, b.Title, b.Volume, b.PayloadVolume(), b.AvailableVolume(), b.Cuboids,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to marshal Bag. %w", err)
	}

	return j, nil
}
