package migrator

import (
	"database/sql"
	"github.com/pressly/goose"
	"log"
)

const (
	migrDir = "./migrations"
)

type GooseMigrator struct {
	db *sql.DB
}

func NewGooseMigrator(db *sql.DB) *GooseMigrator {
	return &GooseMigrator{
		db: db,
	}
}

func (m *GooseMigrator) Commands(command string, args ...string) error {
	return m.gooseRunner(command, args...)
}

func (m *GooseMigrator) gooseRunner(command string, args ...string) error {
	if err := goose.Run(command, m.db, migrDir, args...); err != nil {
		log.Println(err)
		return err
	}
	return nil
}
