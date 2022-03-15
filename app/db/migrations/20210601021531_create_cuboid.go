package migrations

import (
	"cuboid-challenge/app/models"
	"fmt"

	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

func init() {
	migrations = append(migrations, &gormigrate.Migration{
		ID: "20210601021531",
		Migrate: func(trx *gorm.DB) error {
			fmt.Println("Running migration create_cuboid")
			type Cuboid struct {
				models.Model
				Width  uint
				Height uint
				Depth  uint
				BagID  uint
			}

			return trx.AutoMigrate(&Cuboid{})
		},
		Rollback: func(trx *gorm.DB) error {
			fmt.Println("Rollback migration create_cuboid")
			type Cuboid struct{}

			return trx.Migrator().DropTable(&Cuboid{})
		},
	})
}
