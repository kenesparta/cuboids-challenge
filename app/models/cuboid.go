package models

import (
	"encoding/json"
	"fmt"
)

type Cuboid struct {
	Model

	Width  uint `validate:"gt=0"`
	Height uint `validate:"gt=0"`
	Depth  uint `validate:"gt=0"`

	BagID uint
	Bag   *Bag
}

func (c *Cuboid) PayloadVolume() uint {
	return c.Width * c.Depth * c.Height
}

func (c *Cuboid) MarshalJSON() ([]byte, error) {
	jsonCuboid, err := json.Marshal(struct {
		ID     uint `json:"id"`
		Width  uint `json:"width"`
		Height uint `json:"height"`
		Depth  uint `json:"depth"`
		Volume uint `json:"volume"`
		BagID  uint `json:"bagId"`
	}{
		c.ID,
		c.Width,
		c.Height,
		c.Depth,
		c.PayloadVolume(),
		c.BagID,
	})
	if err != nil {
		err = fmt.Errorf("failed to marshal Cuboid. %w", err)
	}

	return jsonCuboid, err
}
