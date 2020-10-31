package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	goMigrate "github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/urfave/cli/v2"
)

func main() {
	app := cli.App{
		Name:  "migrate",
		Usage: "Database migration tool",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "db",
				Aliases:  []string{"d"},
				Usage:    "Database uri: driver://username:password@host:port/dbname?param1=true&param2=false",
				Required: true,
				EnvVars:  []string{"DB_URI"},
			},
			&cli.StringFlag{
				Name:     "path",
				Aliases:  []string{"p"},
				Usage:    "migration files: /path/to/migration/files",
				Required: true,
				EnvVars:  []string{"MIGRATIONS_PATH"},
			},
		},
		Commands: []*cli.Command{
			{
				Name:  "up",
				Usage: "Apply all or N up migrations",
				Action: func(c *cli.Context) error {
					var migrationSteps int
					var err error

					if c.Args().First() != "" {
						migrationSteps, err = strconv.Atoi(c.Args().First())
						if err != nil {
							log.Fatal(err)
						}
					}

					Up(c.String("path"), c.String("db"), migrationSteps)

					log.Println("Successfully applied migrations")

					return nil
				},
			},
			{
				Name:  "down",
				Usage: "Apply all or N down migrations",
				Action: func(c *cli.Context) error {
					var migrationSteps int
					var err error

					if c.Args().First() != "" {
						migrationSteps, err = strconv.Atoi(c.Args().First())
						if err != nil {
							log.Fatal(err)
						}
					}

					Down(c.String("path"), c.String("db"), migrationSteps)

					log.Println("Successfully applied migrations")

					return nil
				},
			},
			{
				Name:  "goto",
				Usage: "Migrate to version V",
				Action: func(c *cli.Context) error {
					var version int
					var err error

					if c.Args().First() != "" {
						version, err = strconv.Atoi(c.Args().First())
						if err != nil {
							log.Fatal(err)
						}
					}

					Migrate(c.String("path"), c.String("db"), uint(version))

					log.Println("Successfully applied migrations")

					return nil
				},
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func Up(path, db string, steps int) {
	m, err := goMigrate.New(fmt.Sprintf("file://%s", path), db)
	defer func() {
		serr, derr := m.Close()
		if serr != nil || derr != nil {
			log.Fatal(err)
		}
	}()
	if err != nil {
		log.Fatal(err)
	}

	if steps > 0 {
		err := m.Steps(steps)
		if err != nil {
			log.Fatal(err)
		}

		return
	}

	err = m.Up()
	if err == goMigrate.ErrNoChange {
		log.Println("No changes to apply")
		return
	}
	if err != nil {
		log.Fatal(err)
	}
}

func Down(path, db string, steps int) {
	m, err := goMigrate.New(fmt.Sprintf("file://%s", path), db)
	defer func() {
		serr, derr := m.Close()
		if serr != nil || derr != nil {
			log.Fatal(err)
		}
	}()
	if err != nil {
		log.Fatal(err)
	}

	if steps < 0 {
		err := m.Steps(steps)
		if err != nil {
			log.Fatal(err)
		}

		return
	}

	err = m.Down()
	if err == goMigrate.ErrNoChange {
		log.Println("No changes to apply")
		return
	}
	if err != nil {
		log.Fatal(err)
	}
}

func Migrate(path, db string, version uint) {
	m, err := goMigrate.New(fmt.Sprintf("file://%s", path), db)
	defer func() {
		serr, derr := m.Close()
		if serr != nil || derr != nil {
			log.Fatal(err)
		}
	}()
	if err != nil {
		log.Fatal(err)
	}

	if err := m.Migrate(version); err != nil {
		log.Fatal(err)
	}
}
