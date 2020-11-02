package main

import (
	"fmt"
	"log"

	goMigrate "github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

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
