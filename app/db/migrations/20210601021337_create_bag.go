package migrations

import (
	"cuboid-challenge/app/models"
	"fmt"

	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

func init() {
	migrations = append(migrations, &gormigrate.Migration{
		ID: "20210601021337",
		Migrate: func(trx *gorm.DB) error {
			fmt.Println("Running migration create_bag")
			type Bag struct {
				models.Model
				Title  string
				Volume uint
			}

			return trx.AutoMigrate(&Bag{})
		},
		Rollback: func(trx *gorm.DB) error {
			fmt.Println("Rollback migration create_bag")
			type Bag struct{}

			return trx.Migrator().DropTable(&Bag{})
		},
	})
}
