package migrations

import (
	"fmt"

	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

func init() {
	migrations = append(migrations, &gormigrate.Migration{
		ID: "20220314230604",
		Migrate: func(tx *gorm.DB) error {
			fmt.Println("Running migration new_bag_disabled_field")
			type Bag struct{ Disabled bool }

			return tx.Migrator().AddColumn(&Bag{}, "disabled")
		},
		Rollback: func(tx *gorm.DB) error {
			fmt.Println("Rollback migration new_bag_disabled_field")
			type Bag struct{ Disabled bool }

			return tx.Migrator().DropColumn(&Bag{}, "disabled")
		},
	})
}
