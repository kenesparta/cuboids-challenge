package cmd

import (
	"bytes"
	"cuboid-challenge/app/config"
	"cuboid-challenge/app/db"
	"cuboid-challenge/app/db/migrations"
	"fmt"
	"log"
	"os"
	"text/template"
	"time"

	"github.com/spf13/cobra"
)

func migrateCmd() *cobra.Command { // nolint:funlen
	migrate := &cobra.Command{
		Use:   "migrate",
		Short: "database migrations tool",
		Run:   func(cmd *cobra.Command, args []string) {},
	}

	migrateGenerate := &cobra.Command{
		Use:   "generate",
		Short: "Generte a new migrations file",
		Args:  cobra.ExactValidArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			dataIn := struct {
				ID   string
				NAME string
			}{
				time.Now().Format("20060102150405"),
				args[0],
			}

			file, err := os.Create(fmt.Sprintf("app/db/migrations/%s_%s.go", dataIn.ID, dataIn.NAME))
			if err != nil {
				fmt.Println("Unable to create migration file:" + err.Error())

				return
			}
			defer file.Close()

			var out bytes.Buffer
			t := template.Must(template.ParseFiles("./cmd/migration.template"))
			if err := t.Execute(&out, dataIn); err != nil {
				fmt.Println("Unable to execute template:" + err.Error())

				return
			}

			if _, err := file.WriteString(out.String()); err != nil {
				fmt.Println("Unable to write to migration file:" + err.Error())

				return
			}

			fmt.Println("Generated new migration file ", file.Name())
		},
	}

	migrateUp := &cobra.Command{
		Use:   "up",
		Short: "Run migrations",
		Run: func(cmd *cobra.Command, args []string) {
			config.Load()
			m := migrations.Migrator(db.Connect())
			if err := m.Migrate(); err != nil {
				log.Fatalf("Could not migrate: %v", err)
			}
			log.Printf("Migration did run successfully")
		},
	}

	migrateDown := &cobra.Command{
		Use:   "down",
		Short: "Rollback last migrations",
		Run: func(cmd *cobra.Command, args []string) {
			config.Load()
			m := migrations.Migrator(db.Connect())
			if err := m.RollbackLast(); err != nil {
				log.Fatalf("Could not rollback last migration: %v", err)
			}
			log.Printf("Rollback run successfully")
		},
	}

	migrate.AddCommand(migrateGenerate, migrateUp, migrateDown)

	return migrate
}
