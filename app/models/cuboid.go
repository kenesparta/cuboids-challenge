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
	return 0
}

func (c *Cuboid) MarshalJSON() ([]byte, error) {
	b, err := json.Marshal(struct {
		ID     uint `json:"id"`
		Width  uint `json:"width"`
		Height uint `json:"height"`
		Depth  uint `json:"depth"`
	}{
		c.ID, c.Width, c.Height, c.Depth,
	})
	if err != nil {
		err = fmt.Errorf("failed to marshal Cuboid. %w", err)
	}

	return b, err
}
