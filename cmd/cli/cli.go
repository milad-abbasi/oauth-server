package main

import (
	"log"
	"os"
	"strconv"

	"github.com/urfave/cli/v2"
)

func main() {
	app := cli.App{
		Name:  "oauth",
		Usage: "Oauth cli",
		Commands: []*cli.Command{
			{
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
				Subcommands: []*cli.Command{
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
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
